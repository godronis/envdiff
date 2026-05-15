package linter_test

import (
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func TestLint_NoWarnings(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"APP_PORT":     "8080",
	}
	warnings := linter.Lint("prod.env", env)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %d", len(warnings))
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"app_secret": "abc123",
	}
	warnings := linter.Lint("dev.env", env)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Severity != linter.SeverityWarn {
		t.Errorf("expected warn severity, got %s", warnings[0].Severity)
	}
	if warnings[0].Key != "app_secret" {
		t.Errorf("unexpected key: %s", warnings[0].Key)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{
		"API_KEY": "",
	}
	warnings := linter.Lint("staging.env", env)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Severity != linter.SeverityWarn {
		t.Errorf("expected warn severity, got %s", warnings[0].Severity)
	}
}

func TestLint_SuspiciousPlaceholder(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "<your-password-here>",
	}
	warnings := linter.Lint("example.env", env)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Severity != linter.SeverityError {
		t.Errorf("expected error severity, got %s", warnings[0].Severity)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"lower_key":  "",
		"PLACEHOLDER": "<fill-me>",
		"GOOD_KEY":   "good-value",
	}
	warnings := linter.Lint("test.env", env)
	// lower_key triggers empty + lowercase = 2 warnings; PLACEHOLDER triggers 1
	if len(warnings) < 3 {
		t.Errorf("expected at least 3 warnings, got %d", len(warnings))
	}
}

func TestLint_FileNameAttached(t *testing.T) {
	env := map[string]string{
		"EMPTY_VAL": "",
	}
	warnings := linter.Lint("myfile.env", env)
	if len(warnings) == 0 {
		t.Fatal("expected at least one warning")
	}
	if warnings[0].File != "myfile.env" {
		t.Errorf("expected file 'myfile.env', got %q", warnings[0].File)
	}
}
