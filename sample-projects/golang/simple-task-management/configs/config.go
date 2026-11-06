package configs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config represents runtime configuration parameters.
type Config struct {
	DBDriver    string `yaml:"db_driver"`
	DatabaseURL string `yaml:"database_url"`
	APIHost     string `yaml:"api_host"`
	APIPort     int    `yaml:"api_port"`
	LogLevel    string `yaml:"log_level"`
}

// Default returns a Config populated with default values.
func Default() Config {
	return Config{
		DBDriver:    "sqlite",
		DatabaseURL: "tasks.db",
		APIHost:     "0.0.0.0",
		APIPort:     8000,
		LogLevel:    "info",
	}
}

// Load reads base.yaml, optional environment file, and overrides from environment variables.
func Load(dir, env string) (Config, error) {
	cfg := Default()

	if dir == "" {
		dir = "configs"
	}

	files := []string{filepath.Join(dir, "base.yaml")}
	if env != "" {
		files = append(files, filepath.Join(dir, env, "app.yaml"))
	}

	for _, f := range files {
		b, err := os.ReadFile(f)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return Config{}, fmt.Errorf("read %s: %w", f, err)
		}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse %s: %w", f, err)
		}
	}

	// Environment variable overrides
	if v, ok := os.LookupEnv("DB_DRIVER"); ok && v != "" {
		cfg.DBDriver = v
	}
	if v, ok := os.LookupEnv("DATABASE_URL"); ok && v != "" {
		cfg.DatabaseURL = v
	}
	if v, ok := os.LookupEnv("API_HOST"); ok && v != "" {
		cfg.APIHost = v
	}
	if v, ok := os.LookupEnv("API_PORT"); ok && v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.APIPort = p
		}
	}
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok && v != "" {
		cfg.LogLevel = v
	}

	return cfg, nil
}

