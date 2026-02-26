package daemon_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/etc/db/migrations"
	"github.com/jsnfwlr/filamate/internal/cmd/daemon"
	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/test/containers"
)

func TestDaemonErrors(t *testing.T) {
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

	cases := []struct {
		name             string
		PostgresHostname string
		PostgresPort     string
		PostgresUsername string
		PostgresPassword string
		PostgresDatabase string
		APIHost          string
		APIPort          string
		APIStaticType    string
		StartDB          bool
		Error            error
	}{
		{
			name:             "invalid db host",
			PostgresHostname: "local host",
			PostgresPort:     "5432",
			PostgresUsername: "filament",
			PostgresPassword: "filament",
			PostgresDatabase: "filament",
			StartDB:          false,
			Error:            fmt.Errorf("invalid config: invalid value %q for field %q", "local host", "POSTGRES_HOST"),
		},
		{
			name:             "invalid db port",
			PostgresHostname: "localhost",
			PostgresPort:     "invalid",
			PostgresUsername: "filament",
			PostgresPassword: "filament",
			PostgresDatabase: "filament",
			StartDB:          false,
			Error:            fmt.Errorf("invalid config: invalid value %q for field %q", "invalid", "POSTGRES_PORT"),
		},
		{
			name:             "db connection error",
			PostgresHostname: "localhost",
			PostgresPort:     "5432",
			PostgresUsername: "filament",
			PostgresPassword: "filament",
			PostgresDatabase: "filament",
			StartDB:          false,
			Error:            errors.New("failed to create migrator: could not connect to database \"postgres://filament:REDACTED_SECRET@localhost:5432/filament?sslmode=disable\": failed to connect to `user=filament database=filament`: 127.0.0.1:5432 (localhost): failed SASL auth: FATAL: password authentication failed for user \"filament\" (SQLSTATE 28P01)"),
		},
		{
			name:          "invalid API host",
			StartDB:       true,
			APIHost:       "http://invalid-host",
			APIPort:       "123",
			APIStaticType: "embedded",
			Error:         errors.New("config validation failed: Key: 'EnvConfig.Host' Error:Field validation for 'Host' failed on the 'hostname|ip' tag"),
		},
	}

	var ctr containers.DatabaseContainer
	var cfg db.Config
	var containerStarted bool

	for _, c := range cases {
		waitTime := 1000 * time.Millisecond
		if c.StartDB {
			if !containerStarted {
				ctr, cfg, err = containers.Postgres(t, ctx, "db_version", "17", new("daemon-progressive-test"))
				if err != nil {
					t.Fatalf("could not start the Postgres container: %v", err)
				}
				containerStarted = true

				t.Cleanup(func() {
					ctr.Cleanup(t)
				})
			}
			waitTime = 1 * time.Second
		}

		t.Run(c.name, func(t *testing.T) {
			t.Setenv("API_HOST", c.APIHost)
			t.Setenv("API_PORT", c.APIPort)
			t.Setenv("API_STATIC_TYPE", c.APIStaticType)

			if !containerStarted {
				t.Setenv("POSTGRES_HOST", c.PostgresHostname)
				t.Setenv("POSTGRES_PORT", c.PostgresPort)
				t.Setenv("POSTGRES_USER", c.PostgresUsername)
				t.Setenv("POSTGRES_PASSWORD", c.PostgresPassword)
				t.Setenv("POSTGRES_DATABASE", c.PostgresDatabase)
			} else {
				t.Setenv("POSTGRES_HOST", cfg.Host)
				t.Setenv("POSTGRES_PORT", cfg.Port)
				t.Setenv("POSTGRES_USER", cfg.Username)
				t.Setenv("POSTGRES_PASSWORD", string(cfg.Password))
				t.Setenv("POSTGRES_DATABASE", cfg.Database)
			}

			cmd := daemon.StartCmd
			cmd.SetContext(ctx)
			err = daemon.StartRun(cmd, []string{})
			switch {
			case c.Error == nil && err != nil:
				t.Errorf("unexpected error: %v", err)
			case c.Error != nil && err == nil:
				t.Errorf("expected error but got none")
			case c.Error != nil && err != nil:
				if err.Error() != c.Error.Error() {
					t.Logf("expected error: %v", c.Error)
					t.Logf("received error: %v", err)
					t.Errorf("Incorrect error:\nexpected: %v\nreceived: %v", c.Error, err)
				}
			}

			if c.StartDB {
				time.Sleep(waitTime)
			}
		})
	}
}

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

	ctr, cfg, err := containers.Postgres(t, ctx, "db_version", "17", new("daemon-start-cmd-test"))
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
		cmd := daemon.StartCmd
		cmd.SetContext(ctx)
		_ = daemon.StartRun(cmd, []string{})
	}()

	time.Sleep(2 * time.Second)

	err = daemon.Healthcheck(ctx)
	if err != nil {
		t.Fatalf("healthcheck failed: %v", err)
	}

	go func() {
		cmd := daemon.CheckCmd
		cmd.SetContext(ctx)
		err = daemon.CheckRun(cmd, []string{})
		if err != nil {
			t.Errorf("check command failed: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)
}
