package validator_test

import (
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func TestValidate_NoIssues(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {
			"APP_NAME": "myapp",
			"PORT":     "8080",
		},
	}
	result := validator.Validate(envs)
	if len(result.Issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(result.Issues))
	}
}

func TestValidate_SecretExposed(t *testing.T) {
	envs := map[string]map[string]string{
		".env.prod": {
			"DB_PASSWORD": "supersecret123",
		},
	}
	result := validator.Validate(envs)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %s", result.Issues[0].Severity)
	}
}

func TestValidate_SecretRedacted_NoIssue(t *testing.T) {
	envs := map[string]map[string]string{
		".env.example": {
			"API_KEY": "<your-api-key>",
		},
	}
	result := validator.Validate(envs)
	if len(result.Issues) != 0 {
		t.Fatalf("expected no issues for redacted secret, got %d", len(result.Issues))
	}
}

func TestValidate_InsecureURL(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {
			"DATABASE_URL": "http://localhost:5432/db",
		},
	}
	result := validator.Validate(envs)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", result.Issues[0].Severity)
	}
}

func TestValidate_UnresolvedTemplate(t *testing.T) {
	envs := map[string]map[string]string{
		".env": {
			"CALLBACK_URL": "https://example.com/${{BASE_PATH}}/callback",
		},
	}
	result := validator.Validate(envs)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", result.Issues[0].Severity)
	}
}

func TestResult_HasErrors_True(t *testing.T) {
	r := validator.Result{
		Issues: []validator.Issue{
			{Key: "X", File: "f", Severity: "error", Message: "bad"},
		},
	}
	if !r.HasErrors() {
		t.Error("expected HasErrors to return true")
	}
}

func TestResult_HasErrors_False(t *testing.T) {
	r := validator.Result{
		Issues: []validator.Issue{
			{Key: "X", File: "f", Severity: "warning", Message: "maybe"},
		},
	}
	if r.HasErrors() {
		t.Error("expected HasErrors to return false for warnings only")
	}
}
