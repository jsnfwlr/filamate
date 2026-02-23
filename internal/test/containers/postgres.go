// Package containers provides utilities for managing test containers.
package containers

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsnfwlr/filamate/internal/db"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type DatabaseContainer struct {
	Postgres *postgres.PostgresContainer
}

// Postgres starts a Postgres container for testing purposes.
func Postgres(t *testing.T, ctx context.Context, versionTable, version string, name *string) (container DatabaseContainer, config db.Config, fault error) {
	t.Helper()

	var err error
	cfg := db.Config{
		Database:     "filamate_test",
		Username:     "user",
		Password:     "password",
		VersionTable: versionTable,
	}

	dbContainer := DatabaseContainer{}

	// name := fmt.Sprintf("filamate-test-postgres-%s", version)

	n := fmt.Sprintf("filamate-test-postgres-%s", version)
	if name != nil {
		n = *name
	}

	opts := []testcontainers.ContainerCustomizer{
		postgres.WithDatabase(cfg.Database),
		postgres.WithUsername(cfg.Username),
		postgres.WithPassword(string(cfg.Password)),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
		testcontainers.WithName(n),
	}

	if name != nil {
		opts = append(opts, testcontainers.WithReuseByName(n))
	}

	dbContainer.Postgres, err = postgres.Run(
		ctx,
		fmt.Sprintf("postgres:%s", version),
		opts...,
	)
	if err != nil {
		return DatabaseContainer{}, db.Config{}, err
	}

	cfg.Host, _ = dbContainer.Postgres.Host(ctx)
	port, _ := dbContainer.Postgres.MappedPort(ctx, "5432")
	cfg.Port = port.Port()

	return dbContainer, cfg, nil
}

func (c DatabaseContainer) Cleanup(t *testing.T) {
	if c.Postgres != nil {
		testcontainers.CleanupContainer(t, c.Postgres)
	}
}
