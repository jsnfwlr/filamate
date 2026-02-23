package containers_test

import (
	"context"
	"strings"
	"testing"

	"github.com/jsnfwlr/filamate/internal/test/containers"
	. "github.com/jsnfwlr/filamate/internal/types"
)

func TestPostgresContainer(t *testing.T) {
	ctx := context.Background()

	t.Run("valid_version_with_name", func(t *testing.T) {
		ctr, cfg, err := containers.Postgres(t, ctx, "db_version", "17", PointerOf("container-test"))
		if err != nil {
			t.Fatalf("error starting Postgres container: %v", err)
		}
		defer ctr.Cleanup(t)

		if cfg.Host == "" {
			t.Error("expected non-empty host")
		}

		if cfg.Port == "" {
			t.Error("expected non-empty port")
		}

		if cfg.Database != "filamate_test" {
			t.Errorf("expected database 'filamate_test', got %s", cfg.Database)
		}

		if cfg.Username != "user" {
			t.Errorf("expected username 'user', got %s", cfg.Username)
		}

		if cfg.Password != "password" {
			t.Errorf("expected password 'password', got %s", cfg.Password)
		}

		if cfg.VersionTable != "db_version" {
			t.Errorf("expected version table 'db_version', got %s", cfg.VersionTable)
		}

		t.Logf("Postgres container started successfully with postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	})

	t.Run("valid_version_without_name", func(t *testing.T) {
		ctr, cfg, err := containers.Postgres(t, ctx, "db_version", "16", nil)
		if err != nil {
			t.Fatalf("error starting Postgres container: %v", err)
		}
		defer ctr.Cleanup(t)

		if cfg.Host == "" {
			t.Error("expected non-empty host")
		}

		if cfg.Port == "" {
			t.Error("expected non-empty port")
		}
	})

	t.Run("invalid_version", func(t *testing.T) {
		ctr, _, err := containers.Postgres(t, ctx, "db_version", "bad_version", nil)
		if err == nil {
			t.Error("expected error for invalid postgres version")
		}
		defer ctr.Cleanup(t)
	})

	t.Run("empty_version", func(t *testing.T) {
		ctr, _, err := containers.Postgres(t, ctx, "db_version", "", nil)
		if err == nil {
			t.Error("expected error for empty postgres version")
		}
		defer ctr.Cleanup(t)
	})

	t.Run("invalid_version_table", func(t *testing.T) {
		ctr, _, err := containers.Postgres(t, ctx, "invalid table name", "17", PointerOf("test-invalid-table"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer ctr.Cleanup(t)

		// Container should start but the version table name will cause issues later
		// This tests that the container creation itself doesn't validate table names
	})

	t.Run("very_long_name", func(t *testing.T) {
		longName := strings.Repeat("a", 100)
		ctr, cfg, err := containers.Postgres(t, ctx, "db_version", "17", PointerOf(longName))
		if err != nil {
			// This might fail due to Docker name limitations, which is expected
			t.Logf("expected failure for very long name: %v", err)
			return
		}
		defer ctr.Cleanup(t)

		if cfg.Database != "filamate_test" {
			t.Errorf("expected database 'filamate_test', got %s", cfg.Database)
		}
	})

	t.Run("empty_context", func(t *testing.T) {
		// Test with background context to ensure no panics
		emptyCtx := context.Background()
		ctr, cfg, err := containers.Postgres(t, emptyCtx, "db_version", "17", PointerOf("empty-context-test"))
		if err != nil {
			t.Fatalf("error with background context: %v", err)
		}
		defer ctr.Cleanup(t)

		if cfg.Host == "" {
			t.Error("expected non-empty host even with background context")
		}
	})
}

func TestDatabaseContainerCleanup(t *testing.T) {
	ctx := context.Background()

	t.Run("cleanup_nil_postgres", func(t *testing.T) {
		ctr := containers.DatabaseContainer{}

		// Should not panic when calling cleanup on nil Postgres
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("cleanup panicked with nil postgres: %v", r)
			}
		}()

		ctr.Cleanup(t)
	})

	t.Run("cleanup_valid_container", func(t *testing.T) {
		ctr, _, err := containers.Postgres(t, ctx, "db_version", "17", PointerOf("cleanup-test"))
		if err != nil {
			t.Fatalf("error starting container for cleanup test: %v", err)
		}

		// Should not panic when calling cleanup on valid container
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("cleanup panicked with valid container: %v", r)
			}
		}()

		ctr.Cleanup(t)
	})
}
