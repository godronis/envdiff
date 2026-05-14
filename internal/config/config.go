// Package config handles loading and validating CLI configuration
// for the envdiff tool, including file paths, output format, and filter options.
package config

import (
	"errors"
	"fmt"
	"strings"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Config holds all configuration options for an envdiff run.
type Config struct {
	Files         []string
	OutputFormat  Format
	PrefixFilter  string
	ExcludeKeys   []string
	OnlyMissing   bool
	OnlyMismatched bool
}

// Validate checks that the Config is valid and returns an error if not.
func Validate(cfg *Config) error {
	if len(cfg.Files) < 2 {
		return errors.New("at least two .env files must be provided")
	}

	for _, f := range cfg.Files {
		if strings.TrimSpace(f) == "" {
			return errors.New("file path must not be empty")
		}
	}

	switch cfg.OutputFormat {
	case FormatText, FormatJSON:
		// valid
	case "":
		cfg.OutputFormat = FormatText
	default:
		return fmt.Errorf("unsupported output format: %q (use \"text\" or \"json\")", cfg.OutputFormat)
	}

	if cfg.OnlyMissing && cfg.OnlyMismatched {
		return errors.New("--only-missing and --only-mismatched are mutually exclusive")
	}

	return nil
}
