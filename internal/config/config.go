package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey        string `yaml:"api_key"`
	DefaultFormat string `yaml:"default_format,omitempty"`
}

func configDir() string {
	if dir := os.Getenv(EnvConfigDir); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "edinet")
}

func configPath() string {
	return filepath.Join(configDir(), "config.yaml")
}

func Load() (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Save(cfg *Config) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0600)
}

func ResolveAPIKey(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if env := os.Getenv(EnvAPIKey); env != "" {
		return env
	}
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.APIKey
}

func ConfigDir() string {
	return configDir()
}
