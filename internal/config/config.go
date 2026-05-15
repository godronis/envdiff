package config

import (
	"errors"
	"fmt"

	"github.com/yourusername/envdiff/internal/sorter"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Config holds all runtime configuration for an envdiff run.
type Config struct {
	Files          []string
	Format         Format
	Prefix         string
	ExcludeKeys    []string
	OnlyMissing    bool
	OnlyMismatched bool
	SortBy         sorter.SortBy
}

// Validate checks that the Config is well-formed and applies sensible defaults.
func Validate(c *Config) error {
	if len(c.Files) < 2 {
		return errors.New("at least two .env files must be provided")
	}

	for _, f := range c.Files {
		if f == "" {
			return errors.New("file path must not be empty")
		}
	}

	if c.Format == "" {
		c.Format = FormatText
	}

	switch c.Format {
	case FormatText, FormatJSON:
		// valid
	default:
		return fmt.Errorf("unsupported format %q: must be \"text\" or \"json\"", c.Format)
	}

	if c.SortBy == "" {
		c.SortBy = sorter.SortByKey
	}

	switch c.SortBy {
	case sorter.SortByKey, sorter.SortByStatus, sorter.SortByFile:
		// valid
	default:
		return fmt.Errorf("unsupported sort-by value %q: must be \"key\", \"status\", or \"file\"", c.SortBy)
	}

	return nil
}
