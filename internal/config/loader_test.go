package config

import (
	"testing"
)

func TestLoadFromArgs_BasicFiles(t *testing.T) {
	cfg, err := LoadFromArgs([]string{".env.dev", ".env.prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.Files))
	}
	if cfg.OutputFormat != FormatText {
		t.Errorf("expected default format text, got %q", cfg.OutputFormat)
	}
}

func TestLoadFromArgs_JSONFormat(t *testing.T) {
	cfg, err := LoadFromArgs([]string{"--format", "json", ".env.a", ".env.b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != FormatJSON {
		t.Errorf("expected json format, got %q", cfg.OutputFormat)
	}
}

func TestLoadFromArgs_PrefixFilter(t *testing.T) {
	cfg, err := LoadFromArgs([]string{"--prefix", "APP_", ".env.a", ".env.b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PrefixFilter != "APP_" {
		t.Errorf("expected prefix APP_, got %q", cfg.PrefixFilter)
	}
}

func TestLoadFromArgs_ExcludeKeys(t *testing.T) {
	cfg, err := LoadFromArgs([]string{"--exclude", "SECRET", "--exclude", "TOKEN", ".env.a", ".env.b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.ExcludeKeys) != 2 {
		t.Errorf("expected 2 excluded keys, got %d", len(cfg.ExcludeKeys))
	}
}

func TestLoadFromArgs_OnlyMissingFlag(t *testing.T) {
	cfg, err := LoadFromArgs([]string{"--only-missing", ".env.a", ".env.b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.OnlyMissing {
		t.Error("expected OnlyMissing to be true")
	}
}

func TestLoadFromArgs_MutuallyExclusiveFlags(t *testing.T) {
	_, err := LoadFromArgs([]string{"--only-missing", "--only-mismatched", ".env.a", ".env.b"})
	if err == nil {
		t.Fatal("expected error for mutually exclusive flags")
	}
}

func TestLoadFromArgs_TooFewFiles(t *testing.T) {
	_, err := LoadFromArgs([]string{".env.only"})
	if err == nil {
		t.Fatal("expected error for fewer than two files")
	}
}
