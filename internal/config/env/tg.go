package env

import (
	"errors"
	"os"
)

const (
	botToken = "BOT_TOKEN"
	botHost  = "BOT_HOST"
)

type botConfig struct {
	token string
	host  string
}

func NewBotConfig() (*botConfig, error) {
	token := os.Getenv(botToken)
	if len(token) == 0 {
		return nil, errors.New("bot token not found")
	}

	host := os.Getenv(botHost)
	if len(host) == 0 {
		return nil, errors.New("tg bot host not found")
	}

	return &botConfig{
		token: token,
		host:  host,
	}, nil
}

func (cfg *botConfig) Token() string {
	return cfg.token
}

func (cfg *botConfig) Host() string {
	return cfg.host
}
