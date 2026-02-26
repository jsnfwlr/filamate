package cmd_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/etc/db/migrations"
	"github.com/jsnfwlr/filamate/internal/cmd"
	"github.com/jsnfwlr/filamate/internal/cmd/daemon"
	"github.com/jsnfwlr/filamate/internal/cmd/database"
	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/test/containers"
)

func TestStartDaemon(t *testing.T) {
	gCfg, err := go11y.LoadConfig()
	if err != nil {
		t.Fatalf("could not load go11y config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
	})

	ctx, _, err = go11y.Initialise(ctx, gCfg, os.Stdout)
	if err != nil {
		t.Fatalf("could not initialise go11y: %v", err)
	}

	ctr, cfg, err := containers.Postgres(t, ctx, "db_version", "17", new("root-cmd-test"))
	if err != nil {
		t.Fatalf("could not start the Postgres container: %v", err)
	}

	t.Cleanup(func() {
		ctr.Cleanup(t)
	})

	col, err := migrations.New()
	if err != nil {
		t.Fatalf("could not create the embedded filesystem: %v", err)
	}

	t.Setenv("POSTGRES_HOST", cfg.Host)
	t.Setenv("POSTGRES_PORT", cfg.Port)
	t.Setenv("POSTGRES_DATABASE", cfg.Database)
	t.Setenv("POSTGRES_USER", cfg.Username)
	t.Setenv("POSTGRES_PASSWORD", string(cfg.Password))

	m, err := db.NewMigrator(ctx, cfg, col)
	if err != nil {
		t.Fatalf("failed to create migrator: %v", err)
	}

	_, err = m.Migrate(ctx)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	go func() {
		buf := bytes.Buffer{}
		err = cmd.Execute(ctx, nil, nil, database.MigrateCmd, "--version=-1")
		err = cmd.Execute(ctx, &buf, &buf, daemon.StartCmd)
	}()

	time.Sleep(5 * time.Second)

	err = daemon.Healthcheck(ctx)
	if err != nil {
		t.Fatalf("healthcheck failed: %v", err)
	}
}
