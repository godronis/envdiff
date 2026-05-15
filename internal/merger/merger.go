// Package merger provides functionality to merge multiple .env file
// maps into a single unified map, with configurable conflict resolution
// strategies (first-wins or last-wins).
package merger

import "sort"

// Strategy defines how conflicting keys are resolved during a merge.
type Strategy int

const (
	// FirstWins keeps the value from the first file that defines a key.
	FirstWins Strategy = iota
	// LastWins overwrites with the value from the last file that defines a key.
	LastWins
)

// Result holds the merged key-value pairs and metadata about the merge.
type Result struct {
	// Merged is the final merged map of key to value.
	Merged map[string]string
	// Conflicts lists keys that appeared in more than one source file.
	Conflicts []string
	// Sources maps each key to the filename that contributed its final value.
	Sources map[string]string
}

// Merge combines multiple named env maps into a single Result.
// The names slice must correspond positionally to the envs slice.
// Strategy controls which value wins when the same key appears in multiple files.
func Merge(names []string, envs []map[string]string, strategy Strategy) Result {
	merged := make(map[string]string)
	sources := make(map[string]string)
	conflictSet := make(map[string]bool)

	for i, env := range envs {
		name := ""
		if i < len(names) {
			name = names[i]
		}
		for k, v := range env {
			if _, exists := merged[k]; exists {
				conflictSet[k] = true
				if strategy == LastWins {
					merged[k] = v
					sources[k] = name
				}
			} else {
				merged[k] = v
				sources[k] = name
			}
		}
	}

	conflicts := make([]string, 0, len(conflictSet))
	for k := range conflictSet {
		conflicts = append(conflicts, k)
	}
	sort.Strings(conflicts)

	return Result{
		Merged:    merged,
		Conflicts: conflicts,
		Sources:   sources,
	}
}
