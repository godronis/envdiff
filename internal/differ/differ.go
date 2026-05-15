// Package differ provides utilities for generating a unified diff-style
// representation of changes between two .env file contents.
package differ

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Line represents a single line in a diff output.
type Line struct {
	Key    string
	Old    string
	New    string
	Status string // "added", "removed", "changed", "unchanged"
}

// Diff holds the result of comparing two environments.
type Diff struct {
	FileA string
	FileB string
	Lines []Line
}

// Generate produces a Diff from a slice of comparator.Result entries,
// focusing on the first two files present in results.
func Generate(fileA, fileB string, results []comparator.Result) Diff {
	d := Diff{FileA: fileA, FileB: fileB}

	for _, r := range results {
		valA, hasA := r.Values[fileA]
		valB, hasB := r.Values[fileB]

		switch {
		case !hasA && hasB:
			d.Lines = append(d.Lines, Line{Key: r.Key, Old: "", New: valB, Status: "added"})
		case hasA && !hasB:
			d.Lines = append(d.Lines, Line{Key: r.Key, Old: valA, New: "", Status: "removed"})
		case hasA && hasB && valA != valB:
			d.Lines = append(d.Lines, Line{Key: r.Key, Old: valA, New: valB, Status: "changed"})
		default:
			d.Lines = append(d.Lines, Line{Key: r.Key, Old: valA, New: valB, Status: "unchanged"})
		}
	}

	return d
}

// FormatUnified returns a human-readable unified diff string.
func FormatUnified(d Diff) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "--- %s\n", d.FileA)
	fmt.Fprintf(&sb, "+++ %s\n", d.FileB)

	for _, l := range d.Lines {
		switch l.Status {
		case "added":
			fmt.Fprintf(&sb, "+ %s=%s\n", l.Key, l.New)
		case "removed":
			fmt.Fprintf(&sb, "- %s=%s\n", l.Key, l.Old)
		case "changed":
			fmt.Fprintf(&sb, "~ %s: %q -> %q\n", l.Key, l.Old, l.New)
		default:
			fmt.Fprintf(&sb, "  %s=%s\n", l.Key, l.Old)
		}
	}

	return sb.String()
}
