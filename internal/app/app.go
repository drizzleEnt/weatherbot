package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weatherbot/internal/config"
	"weatherbot/internal/config/env"
	"weatherbot/internal/routes"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type App struct {
	httpServer *http.Server

	httpCfg config.HTTPConfig

	log *zap.Logger
}

func New(ctx context.Context) (*App, error) {
	a := App{}

	inits := []func(context.Context) error{
		a.initLogger,
		a.initConfig,
		a.initHttpServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return nil, err
		}
	}
	return &a, nil
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.runHttpServer(); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			a.log.Error("failed run server", zap.Error(err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	a.log.Info("Shutting down HTTP server...")
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		a.log.Error("failed shutdown server", zap.Error(err))
		return err
	}
	a.log.Info("HTTP server stopped")

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	if a.log == nil {
		logger := zap.New(getCore(getAtomicLevel()))

		a.log = logger
	}

	return nil
}

func (a *App) initHttpServer(ctx context.Context) error {
	if a.httpServer == nil {
		mux := routes.SetupRoutes()

		srv := &http.Server{
			Addr:           a.httpConfig().Address(),
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		a.httpServer = srv
	}

	return nil
}

func (a *App) runHttpServer() error {
	a.log.Info("http server start", zap.String("address", a.httpServer.Addr))

	return a.httpServer.ListenAndServe()
}

func (a *App) httpConfig() config.HTTPConfig {
	if a.httpCfg == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			a.log.Error("failed load http config", zap.Error(err))
			os.Exit(1)
		}

		a.httpCfg = cfg
	}

	return a.httpCfg
}

var logFilePath string

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10,
		MaxAge:     7,
		MaxBackups: 3,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	develomentCfg := zap.NewDevelopmentEncoderConfig()
	develomentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(develomentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

var logLevel = flag.String("l", "info", "log level")

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level

	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
