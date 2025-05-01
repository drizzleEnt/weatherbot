package psql

import (
	"context"
	"fmt"
	"weatherbot/internal/clients"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

type dbClient struct {
	db *pgx.Conn
}

var _ clients.DBClient = (*dbClient)(nil)

func NewClient(ctx context.Context, dsn string) (*dbClient, error) {
	cl, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed open sql client: %w", err)
	}

	err = cl.Ping(ctx)
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
func (d *dbClient) DB() *pgx.Conn {
	return d.db
}
