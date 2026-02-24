package db

import (
	"context"
	"fmt"

	"github.com/jsnfwlr/filamate/etc/db/migrations"

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

func (c *Client) LoadDemoData(ctx context.Context) error {
	done, err := c.Queries.CheckDemoData(ctx)
	if err != nil {
		return fmt.Errorf("could not check for demo data: %w", err)
	}

	if done {
		return nil
	}

	dd, err := migrations.DemoData.ReadFile("demo_data.sql")
	if err != nil {
		return fmt.Errorf("could not read demo data: %w", err)
	}

	_, err = c.pool.Exec(ctx, string(dd))
	if err != nil {
		return fmt.Errorf("could not execute demo data: %w", err)
	}

	err = c.Queries.SetDemoData(ctx)
	if err != nil {
		return fmt.Errorf("could not set demo data: %w", err)
	}

	return nil
}
