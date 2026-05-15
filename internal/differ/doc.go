// Package differ generates a unified diff-style comparison between two
// .env environments based on comparator results.
//
// It produces a structured Diff containing Line entries that describe
// whether each key was added, removed, changed, or unchanged between
// the two files. The FormatUnified function renders this as a
// human-readable diff string using +, -, ~, and space prefixes.
//
// Example usage:
//
//	results, _ := comparator.Compare(envs)
//	d := differ.Generate("dev.env", "prod.env", results)
//	fmt.Print(differ.FormatUnified(d))
package differ
