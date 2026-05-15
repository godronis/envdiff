// Package linter provides basic linting checks for parsed .env file entries,
// such as detecting suspicious values, empty values, or keys with unusual formatting.
package linter

import (
	"fmt"
	"strings"
)

// Severity represents the severity level of a lint warning.
type Severity string

const (
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

// Warning represents a single lint warning for a specific key.
type Warning struct {
	Key      string
	File     string
	Message  string
	Severity Severity
}

// Lint checks the provided env map (from a single file) for common issues
// and returns a slice of warnings.
func Lint(file string, env map[string]string) []Warning {
	var warnings []Warning

	for key, value := range env {
		if w := checkKeyFormat(file, key); w != nil {
			warnings = append(warnings, *w)
		}
		if w := checkEmptyValue(file, key, value); w != nil {
			warnings = append(warnings, *w)
		}
		if w := checkSuspiciousValue(file, key, value); w != nil {
			warnings = append(warnings, *w)
		}
	}

	return warnings
}

// checkKeyFormat warns if a key contains lowercase letters (non-standard).
func checkKeyFormat(file, key string) *Warning {
	if key != strings.ToUpper(key) {
		return &Warning{
			Key:      key,
			File:     file,
			Message:  fmt.Sprintf("key %q is not uppercase; env keys are conventionally uppercase", key),
			Severity: SeverityWarn,
		}
	}
	return nil
}

// checkEmptyValue warns if a value is empty.
func checkEmptyValue(file, key, value string) *Warning {
	if strings.TrimSpace(value) == "" {
		return &Warning{
			Key:      key,
			File:     file,
			Message:  fmt.Sprintf("key %q has an empty value", key),
			Severity: SeverityWarn,
		}
	}
	return nil
}

// checkSuspiciousValue warns if a value looks like an unresolved placeholder.
func checkSuspiciousValue(file, key, value string) *Warning {
	if strings.Contains(value, "<") && strings.Contains(value, ">") {
		return &Warning{
			Key:      key,
			File:     file,
			Message:  fmt.Sprintf("key %q appears to contain an unresolved placeholder: %q", key, value),
			Severity: SeverityError,
		}
	}
	return nil
}
