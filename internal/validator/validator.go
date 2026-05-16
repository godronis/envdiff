// Package validator checks .env file entries for common
// security and correctness issues beyond basic linting.
package validator

import (
	"fmt"
	"strings"
)

// Issue represents a single validation problem found in an env entry.
type Issue struct {
	Key      string
	File     string
	Severity string // "error" or "warning"
	Message  string
}

// Result holds all issues found after validating one or more env maps.
type Result struct {
	Issues []Issue
}

// HasErrors returns true if any issue has severity "error".
func (r Result) HasErrors() bool {
	for _, i := range r.Issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}

// Validate runs all checks against the provided env maps (filename -> key/value pairs).
func Validate(envs map[string]map[string]string) Result {
	var issues []Issue
	for file, entries := range envs {
		for key, value := range entries {
			issues = append(issues, checkSecretExposed(file, key, value)...)
			issues = append(issues, checkURLScheme(file, key, value)...)
			issues = append(issues, checkDuplicateSemantics(file, key, value)...)
		}
	}
	return Result{Issues: issues}
}

var secretKeywords = []string{"PASSWORD", "SECRET", "TOKEN", "API_KEY", "PRIVATE_KEY"}

func checkSecretExposed(file, key, value string) []Issue {
	if value == "" {
		return nil
	}
	upper := strings.ToUpper(key)
	for _, kw := range secretKeywords {
		if strings.Contains(upper, kw) {
			if !looksRedacted(value) {
				return []Issue{{
					Key:      key,
					File:     file,
					Severity: "warning",
					Message:  fmt.Sprintf("key %q may contain a sensitive value", key),
				}}
			}
		}
	}
	return nil
}

func looksRedacted(value string) bool {
	v := strings.ToUpper(value)
	return v == "CHANGEME" || v == "REDACTED" || v == "***" || strings.HasPrefix(v, "<") && strings.HasSuffix(v, ">")
}

func checkURLScheme(file, key, value string) []Issue {
	upper := strings.ToUpper(key)
	if !strings.Contains(upper, "URL") && !strings.Contains(upper, "HOST") {
		return nil
	}
	if strings.HasPrefix(value, "http://") {
		return []Issue{{
			Key:      key,
			File:     file,
			Severity: "warning",
			Message:  fmt.Sprintf("key %q uses insecure http:// scheme", key),
		}}
	}
	return nil
}

func checkDuplicateSemantics(file, key, value string) []Issue {
	if strings.Contains(value, "${{") || strings.Contains(value, "%%") {
		return []Issue{{
			Key:      key,
			File:     file,
			Severity: "error",
			Message:  fmt.Sprintf("key %q contains unresolved template expression in value", key),
		}}
	}
	return nil
}
