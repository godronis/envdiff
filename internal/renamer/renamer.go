// Package renamer provides functionality to rename or remap keys across
// parsed environment maps using a mapping definition.
package renamer

import "fmt"

// Mapping defines a source key and the target key it should be renamed to.
type Mapping struct {
	From string
	To   string
}

// Result holds the outcome of a rename operation on a single env map.
type Result struct {
	// Renamed contains keys that were successfully renamed.
	Renamed []Mapping
	// NotFound contains source keys from the mapping that were absent in the env.
	NotFound []string
	// Env is the transformed environment map after renaming.
	Env map[string]string
}

// Apply renames keys in the provided env map according to the given mappings.
// Keys not present in the env are recorded in Result.NotFound.
// If a target key already exists, an error is returned to prevent silent overwrites.
func Apply(env map[string]string, mappings []Mapping) (Result, error) {
	if env == nil {
		return Result{Env: map[string]string{}}, nil
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var renamed []Mapping
	var notFound []string

	for _, m := range mappings {
		val, exists := out[m.From]
		if !exists {
			notFound = append(notFound, m.From)
			continue
		}
		if _, conflict := out[m.To]; conflict && m.From != m.To {
			return Result{}, fmt.Errorf("renamer: target key %q already exists in env", m.To)
		}
		delete(out, m.From)
		out[m.To] = val
		renamed = append(renamed, m)
	}

	return Result{
		Renamed:  renamed,
		NotFound: notFound,
		Env:      out,
	}, nil
}
