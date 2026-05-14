package config

import (
	"flag"
	"strings"
)

// excludeList is a helper type for collecting repeated --exclude flags.
type excludeList []string

func (e *excludeList) String() string {
	return strings.Join(*e, ",")
}

func (e *excludeList) Set(value string) error {
	*e = append(*e, value)
	return nil
}

// LoadFromArgs parses os.Args-style arguments and returns a Config.
// It does not call os.Exit; callers should handle the returned error.
func LoadFromArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)

	format := fs.String("format", "text", "Output format: text or json")
	prefix := fs.String("prefix", "", "Only compare keys with this prefix")
	onlyMissing := fs.Bool("only-missing", false, "Report only missing keys")
	onlyMismatched := fs.Bool("only-mismatched", false, "Report only mismatched values")

	var exclude excludeList
	fs.Var(&exclude, "exclude", "Exclude a specific key (repeatable)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		Files:          fs.Args(),
		OutputFormat:   Format(*format),
		PrefixFilter:   *prefix,
		ExcludeKeys:    []string(exclude),
		OnlyMissing:    *onlyMissing,
		OnlyMismatched: *onlyMismatched,
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
