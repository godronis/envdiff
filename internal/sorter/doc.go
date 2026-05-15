// Package sorter provides utilities for ordering comparator.DiffResult slices
// produced by the envdiff comparison pipeline.
//
// Supported sort fields:
//
//   - SortByKey    – alphabetical order by environment variable name.
//   - SortByStatus – groups missing keys before mismatched keys.
//   - SortByFile   – groups results by the first file they appear in,
//     with a secondary sort by key name within each group.
//
// Apply always returns a new slice, leaving the original untouched.
package sorter
