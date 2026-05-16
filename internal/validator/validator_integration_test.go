package validator_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/validator"
)

const sampleEnvContent = `
APP_NAME=myservice
DB_PASSWORD=hunter2
DATABASE_URL=http://db.internal:5432/prod
REDIS_TOKEN=<redacted>
CALLBACK=${{UNRESOLVED}}
`

func TestValidate_Integration_ParsedEnv(t *testing.T) {
	entries, err := parser.ParseString(strings.TrimSpace(sampleEnvContent))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	envs := map[string]map[string]string{
		".env.prod": entries,
	}

	result := validator.Validate(envs)

	if len(result.Issues) == 0 {
		t.Fatal("expected validation issues but got none")
	}

	severities := map[string]int{}
	for _, issue := range result.Issues {
		severities[issue.Severity]++
		if issue.Key == "" {
			t.Errorf("issue missing key: %+v", issue)
		}
		if issue.File == "" {
			t.Errorf("issue missing file: %+v", issue)
		}
	}

	if severities["error"] < 1 {
		t.Errorf("expected at least 1 error-level issue, got %d", severities["error"])
	}
	if severities["warning"] < 2 {
		t.Errorf("expected at least 2 warning-level issues, got %d", severities["warning"])
	}
}

func TestValidate_Integration_CleanEnv(t *testing.T) {
	clean := `APP_NAME=myservice
PORT=8080
LOG_LEVEL=info
`
	entries, err := parser.ParseString(clean)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	envs := map[string]map[string]string{
		".env": entries,
	}

	result := validator.Validate(envs)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues for clean env, got: %+v", result.Issues)
	}
	if result.HasErrors() {
		t.Error("expected HasErrors() == false for clean env")
	}
}
