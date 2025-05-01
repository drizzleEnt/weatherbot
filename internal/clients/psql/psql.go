package psql

import (
	"context"
	"database/sql"
	"fmt"
	"weatherbot/internal/clients"
)

type dbClient struct {
	db *sql.DB
}

var _ clients.DBClient = (*dbClient)(nil)

func NewClient(ctx context.Context, dsn string) (*dbClient, error) {
	cl, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed open sql client: %w", err)
	}

	err = cl.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed ping sql db: %w", err)
	}

	return &dbClient{
		db: cl,
	}, nil
}

// Close implements clients.DBClient.
func (d *dbClient) Close() error {
	return d.Close()
}

// DB implements clients.DBClient.
func (d *dbClient) DB() *sql.DB {
	return d.db
}
