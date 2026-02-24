// Package db implements the database connection and related functionality.
// This includes configuration creation/loading, connection URI generation, and migration handling.
package db

import (
	"fmt"
	"log/slog"
	"regexp"

	env "github.com/caarlos0/env/v10"
)

type Secret string

func (Secret) LogValue() slog.Value {
	return slog.StringValue("REDACTED_SECRET")
}

type Config struct {
	Host         string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port         string `env:"POSTGRES_PORT" envDefault:"5432"`
	Database     string `env:"POSTGRES_DATABASE" envDefault:"filament"` // Database name
	Username     string `env:"POSTGRES_USER" envDefault:"filament"`     // Username
	Password     Secret `env:"POSTGRES_PASSWORD" envDefault:"filament"` // Password
	VersionTable string `env:"VERSION_TABLE" envDefault:"db_version"`   // Migrations version table
	DemoData     bool   `env:"DEMO_DATA" envDefault:"false"`            // Whether to include demo data in the database
}

func (c Config) GetURI() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Username, c.Password, c.Host, c.Port, c.Database)
}

func (c Config) GetRedactedURI() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Username, c.Password.LogValue(), c.Host, c.Port, c.Database)
}

func (c Config) GetHost() string {
	return c.Host
}

func (c Config) GetPort() string {
	return c.Port
}

func (c Config) GetDatabase() string {
	return c.Database
}

func (c Config) GetPassword() string {
	if c.Password == "" {
		return "UNSET"
	}

	return "REDACTED_SECRET"
}

func (c Config) GetUsername() string {
	return c.Username
}

func (c Config) GetVersionTable() string {
	return c.VersionTable
}

func (c Config) GetIncDemo() bool {
	return c.DemoData
}

func LoadConfig() (cfg Config, fault error) {
	c := Config{}
	if err := env.Parse(&c); err != nil {
		return Config{}, fmt.Errorf("could not load config: %w", err)
	}

	if err := c.Validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return c, nil
}

type ConfigError struct {
	field string
	value string
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("invalid value %q for field %q", e.value, e.field)
}

func NewConfigError(field, value string) error {
	return ConfigError{field: field, value: value}
}

var (
	alphaNum = regexp.MustCompile(`^[a-zA-Z0-9\._-]+$`) // Host, Database, Username
	numeric  = regexp.MustCompile(`^[0-9]+$`)           // Port
)

func (c Config) Validate() error {
	if !alphaNum.MatchString(c.Host) {
		return ConfigError{field: "POSTGRES_HOST", value: c.Host}
	}

	if !numeric.MatchString(c.Port) {
		return ConfigError{field: "POSTGRES_PORT", value: c.Port}
	}

	if !alphaNum.MatchString(c.Database) {
		return ConfigError{field: "POSTGRES_DATABASE", value: c.Database}
	}

	if !alphaNum.MatchString(c.Username) {
		return ConfigError{field: "POSTGRES_USER", value: c.Username}
	}

	if c.Password == "" {
		return ConfigError{field: "POSTGRES_PASSWORD", value: ""}
	}

	return nil
}
