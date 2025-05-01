package env

import (
	"errors"
	"fmt"
	"os"
)

const (
	dbname     = "PG_DATABASE_NAME"
	dbuser     = "PG_USER"
	dbpassword = "PG_PASSWORD"
	dbport     = "PG_PORT"
	dbhost     = "PG_HOST"
	dbssl      = "PG_SSL"
)

type pgConfig struct {
	dsn string
}

func NewPGConfig() (*pgConfig, error) {
	name := os.Getenv(dbname)
	if len(name) == 0 {
		return nil, errors.New("pg name not found")
	}

	user := os.Getenv(dbuser)
	if len(user) == 0 {
		return nil, errors.New("pg user not found")
	}

	port := os.Getenv(dbport)
	if len(port) == 0 {
		return nil, errors.New("pg port not found")
	}

	host := os.Getenv(dbhost)
	if len(host) == 0 {
		return nil, errors.New("pg host not found")
	}

	password := os.Getenv(dbpassword)

	if len(password) == 0 {
		return nil, errors.New("pg password not found")
	}

	ssl := os.Getenv(dbssl)
	if len(ssl) == 0 {
		return nil, errors.New("pg sslmode not found")
	}

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		host, port, name, user, password, ssl)

	return &pgConfig{
		dsn: dsn,
	}, nil

}

func (cfg *pgConfig) Address() string {
	return cfg.dsn
}
