package eventconsumer

import (
	"errors"
	"fmt"
	"time"
	"weatherbot/internal/events"

	"go.uber.org/zap"
)

type consumer struct {
	log       *zap.Logger
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int, log *zap.Logger) *consumer {
	return &consumer{
		log:       log,
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *consumer) Start() error {
	for {
		gotevent, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			c.log.Error("failed fetch", zap.Error(err))
		}

		if len(gotevent) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvent(gotevent); err != nil {
			c.log.Error("failed handle msg", zap.Error(err))
			continue
		}
	}
}

func (c *consumer) handleEvent(events []events.Event) error {
	c.log.Info("Start process ...")
	errs := make([]error, 0)
	for _, e := range events {
		c.log.Info("got new event:", zap.String("msg", e.Text))

		if err := c.processor.Process(e); err != nil {
			err = fmt.Errorf("cant process event: %w", err)
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}
