package merger_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/merger"
)

func TestMerge_NoConflicts(t *testing.T) {
	envs := []map[string]string{
		{"APP_HOST": "localhost"},
		{"APP_PORT": "8080"},
	}
	res := merger.Merge([]string{"a.env", "b.env"}, envs, merger.FirstWins)
	if res.Merged["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %s", res.Merged["APP_HOST"])
	}
	if res.Merged["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT=8080, got %s", res.Merged["APP_PORT"])
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_FirstWins(t *testing.T) {
	envs := []map[string]string{
		{"DB_PASS": "secret1"},
		{"DB_PASS": "secret2"},
	}
	res := merger.Merge([]string{"a.env", "b.env"}, envs, merger.FirstWins)
	if res.Merged["DB_PASS"] != "secret1" {
		t.Errorf("expected DB_PASS=secret1, got %s", res.Merged["DB_PASS"])
	}
	if res.Sources["DB_PASS"] != "a.env" {
		t.Errorf("expected source a.env, got %s", res.Sources["DB_PASS"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "DB_PASS" {
		t.Errorf("expected conflict on DB_PASS, got %v", res.Conflicts)
	}
}

func TestMerge_LastWins(t *testing.T) {
	envs := []map[string]string{
		{"DB_PASS": "secret1"},
		{"DB_PASS": "secret2"},
	}
	res := merger.Merge([]string{"a.env", "b.env"}, envs, merger.LastWins)
	if res.Merged["DB_PASS"] != "secret2" {
		t.Errorf("expected DB_PASS=secret2, got %s", res.Merged["DB_PASS"])
	}
	if res.Sources["DB_PASS"] != "b.env" {
		t.Errorf("expected source b.env, got %s", res.Sources["DB_PASS"])
	}
}

func TestMerge_ConflictsSorted(t *testing.T) {
	envs := []map[string]string{
		{"Z_KEY": "1", "A_KEY": "1", "M_KEY": "1"},
		{"Z_KEY": "2", "A_KEY": "2", "M_KEY": "2"},
	}
	res := merger.Merge([]string{"a.env", "b.env"}, envs, merger.FirstWins)
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, k := range expected {
		if res.Conflicts[i] != k {
			t.Errorf("expected conflicts[%d]=%s, got %s", i, k, res.Conflicts[i])
		}
	}
}

func TestMerge_EmptyInputs(t *testing.T) {
	res := merger.Merge([]string{}, []map[string]string{}, merger.FirstWins)
	if len(res.Merged) != 0 {
		t.Errorf("expected empty merged map")
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}
