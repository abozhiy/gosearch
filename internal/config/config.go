package config

import (
	"fmt"
	"os"
	"strings"

	"gosearch/internal/store"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database store.Database `yaml:"database"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	settings, err := loadSettings()
	if err != nil {
		return nil, fmt.Errorf("settings load failed: %w", err)
	}

	cfg, err := parseYAML(os.ExpandEnv(string(settings)))
	if err != nil {
		return nil, fmt.Errorf("config parse failed > %w", err)
	}

	return cfg, nil
}

func loadSettings() ([]byte, error) {
	path := os.Getenv("SETTINGS_PATH")

	if strings.TrimSpace(path) == "" || !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("invalid settings path: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file at %s > %w", path, err)
	}

	return data, nil
}

func parseYAML(raw string) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		return nil, fmt.Errorf("yaml parse error > %w", err)
	}
	return &cfg, nil
}
