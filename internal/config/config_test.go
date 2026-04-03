package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvConfigDir, dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKey != "" {
		t.Errorf("expected empty api key, got %q", cfg.APIKey)
	}
	if cfg.DefaultFormat != "" {
		t.Errorf("expected empty default format, got %q", cfg.DefaultFormat)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvConfigDir, dir)

	cfg := &Config{
		APIKey:        "test-key-12345",
		DefaultFormat: "json",
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %o", info.Mode().Perm())
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.APIKey != "test-key-12345" {
		t.Errorf("expected api key 'test-key-12345', got %q", loaded.APIKey)
	}
	if loaded.DefaultFormat != "json" {
		t.Errorf("expected default format 'json', got %q", loaded.DefaultFormat)
	}
}

func TestResolveAPIKey_EnvOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvConfigDir, dir)
	t.Setenv(EnvAPIKey, "env-key")

	cfg := &Config{APIKey: "config-key"}
	if err := Save(cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	got := ResolveAPIKey("")
	if got != "env-key" {
		t.Errorf("expected 'env-key', got %q", got)
	}
}

func TestResolveAPIKey_FlagOverridesEnv(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-key")

	got := ResolveAPIKey("flag-key")
	if got != "flag-key" {
		t.Errorf("expected 'flag-key', got %q", got)
	}
}

func TestResolveAPIKey_FallbackToConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvConfigDir, dir)
	t.Setenv(EnvAPIKey, "")

	cfg := &Config{APIKey: "config-key"}
	if err := Save(cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	got := ResolveAPIKey("")
	if got != "config-key" {
		t.Errorf("expected 'config-key', got %q", got)
	}
}
