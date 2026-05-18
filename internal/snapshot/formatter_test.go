package snapshot_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/snapshot"
)

func TestFormatText_NoChanges(t *testing.T) {
	d := snapshot.DiffResult{BaseLabel: "a", CurrentLabel: "b"}
	out := snapshot.FormatText(d)
	if !strings.Contains(out, "No changes detected") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
}

func TestFormatText_ShowsAddedRemovedChanged(t *testing.T) {
	d := snapshot.DiffResult{
		BaseLabel:    "v1",
		CurrentLabel: "v2",
		Added:        []string{"NEW_KEY"},
		Removed:      []string{"OLD_KEY"},
		Changed:      []string{"CHANGED_KEY"},
	}
	out := snapshot.FormatText(d)
	for _, want := range []string{"+ NEW_KEY", "- OLD_KEY", "~ CHANGED_KEY", "v1", "v2"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestFormatJSON_ValidStructure(t *testing.T) {
	d := snapshot.DiffResult{
		BaseLabel:    "base",
		CurrentLabel: "current",
		Added:        []string{"A"},
		Removed:      []string{"B"},
		Changed:      []string{"C"},
	}
	out, err := snapshot.FormatJSON(d)
	if err != nil {
		t.Fatalf("FormatJSON error: %v", err)
	}
	for _, want := range []string{`"base"`, `"current"`, `"A"`, `"B"`, `"C"`, `"added"`, `"removed"`, `"changed"`} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output, got:\n%s", want, out)
		}
	}
}

func TestFormatJSON_EmptySlicesNotNull(t *testing.T) {
	d := snapshot.DiffResult{BaseLabel: "x", CurrentLabel: "y"}
	out, err := snapshot.FormatJSON(d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "null") {
		t.Errorf("expected no null in output, got:\n%s", out)
	}
}
