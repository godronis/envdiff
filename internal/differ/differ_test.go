package differ_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/differ"
)

func makeResults() []comparator.Result {
	return []comparator.Result{
		{
			Key:      "APP_NAME",
			Status:   "ok",
			Values:   map[string]string{"a.env": "myapp", "b.env": "myapp"},
		},
		{
			Key:      "DB_HOST",
			Status:   "mismatch",
			Values:   map[string]string{"a.env": "localhost", "b.env": "db.prod"},
		},
		{
			Key:      "SECRET_KEY",
			Status:   "missing",
			Values:   map[string]string{"a.env": "abc123"},
		},
		{
			Key:      "NEW_FLAG",
			Status:   "missing",
			Values:   map[string]string{"b.env": "true"},
		},
	}
}

func TestGenerate_StatusUnchanged(t *testing.T) {
	results := makeResults()
	d := differ.Generate("a.env", "b.env", results)

	if d.Lines[0].Status != "unchanged" {
		t.Errorf("expected unchanged, got %s", d.Lines[0].Status)
	}
}

func TestGenerate_StatusChanged(t *testing.T) {
	results := makeResults()
	d := differ.Generate("a.env", "b.env", results)

	if d.Lines[1].Status != "changed" {
		t.Errorf("expected changed, got %s", d.Lines[1].Status)
	}
	if d.Lines[1].Old != "localhost" || d.Lines[1].New != "db.prod" {
		t.Errorf("unexpected values: old=%s new=%s", d.Lines[1].Old, d.Lines[1].New)
	}
}

func TestGenerate_StatusRemoved(t *testing.T) {
	results := makeResults()
	d := differ.Generate("a.env", "b.env", results)

	if d.Lines[2].Status != "removed" {
		t.Errorf("expected removed, got %s", d.Lines[2].Status)
	}
}

func TestGenerate_StatusAdded(t *testing.T) {
	results := makeResults()
	d := differ.Generate("a.env", "b.env", results)

	if d.Lines[3].Status != "added" {
		t.Errorf("expected added, got %s", d.Lines[3].Status)
	}
}

func TestFormatUnified_ContainsHeaders(t *testing.T) {
	results := makeResults()
	d := differ.Generate("a.env", "b.env", results)
	out := differ.FormatUnified(d)

	if !strings.Contains(out, "--- a.env") {
		t.Error("missing --- header")
	}
	if !strings.Contains(out, "+++ b.env") {
		t.Error("missing +++ header")
	}
}

func TestFormatUnified_ContainsDiffMarkers(t *testing.T) {
	results := makeResults()
	d := differ.Generate("a.env", "b.env", results)
	out := differ.FormatUnified(d)

	if !strings.Contains(out, "+ NEW_FLAG") {
		t.Error("expected + marker for added key")
	}
	if !strings.Contains(out, "- SECRET_KEY") {
		t.Error("expected - marker for removed key")
	}
	if !strings.Contains(out, "~ DB_HOST") {
		t.Error("expected ~ marker for changed key")
	}
}
