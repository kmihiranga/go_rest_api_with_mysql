package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse() *Config {
	return &Config{
		AppConfig: parseAppConfig(),
		DBConfig:  parseDBConfig(),
	}
}

// parse application configurations
func parseAppConfig() *AppConfig {
	content := read("app.yaml")

	cfg := &AppConfig{}

	err := yaml.Unmarshal(content, &cfg)

	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return cfg
}

// parse database configurations
func parseDBConfig() *DBConfig {
	content := read("database.yaml")

	cfg := &DBConfig{}

	err := yaml.Unmarshal(content, &cfg)
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return cfg
}
