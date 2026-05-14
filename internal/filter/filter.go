// Package filter provides utilities for filtering comparison results
// based on key patterns, prefixes, or exclusion rules.
package filter

import (
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Options holds the filtering configuration.
type Options struct {
	// Prefix restricts results to keys with the given prefix.
	Prefix string
	// ExcludeKeys is a list of exact key names to exclude from results.
	ExcludeKeys []string
	// OnlyMissing restricts results to missing-key differences only.
	OnlyMissing bool
	// OnlyMismatched restricts results to value-mismatch differences only.
	OnlyMismatched bool
}

// Apply filters a slice of Diff entries according to the given Options.
// It returns a new slice containing only the entries that match all criteria.
func Apply(diffs []comparator.Diff, opts Options) []comparator.Diff {
	excludeSet := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excludeSet[k] = true
	}

	var result []comparator.Diff
	for _, d := range diffs {
		if excludeSet[d.Key] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(d.Key, opts.Prefix) {
			continue
		}
		if opts.OnlyMissing && d.Type != comparator.DiffTypeMissing {
			continue
		}
		if opts.OnlyMismatched && d.Type != comparator.DiffTypeMismatch {
			continue
		}
		result = append(result, d)
	}
	return result
}
