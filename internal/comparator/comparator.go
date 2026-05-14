package comparator

import "sort"

// MissingKey represents a key that is present in at least one env file
// but absent in another.
type MissingKey struct {
	Key  string `json:"key"`
	File string `json:"file"`
}

// MismatchedKey represents a key that exists in all files but has
// differing values across them.
type MismatchedKey struct {
	Key    string            `json:"key"`
	Values map[string]string `json:"values"`
}

// Result holds the full output of a comparison between env files.
type Result struct {
	Missing    []MissingKey    `json:"missing"`
	Mismatched []MismatchedKey `json:"mismatched"`
}

// Compare takes a map of filename -> key/value pairs and returns a Result
// describing all missing and mismatched keys across the provided env files.
func Compare(envs map[string]map[string]string) Result {
	keys := collectAllKeys(envs)
	var result Result

	for _, key := range keys {
		presentIn := []string{}
		values := map[string]string{}

		for file, env := range envs {
			if val, ok := env[key]; ok {
				presentIn = append(presentIn, file)
				values[file] = val
			}
		}

		if len(presentIn) < len(envs) {
			for file := range envs {
				if _, ok := envs[file][key]; !ok {
					result.Missing = append(result.Missing, MissingKey{
						Key:  key,
						File: file,
					})
				}
			}
		} else if hasMismatch(values) {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:    key,
				Values: values,
			})
		}
	}

	return result
}

func collectAllKeys(envs map[string]map[string]string) []string {
	seen := map[string]struct{}{}
	for _, env := range envs {
		for key := range env {
			seen[key] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

func hasMismatch(values map[string]string) bool {
	var first string
	var set bool
	for _, v := range values {
		if !set {
			first = v
			set = true
			continue
		}
		if v != first {
			return true
		}
	}
	return false
}

func sortStrings(s []string) {
	sort.Strings(s)
}
