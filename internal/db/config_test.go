package db_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jsnfwlr/filamate/etc/db/migrations"
	"github.com/jsnfwlr/filamate/internal/db"
)

func TestMigrationFS(t *testing.T) {
	col, err := migrations.New()
	if err != nil {
		t.Fatalf("could not create the embedded filesystem: %v", err)
	}

	fi, err := col.ReadDir(".")
	if err != nil {
		t.Fatalf("could not read the directory: %v", err)
	}

	if len(fi) == 0 {
		t.Fatalf("no files found in the directory")
	}

	for _, f := range fi {
		t.Logf("name: %s, size: %d, mode: %s, modTime: %v, isDir: %t", f.Name(), f.Size(), f.Mode(), f.ModTime(), f.IsDir())
	}

	sharedPaths, err := col.Glob(filepath.Join("*", "*.sql"))
	if err != nil {
		t.Errorf("could not get globs: %s", err)
	}

	for _, p := range sharedPaths {
		t.Logf("path: %s", p)
	}

	// Test invalid glob patterns
	_, err = col.Glob("invalid[")
	if err == nil {
		t.Error("expected error for invalid glob pattern")
	}

	// Test non-existent directory
	_, err = col.ReadDir("nonexistent")
	if err == nil {
		t.Error("expected error reading non-existent directory")
	}
}

func TestConfig(t *testing.T) {
	def, err := db.LoadConfig()
	if err != nil {
		t.Fatalf("could not load config: %v", err)
	}

	if def.Host != "localhost" {
		t.Errorf("incorrect host: expected 'localhost', got '%s'", def.Host)
	}

	if def.Port != "5432" {
		t.Errorf("incorrect port: expected '5432', got '%s'", def.Port)
	}

	if def.Database != "filament" {
		t.Errorf("incorrect database: expected 'filament', got '%s'", def.Database)
	}

	if def.Username != "filament" {
		t.Errorf("incorrect username: expected 'filament', got '%s'", def.Username)
	}

	expectedURI := "postgres://filament:filament@localhost:5432/filament?sslmode=disable"

	if def.GetURI() != expectedURI {
		t.Errorf("incorrect URI: expected '%s', got '%s'", expectedURI, def.GetURI())
	}

	expectedRedactedURI := "postgres://filament:REDACTED_SECRET@localhost:5432/filament?sslmode=disable"

	if def.GetRedactedURI() != expectedRedactedURI {
		t.Errorf("incorrect redacted URI: expected '%s', got '%s'", expectedRedactedURI, def.GetRedactedURI())
	}

	t.Setenv("POSTGRES_HOST", "other-host")
	t.Setenv("POSTGRES_PORT", "5443")
	t.Setenv("POSTGRES_DATABASE", "filament2")
	t.Setenv("POSTGRES_USER", "filament2")
	t.Setenv("POSTGRES_PASSWORD", "filament2")

	updated, err := db.LoadConfig()
	if err != nil {
		t.Fatalf("could not load updated config: %v", err)
	}

	if updated.Host != "other-host" {
		t.Errorf("incorrect updated host: expected 'other-host', got '%s'", updated.Host)
	}

	if updated.Port != "5443" {
		t.Errorf("incorrect updated port: expected '5443', got '%s'", updated.Port)
	}

	if updated.Database != "filament2" {
		t.Errorf("incorrect updated database: expected 'filament2', got '%s'", updated.Database)
	}

	if updated.Username != "filament2" {
		t.Errorf("incorrect updated username: expected 'filament2', got '%s'", updated.Username)
	}

	expectedUpdatedURI := "postgres://filament2:filament2@other-host:5443/filament2?sslmode=disable"
	if updated.GetURI() != expectedUpdatedURI {
		t.Errorf("incorrect updated URI: expected '%s', got '%s'", expectedUpdatedURI, updated.GetURI())
	}

	expectedUpdatedRedactedURI := "postgres://filament2:REDACTED_SECRET@other-host:5443/filament2?sslmode=disable"
	if updated.GetRedactedURI() != expectedUpdatedRedactedURI {
		t.Errorf("incorrect updated redacted URI: expected '%s', got '%s'", expectedUpdatedRedactedURI, updated.GetRedactedURI())
	}
}

func TestConfigEdgeCases(t *testing.T) {
	t.Run("empty_values", func(t *testing.T) {
		t.Setenv("POSTGRES_HOST", "")
		t.Setenv("POSTGRES_PORT", "")
		t.Setenv("POSTGRES_DATABASE", "")
		t.Setenv("POSTGRES_USER", "")
		t.Setenv("POSTGRES_PASSWORD", "")
		t.Setenv("POSTGRES_VERSION_TABLE", "")

		cfg, err := db.LoadConfig()
		if err != nil {
			t.Fatalf("unexpected error with empty values: %v", err)
		}

		// Should use defaults
		if cfg.Host != "localhost" {
			t.Errorf("expected default host, got %s", cfg.Host)
		}
		if cfg.Port != "5432" {
			t.Errorf("expected default port, got %s", cfg.Port)
		}
	})

	t.Run("special_characters", func(t *testing.T) {
		t.Setenv("POSTGRES_HOST", "host-with-dashes")
		t.Setenv("POSTGRES_DATABASE", "db_with_underscores")
		t.Setenv("POSTGRES_USER", "user.with.dots")
		t.Setenv("POSTGRES_PASSWORD", "pass@word!123")

		cfg, err := db.LoadConfig()
		if err != nil {
			t.Fatalf("unexpected error with special characters: %v", err)
		}

		uri := cfg.GetURI()
		if !strings.Contains(uri, "host-with-dashes") {
			t.Error("URI should contain hostname with dashes")
		}
		if !strings.Contains(uri, "db_with_underscores") {
			t.Error("URI should contain database name with underscores")
		}
		if !strings.Contains(uri, "user.with.dots") {
			t.Error("URI should contain username with dots")
		}
		if !strings.Contains(uri, "pass@word!123") {
			t.Error("URI should contain password with special characters")
		}
	})

	t.Run("numeric_values", func(t *testing.T) {
		t.Setenv("POSTGRES_HOST", "192.168.1.100")
		t.Setenv("POSTGRES_PORT", "15432")
		t.Setenv("POSTGRES_DATABASE", "db123")

		cfg, err := db.LoadConfig()
		if err != nil {
			t.Fatalf("unexpected error with numeric values: %v", err)
		}

		if cfg.Host != "192.168.1.100" {
			t.Errorf("expected IP address host, got %s", cfg.Host)
		}
		if cfg.Port != "15432" {
			t.Errorf("expected custom port, got %s", cfg.Port)
		}
	})

	t.Run("invalid_port", func(t *testing.T) {
		t.Setenv("POSTGRES_PORT", "not-a-number")

		_, err := db.LoadConfig()
		if err == nil {
			t.Error("expected error for invalid port, but got none")
		}
	})

	t.Run("invalid host", func(t *testing.T) {
		t.Setenv("POSTGRES_HOST", "invalid hostname with spaces")

		_, err := db.LoadConfig()
		if err == nil {
			t.Error("expected error for invalid host, but got none")
		}
	})

	t.Run("invalid_database_name", func(t *testing.T) {
		t.Setenv("POSTGRES_DATABASE", "invalid database name with spaces")

		_, err := db.LoadConfig()
		if err == nil {
			t.Error("expected error for invalid database name, but got none")
		}
	})

	t.Run("invalid_username", func(t *testing.T) {
		t.Setenv("POSTGRES_USER", "invalid username with spaces")

		_, err := db.LoadConfig()
		if err == nil {
			t.Error("expected error for invalid username, but got none")
		}
	})

	t.Run("invalid_password", func(t *testing.T) {
		db, _ := db.LoadConfig()
		db.Password = ""

		err := db.Validate()
		if err == nil {
			t.Error("expected error for empty password, but got none")
		}
	})
}

func TestSecretType(t *testing.T) {
	t.Run("log_value_redaction", func(t *testing.T) {
		secret := db.Secret("super-secret-password")
		logValue := secret.LogValue()

		if logValue.String() != "REDACTED_SECRET" {
			t.Errorf("expected 'REDACTED_SECRET', got %s", logValue.String())
		}
	})

	t.Run("config_password_methods", func(t *testing.T) {
		cfg := db.Config{
			Password: db.Secret("test-password"),
		}

		if cfg.GetPassword() != "REDACTED_SECRET" {
			t.Errorf("expected 'REDACTED_SECRET', got %s", cfg.GetPassword())
		}

		// Test with empty password
		cfg.Password = db.Secret("")
		if cfg.GetPassword() != "UNSET" {
			t.Errorf("expected 'UNSET' for empty password, got %s", cfg.GetPassword())
		}
	})
}

func TestConfigMethods(t *testing.T) {
	cfg := db.Config{
		Host:         "test-host",
		Port:         "9999",
		Database:     "test-db",
		Username:     "test-user",
		Password:     db.Secret("test-pass"),
		VersionTable: "test_version",
	}

	if cfg.GetHost() != "test-host" {
		t.Errorf("expected 'test-host', got %s", cfg.GetHost())
	}

	if cfg.GetPort() != "9999" {
		t.Errorf("expected '9999', got %s", cfg.GetPort())
	}

	if cfg.GetDatabase() != "test-db" {
		t.Errorf("expected 'test-db', got %s", cfg.GetDatabase())
	}

	if cfg.GetUsername() != "test-user" {
		t.Errorf("expected 'test-user', got %s", cfg.GetUsername())
	}

	if cfg.GetVersionTable() != "test_version" {
		t.Errorf("expected 'test_version', got %s", cfg.GetVersionTable())
	}
}

func TestConfigValidate(t *testing.T) {
	testCases := map[string]struct {
		cfg       db.Config
		expectErr error
	}{
		"empty host": {
			cfg: db.Config{
				Host:     "",
				Port:     "5432",
				Database: "filament",
				Username: "filament",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_HOST", ""),
		},
		"loc@lhost": {
			cfg: db.Config{
				Host:     "loc@lhost",
				Port:     "5432",
				Database: "filament",
				Username: "filament",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_HOST", "loc@lhost"),
		},
		"local host": {
			cfg: db.Config{
				Host:     "local host",
				Port:     "5432",
				Database: "filament",
				Username: "filament",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_HOST", "local host"),
		},
		"empty user": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "5432",
				Database: "filament",
				Username: "",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_USER", ""),
		},
		"user!name": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "5432",
				Database: "filament",
				Username: "user!name",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_USER", "user!name"),
		},
		"user name": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "5432",
				Database: "filament",
				Username: "user name",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_USER", "user name"),
		},
		"empty db": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "5432",
				Database: "",
				Username: "filament",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_DATABASE", ""),
		},
		"data!base": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "5432",
				Database: "data!base",
				Username: "user!name",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_DATABASE", "data!base"),
		},
		"data base": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "5432",
				Database: "data base",
				Username: "user name",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_DATABASE", "data base"),
		},
		"empty password": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "5432",
				Database: "filament",
				Username: "filament",
				Password: "",
			},
			expectErr: db.NewConfigError("POSTGRES_PASSWORD", ""),
		},
		"empty port": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "",
				Database: "filament",
				Username: "filament",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_PORT", ""),
		},
		"port with letters": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "54a2",
				Database: "filament",
				Username: "filament",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_PORT", "54a2"),
		},
		"port with spaces": {
			cfg: db.Config{
				Host:     "localhost",
				Port:     "54 32",
				Database: "filament",
				Username: "filament",
				Password: "filament",
			},
			expectErr: db.NewConfigError("POSTGRES_PORT", "54 32"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if err == nil && tc.expectErr != nil {
				t.Errorf("expected error %v, got nil", tc.expectErr)
			} else if err != nil && tc.expectErr == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err != nil && tc.expectErr != nil && err.Error() != tc.expectErr.Error() {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}
		})
	}
}
