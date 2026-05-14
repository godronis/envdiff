// Package comparator implements the core diff logic for envdiff.
//
// It accepts parsed environment maps (produced by the parser package) and
// computes two categories of discrepancies:
//
//   - Missing keys: a key present in at least one environment but absent in
//     another.
//   - Mismatched values: a key present in all environments but with differing
//     values.
//
// Basic usage:
//
//	envs := map[string]map[string]string{
//	    "dev":  {"PORT": "3000", "DB_HOST": "localhost"},
//	    "prod": {"PORT": "80",   "DB_HOST": "db.example.com"},
//	}
//
//	result := comparator.Compare(envs)
//	// result.MissingIn  — keys absent per environment
//	// result.Mismatched — keys whose values differ
package comparator
