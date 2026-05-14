// Package comparator provides functionality for comparing parsed .env files
// across multiple environments, identifying missing and mismatched keys.
package comparator

// Result holds the comparison outcome between two or more env files.
type Result struct {
	// MissingIn maps an environment name to keys missing from that environment.
	MissingIn map[string][]string
	// Mismatched contains keys whose values differ across environments.
	Mismatched []MismatchedKey
}

// MismatchedKey describes a key whose value differs across environments.
type MismatchedKey struct {
	Key    string
	Values map[string]string // env name -> value
}

// Compare takes a map of environment name -> parsed key/value pairs and
// returns a Result describing all missing and mismatched keys.
func Compare(envs map[string]map[string]string) Result {
	allKeys := collectAllKeys(envs)

	result := Result{
		MissingIn:  make(map[string][]string),
		Mismatched: []MismatchedKey{},
	}

	for _, key := range allKeys {
		values := make(map[string]string)
		presentIn := 0

		for envName, pairs := range envs {
			val, ok := pairs[key]
			if !ok {
				result.MissingIn[envName] = append(result.MissingIn[envName], key)
			} else {
				values[envName] = val
				presentIn++
			}
		}

		if presentIn > 1 && hasMismatch(values) {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:    key,
				Values: values,
			})
		}
	}

	return result
}

// collectAllKeys returns a deduplicated, sorted list of all keys across envs.
func collectAllKeys(envs map[string]map[string]string) []string {
	seen := make(map[string]struct{})
	for _, pairs := range envs {
		for k := range pairs {
			seen[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

// hasMismatch returns true if not all values in the map are identical.
func hasMismatch(values map[string]string) bool {
	var first string
	init := false
	for _, v := range values {
		if !init {
			first = v
			init = true
			continue
		}
		if v != first {
			return true
		}
	}
	return false
}

// sortStrings sorts a string slice in place (simple insertion sort).
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
