package summary_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envdiff/internal/comparator"
	"github.com/yourusername/envdiff/internal/summary"
)

func makeResults() []comparator.Result {
	return []comparator.Result{
		{Key: "A", Status: comparator.StatusMatch},
		{Key: "B", Status: comparator.StatusMissing},
		{Key: "C", Status: comparator.StatusMismatch},
		{Key: "D", Status: comparator.StatusMatch},
		{Key: "E", Status: comparator.StatusMissing},
	}
}

func TestCompute_Counts(t *testing.T) {
	s := summary.Compute(makeResults())
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.Matched != 2 {
		t.Errorf("expected Matched=2, got %d", s.Matched)
	}
	if s.Missing != 2 {
		t.Errorf("expected Missing=2, got %d", s.Missing)
	}
	if s.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", s.Mismatched)
	}
}

func TestCompute_Empty(t *testing.T) {
	s := summary.Compute(nil)
	if s.Total != 0 || s.Matched != 0 || s.Missing != 0 || s.Mismatched != 0 {
		t.Errorf("expected all zeros for empty input, got %+v", s)
	}
}

func TestFormatText_ContainsFields(t *testing.T) {
	s := summary.Stats{Total: 5, Matched: 2, Missing: 2, Mismatched: 1}
	out := summary.FormatText(s)
	for _, want := range []string{"Total", "Matched", "Missing", "Mismatched", "5", "2", "1"} {
		if !strings.Contains(out, want) {
			t.Errorf("FormatText output missing %q\n%s", want, out)
		}
	}
}

func TestFormatJSON_ValidFields(t *testing.T) {
	s := summary.Stats{Total: 3, Matched: 1, Missing: 1, Mismatched: 1}
	out := summary.FormatJSON(s)
	for _, want := range []string{`"total":3`, `"matched":1`, `"missing":1`, `"mismatched":1`} {
		if !strings.Contains(out, want) {
			t.Errorf("FormatJSON output missing %q\n%s", want, out)
		}
	}
}
