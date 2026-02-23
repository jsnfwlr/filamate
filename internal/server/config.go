package server

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	validator "github.com/go-playground/validator/v10"
)

type EnvConfig struct {
	StaticType string `env:"STATIC_TYPE" envDefault:"embedded" validate:"oneof=embedded directory"` // embedded or directory
	Host       string `env:"API_HOST" envDefault:"0.0.0.0" validate:"hostname|ip"`
	Port       string `env:"API_PORT" envDefault:"9766" validate:"numeric"`
}

type Config struct {
	staticType string
	host       string
	port       string
}

type Configuration interface {
	StaticType() string
	Host() string
	Port() string
}

func LoadConfig() (cfg Config, fault error) {
	e := EnvConfig{}
	_ = env.Parse(&e)

	err := e.Validate()
	if err != nil {
		return Config{}, err
	}

	c := Config{
		staticType: e.StaticType,
		host:       e.Host,
		port:       e.Port,
	}

	return c, nil
}

func (e EnvConfig) Validate() error {
	validation := validator.New()
	if err := validation.Struct(e); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

func (c Config) StaticType() string {
	return c.staticType
}

func (c Config) Host() string {
	return c.host
}

func (c Config) Port() string {
	return c.port
}
