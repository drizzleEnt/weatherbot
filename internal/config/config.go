package config

import "github.com/subosito/gotenv"

func Load(path string) error {
	err := gotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type HTTPConfig interface {
	Address() string
}

type BotConfig interface {
	Token() string
	Host() string
}

type PGConfig interface {
	Address() string
}
