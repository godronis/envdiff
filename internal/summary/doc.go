// Package summary provides utilities for computing and formatting
// aggregated statistics from envdiff comparison results.
//
// It accepts a slice of comparator.Result values and produces a Stats
// struct containing counts of matched, missing, and mismatched keys.
// The stats can be rendered as plain text or JSON for inclusion in
// reporter output.
//
// Example usage:
//
//	results := comparator.Compare(envs)
//	stats := summary.Compute(results)
//	fmt.Print(summary.FormatText(stats))
package summary
