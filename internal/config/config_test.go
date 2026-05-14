package config

import (
	"testing"
)

func TestValidate_ValidConfig(t *testing.T) {
	cfg := &Config{
		Files:        []string{".env.dev", ".env.prod"},
		OutputFormat: FormatText,
	}
	if err := Validate(cfg); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_DefaultsToText(t *testing.T) {
	cfg := &Config{
		Files: []string{".env.dev", ".env.prod"},
	}
	if err := Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != FormatText {
		t.Errorf("expected default format %q, got %q", FormatText, cfg.OutputFormat)
	}
}

func TestValidate_TooFewFiles(t *testing.T) {
	cfg := &Config{Files: []string{".env.dev"}}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for fewer than two files")
	}
}

func TestValidate_EmptyFilePath(t *testing.T) {
	cfg := &Config{Files: []string{".env.dev", "   "}}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for empty file path")
	}
}

func TestValidate_UnsupportedFormat(t *testing.T) {
	cfg := &Config{
		Files:        []string{".env.dev", ".env.prod"},
		OutputFormat: "yaml",
	}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestValidate_MutuallyExclusiveFlags(t *testing.T) {
	cfg := &Config{
		Files:          []string{".env.dev", ".env.prod"},
		OnlyMissing:    true,
		OnlyMismatched: true,
	}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for mutually exclusive flags")
	}
}

func TestValidate_JSONFormat(t *testing.T) {
	cfg := &Config{
		Files:        []string{".env.dev", ".env.prod"},
		OutputFormat: FormatJSON,
	}
	if err := Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
