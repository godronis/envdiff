package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/filter"
)

var sampleDiffs = []comparator.Diff{
	{Key: "DB_HOST", Type: comparator.DiffTypeMissing, File: "prod.env"},
	{Key: "DB_PORT", Type: comparator.DiffTypeMismatch, Values: map[string]string{"dev.env": "5432", "prod.env": "5433"}},
	{Key: "APP_SECRET", Type: comparator.DiffTypeMissing, File: "staging.env"},
	{Key: "APP_DEBUG", Type: comparator.DiffTypeMismatch, Values: map[string]string{"dev.env": "true", "prod.env": "false"}},
}

func TestApply_NoOptions(t *testing.T) {
	result := filter.Apply(sampleDiffs, filter.Options{})
	if len(result) != len(sampleDiffs) {
		t.Errorf("expected %d diffs, got %d", len(sampleDiffs), len(result))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	result := filter.Apply(sampleDiffs, filter.Options{Prefix: "DB_"})
	if len(result) != 2 {
		t.Errorf("expected 2 diffs with prefix DB_, got %d", len(result))
	}
	for _, d := range result {
		if d.Key[:3] != "DB_" {
			t.Errorf("unexpected key without DB_ prefix: %s", d.Key)
		}
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	result := filter.Apply(sampleDiffs, filter.Options{ExcludeKeys: []string{"APP_SECRET", "DB_HOST"}})
	if len(result) != 2 {
		t.Errorf("expected 2 diffs after exclusion, got %d", len(result))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	result := filter.Apply(sampleDiffs, filter.Options{OnlyMissing: true})
	if len(result) != 2 {
		t.Errorf("expected 2 missing diffs, got %d", len(result))
	}
	for _, d := range result {
		if d.Type != comparator.DiffTypeMissing {
			t.Errorf("expected DiffTypeMissing, got %s", d.Type)
		}
	}
}

func TestApply_OnlyMismatched(t *testing.T) {
	result := filter.Apply(sampleDiffs, filter.Options{OnlyMismatched: true})
	if len(result) != 2 {
		t.Errorf("expected 2 mismatch diffs, got %d", len(result))
	}
	for _, d := range result {
		if d.Type != comparator.DiffTypeMismatch {
			t.Errorf("expected DiffTypeMismatch, got %s", d.Type)
		}
	}
}

func TestApply_EmptyDiffs(t *testing.T) {
	result := filter.Apply(nil, filter.Options{Prefix: "DB_"})
	if result != nil && len(result) != 0 {
		t.Errorf("expected empty result for nil input, got %v", result)
	}
}
