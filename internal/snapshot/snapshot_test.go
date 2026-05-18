package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/snapshot"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	s := snapshot.Snapshot{
		Label:    "v1",
		FilePath: ".env",
		Env:      map[string]string{"KEY": "value", "PORT": "8080"},
	}
	if err := snapshot.Save(s, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Label != s.Label {
		t.Errorf("expected label %q, got %q", s.Label, loaded.Label)
	}
	if loaded.Env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", loaded.Env["KEY"])
	}
}

func TestSave_EmptyLabelReturnsError(t *testing.T) {
	dir := t.TempDir()
	s := snapshot.Snapshot{Env: map[string]string{"K": "v"}}
	if err := snapshot.Save(s, filepath.Join(dir, "snap.json")); err == nil {
		t.Error("expected error for empty label")
	}
}

func TestSave_EmptyEnvReturnsError(t *testing.T) {
	dir := t.TempDir()
	s := snapshot.Snapshot{Label: "v1"}
	if err := snapshot.Save(s, filepath.Join(dir, "snap.json")); err == nil {
		t.Error("expected error for empty env")
	}
}

func TestLoad_InvalidPathReturnsError(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSONReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0644)
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiff_AddedRemovedChanged(t *testing.T) {
	base := snapshot.Snapshot{
		Label: "base",
		Env:   map[string]string{"A": "1", "B": "2", "C": "3"},
	}
	current := snapshot.Snapshot{
		Label: "current",
		Env:   map[string]string{"A": "1", "B": "changed", "D": "new"},
	}
	result := snapshot.Diff(base, current)
	if len(result.Added) != 1 || result.Added[0] != "D" {
		t.Errorf("expected Added=[D], got %v", result.Added)
	}
	if len(result.Removed) != 1 || result.Removed[0] != "C" {
		t.Errorf("expected Removed=[C], got %v", result.Removed)
	}
	if len(result.Changed) != 1 || result.Changed[0] != "B" {
		t.Errorf("expected Changed=[B], got %v", result.Changed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	env := map[string]string{"X": "1"}
	base := snapshot.Snapshot{Label: "a", Env: env}
	current := snapshot.Snapshot{Label: "b", Env: map[string]string{"X": "1"}}
	result := snapshot.Diff(base, current)
	if len(result.Added)+len(result.Removed)+len(result.Changed) != 0 {
		t.Errorf("expected no diff, got %+v", result)
	}
}
