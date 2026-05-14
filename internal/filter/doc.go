// Package filter provides post-processing utilities for envdiff comparison
// results. It allows callers to narrow down a []comparator.Diff slice by
// applying one or more of the following criteria:
//
//   - Prefix: only include keys that start with a given string (e.g. "DB_").
//   - ExcludeKeys: drop specific keys by exact name.
//   - OnlyMissing: keep only entries where a key is absent from one or more
//     environment files.
//   - OnlyMismatched: keep only entries where the key exists in all files but
//     its value differs between them.
//
// Typical usage:
//
//	opts := filter.Options{
//		Prefix:      "DB_",
//		OnlyMissing: true,
//	}
//	filtered := filter.Apply(diffs, opts)
//
// Multiple criteria are combined with AND semantics — a diff entry must
// satisfy every active criterion to appear in the output.
package filter
