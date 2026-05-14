package comparator_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/comparator"
)

func TestCompare_NoMissingNoMismatch(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"DB_HOST": "localhost", "PORT": "8080"},
		"prod": {"DB_HOST": "localhost", "PORT": "8080"},
	}
	result := comparator.Compare(envs)
	if len(result.MissingIn) != 0 {
		t.Errorf("expected no missing keys, got %v", result.MissingIn)
	}
	if len(result.Mismatched) != 0 {
		t.Errorf("expected no mismatched keys, got %v", result.Mismatched)
	}
}

func TestCompare_MissingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"DB_HOST": "localhost", "SECRET": "abc"},
		"prod": {"DB_HOST": "localhost"},
	}
	result := comparator.Compare(envs)
	missing, ok := result.MissingIn["prod"]
	if !ok {
		t.Fatal("expected 'prod' to have missing keys")
	}
	if len(missing) != 1 || missing[0] != "SECRET" {
		t.Errorf("expected missing key SECRET in prod, got %v", missing)
	}
}

func TestCompare_MismatchedValue(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"DB_HOST": "localhost"},
		"prod": {"DB_HOST": "db.prod.example.com"},
	}
	result := comparator.Compare(envs)
	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatched key, got %d", len(result.Mismatched))
	}
	if result.Mismatched[0].Key != "DB_HOST" {
		t.Errorf("expected mismatched key DB_HOST, got %s", result.Mismatched[0].Key)
	}
}

func TestCompare_MultipleMissingAndMismatch(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":     {"APP_ENV": "development", "PORT": "3000", "SECRET": "devsecret"},
		"staging": {"APP_ENV": "staging", "PORT": "3000"},
		"prod":    {"APP_ENV": "production", "PORT": "80", "SECRET": "prodsecret"},
	}
	result := comparator.Compare(envs)

	// SECRET should be missing in staging
	if stagingMissing, ok := result.MissingIn["staging"]; !ok {
		t.Error("expected staging to have missing keys")
	} else if len(stagingMissing) != 1 || stagingMissing[0] != "SECRET" {
		t.Errorf("unexpected missing keys for staging: %v", stagingMissing)
	}

	// APP_ENV, PORT, SECRET should all be mismatched
	mismatchedKeys := make(map[string]bool)
	for _, m := range result.Mismatched {
		mismatchedKeys[m.Key] = true
	}
	for _, key := range []string{"APP_ENV", "PORT"} {
		if !mismatchedKeys[key] {
			t.Errorf("expected key %s to be mismatched", key)
		}
	}
}

func TestCompare_EmptyEnvs(t *testing.T) {
	result := comparator.Compare(map[string]map[string]string{})
	if len(result.MissingIn) != 0 || len(result.Mismatched) != 0 {
		t.Error("expected empty result for empty input")
	}
}
