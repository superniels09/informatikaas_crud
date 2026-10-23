package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type Config struct {
	Database *DatabaseConfig `yaml:"database"`
	JWT      *JWTConfig      `yaml:"jwt"`

	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func Load(name string) (*Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Database == nil {
		return nil, fmt.Errorf("database configuration missing")
	}

	return &cfg, nil
}
