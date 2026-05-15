package sorter_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/comparator"
	"github.com/yourusername/envdiff/internal/sorter"
)

func makeResults() []comparator.DiffResult {
	return []comparator.DiffResult{
		{Key: "ZEBRA", Status: comparator.StatusMismatch, PresentIn: []string{"b.env"}},
		{Key: "ALPHA", Status: comparator.StatusMissing, PresentIn: []string{"a.env"}},
		{Key: "MANGO", Status: comparator.StatusMismatch, PresentIn: []string{"a.env"}},
		{Key: "BETA", Status: comparator.StatusMissing, PresentIn: []string{"b.env"}},
	}
}

func TestApply_SortByKey(t *testing.T) {
	results := sorter.Apply(makeResults(), sorter.SortByKey)
	keys := []string{results[0].Key, results[1].Key, results[2].Key, results[3].Key}
	expected := []string{"ALPHA", "BETA", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestApply_SortByStatus(t *testing.T) {
	results := sorter.Apply(makeResults(), sorter.SortByStatus)
	if results[0].Status != comparator.StatusMissing {
		t.Errorf("expected first result to be StatusMissing, got %v", results[0].Status)
	}
	if results[len(results)-1].Status != comparator.StatusMismatch {
		t.Errorf("expected last result to be StatusMismatch, got %v", results[len(results)-1].Status)
	}
}

func TestApply_SortByFile(t *testing.T) {
	results := sorter.Apply(makeResults(), sorter.SortByFile)
	if results[0].PresentIn[0] != "a.env" {
		t.Errorf("expected first file to be a.env, got %q", results[0].PresentIn[0])
	}
}

func TestApply_UnknownSortBy_PreservesOrder(t *testing.T) {
	input := makeResults()
	results := sorter.Apply(input, sorter.SortBy("unknown"))
	for i, r := range results {
		if r.Key != input[i].Key {
			t.Errorf("index %d: order changed unexpectedly, got %q want %q", i, r.Key, input[i].Key)
		}
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	input := makeResults()
	originalFirst := input[0].Key
	sorter.Apply(input, sorter.SortByKey)
	if input[0].Key != originalFirst {
		t.Errorf("original slice was mutated: got %q, want %q", input[0].Key, originalFirst)
	}
}
