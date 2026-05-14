// Package config provides configuration loading and validation for envdiff.
//
// It defines the Config struct that captures all user-supplied options,
// including file paths, output format, prefix filters, excluded keys,
// and reporting mode flags.
//
// Usage:
//
//	cfg, err := config.LoadFromArgs(os.Args[1:])
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		os.Exit(1)
//	}
//
// Validation is performed automatically by LoadFromArgs, but can also
// be invoked directly via config.Validate for custom loading scenarios.
package config
