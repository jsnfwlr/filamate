package database_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/etc/db/migrations"
	"github.com/jsnfwlr/filamate/internal/cmd"
	"github.com/jsnfwlr/filamate/internal/cmd/database"
	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/test/containers"
)

func TestMigrateCmd(t *testing.T) {
	gCfg, err := go11y.LoadConfig()
	if err != nil {
		t.Fatalf("could not load go11y config: %v", err)
	}

	ctx, _, err := go11y.Initialise(context.Background(), gCfg, os.Stdout)
	if err != nil {
		t.Fatalf("could not initialise go11y: %v", err)
	}

	ctr, cfg, err := containers.Postgres(t, ctx, "db_version", "17", new("database-migrate-cmd-test"))
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

	testCases := []struct {
		name      string
		config    db.Config
		args      []string
		expectErr error
	}{
		{
			name:      "migrate to latest",
			config:    cfg,
			args:      []string{"database", "migrate", "--version", "-1"},
			expectErr: nil,
		},
		{
			name:      "reset to version 0",
			config:    cfg,
			args:      []string{"database", "migrate", "--version", "0"},
			expectErr: nil,
		},
		{
			name:      "migrate to version 1",
			config:    cfg,
			args:      []string{"database", "migrate", "--version", "1"},
			expectErr: nil,
		},
		{
			name:      "migrate to version a",
			config:    cfg,
			args:      []string{"database", "migrate", "--version", "a"},
			expectErr: errors.New(`invalid argument "a" for "--version" flag: strconv.ParseInt: parsing "a": invalid syntax`),
		},
		{
			name:   "migrate to version -2",
			config: cfg,
			args:   []string{"database", "migrate", "--version", "-2"},

			expectErr: fmt.Errorf("error running migrations: failed to set migration version: could not migrate to -2: destination version -2 is outside the valid versions of 0 to %d", col.Steps()),
		},
		{
			name:      "migrate to version 99999",
			config:    cfg,
			args:      []string{"database", "migrate", "--version", "99999"},
			expectErr: fmt.Errorf("error running migrations: failed to set migration version: could not migrate to 99999: destination version 99999 is outside the valid versions of 0 to %d", col.Steps()),
		},
		{
			name:      "migrate to latest version",
			config:    cfg,
			args:      []string{"database", "migrate", "--version", fmt.Sprintf("%d", col.Steps())},
			expectErr: nil,
		},
	}

	buf := bytes.Buffer{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()

			t.Setenv("POSTGRES_HOST", tc.config.Host)
			t.Setenv("POSTGRES_PORT", tc.config.Port)
			t.Setenv("POSTGRES_DATABASE", tc.config.Database)
			t.Setenv("POSTGRES_USER", tc.config.Username)
			t.Setenv("POSTGRES_PASSWORD", string(tc.config.Password))

			err = cmd.Execute(ctx, &buf, io.Discard, database.MigrateCmd, tc.args...)
			switch {
			case tc.expectErr == nil && err != nil:
				t.Errorf("unexpected error: %v", err)
				return
			case tc.expectErr != nil && err == nil:
				t.Errorf("missing error %v", tc.expectErr)
				return
			case tc.expectErr != nil && err != nil && err.Error() != tc.expectErr.Error():
				t.Errorf("incorrect error:\nexpected:\n\t%v\nreceived:\n\t%v", tc.expectErr, err)
				return
			}
		})
	}
}
