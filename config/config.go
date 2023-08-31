package config

import (
	"avito-segments/pkg/db/postgres"
	"avito-segments/pkg/httpserver"
	"avito-segments/pkg/logger"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Logger     logger.Config     `yaml:"log"`
	Storage    postgres.Config   `yaml:"storage"`
	HTTPServer httpserver.Config `yaml:"http_server"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH is not set")
	}

	// check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %s", err)
	}

	return cfg, nil
}
