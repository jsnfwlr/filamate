package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool    *pgxpool.Pool
	Queries *Queries
}

func Connect(ctx context.Context, cfg Config) (dbClient *Client, fault error) {
	p, err := pgxpool.New(ctx, cfg.GetURI())
	if err != nil {
		return nil, fmt.Errorf("could not create connection pool: %w", err)
	}

	q := New(p)

	return &Client{
		pool:    p,
		Queries: q,
	}, nil
}

func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, opts)
}
