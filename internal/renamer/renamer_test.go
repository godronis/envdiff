package renamer_test

import (
	"testing"

	"github.com/user/envdiff/internal/renamer"
)

func TestApply_BasicRename(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value1", "KEEP": "value2"}
	mappings := []renamer.Mapping{{From: "OLD_KEY", To: "NEW_KEY"}}

	res, err := renamer.Apply(env, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if res.Env["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", res.Env["NEW_KEY"])
	}
	if res.Env["KEEP"] != "value2" {
		t.Error("expected KEEP to be preserved")
	}
	if len(res.Renamed) != 1 || res.Renamed[0].From != "OLD_KEY" {
		t.Error("expected one renamed entry")
	}
}

func TestApply_KeyNotFound(t *testing.T) {
	env := map[string]string{"EXISTING": "val"}
	mappings := []renamer.Mapping{{From: "MISSING", To: "TARGET"}}

	res, err := renamer.Apply(env, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.NotFound) != 1 || res.NotFound[0] != "MISSING" {
		t.Errorf("expected MISSING in NotFound, got %v", res.NotFound)
	}
	if len(res.Renamed) != 0 {
		t.Error("expected no renamed entries")
	}
}

func TestApply_ConflictReturnsError(t *testing.T) {
	env := map[string]string{"SRC": "a", "DST": "b"}
	mappings := []renamer.Mapping{{From: "SRC", To: "DST"}}

	_, err := renamer.Apply(env, mappings)
	if err == nil {
		t.Fatal("expected error for conflicting target key")
	}
}

func TestApply_NilEnv(t *testing.T) {
	res, err := renamer.Apply(nil, []renamer.Mapping{{From: "A", To: "B"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Error("expected empty env for nil input")
	}
}

func TestApply_MultipleRenames(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	mappings := []renamer.Mapping{
		{From: "A", To: "ALPHA"},
		{From: "B", To: "BETA"},
	}

	res, err := renamer.Apply(env, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 2 {
		t.Errorf("expected 2 renamed, got %d", len(res.Renamed))
	}
	if res.Env["ALPHA"] != "1" || res.Env["BETA"] != "2" || res.Env["C"] != "3" {
		t.Error("unexpected env state after multiple renames")
	}
}
