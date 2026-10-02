package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
)

type AppConfig struct {
	LogLevel               string `env:"LOG_LEVEL" envDefault:"debug"`
	PostmanApiKey          string `env:"POSTMAN_API_KEY" envDefault:""`
	PostmanWorkspaceID     string `env:"POSTMAN_WORKSPACE_ID" envDefault:""`
	PostmanCollectionID    string `env:"POSTMAN_COLLECTION_ID" envDefault:""`
	SwaggerAuthToken       string `env:"SWAGGER_AUTH_TOKEN" envDefault:""`
	SwaggerDomainPrefix    string `env:"SWAGGER_DOMAIN_PREFIX" envDefault:"/"`
	SwaggerOptionsHtmlPath string
	Oas2pmPath             string
}

func Load() (*AppConfig, error) {
	cfg := &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	cfg.SwaggerOptionsHtmlPath = "./swagger_options.html"
	if _, err := os.Stat(cfg.SwaggerOptionsHtmlPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("swagger options: %w", err)
	}

	cfg.Oas2pmPath = "./tools/oas2pm"
	if _, err := os.Stat(cfg.Oas2pmPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("oas2pm: %w", err)
	}

	return cfg, nil
}
