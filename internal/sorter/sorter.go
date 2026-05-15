package sorter

import (
	"sort"

	"github.com/yourusername/envdiff/internal/comparator"
)

// SortBy defines the field by which results are sorted.
type SortBy string

const (
	// SortByKey sorts results alphabetically by key name.
	SortByKey SortBy = "key"
	// SortByStatus sorts results by status (missing before mismatched).
	SortByStatus SortBy = "status"
	// SortByFile sorts results by the file in which the issue was detected.
	SortByFile SortBy = "file"
)

// Apply sorts the given slice of DiffResult according to the specified SortBy field.
// If an unrecognised field is provided, the original order is preserved.
func Apply(results []comparator.DiffResult, by SortBy) []comparator.DiffResult {
	copied := make([]comparator.DiffResult, len(results))
	copy(copied, results)

	switch by {
	case SortByKey:
		sort.SliceStable(copied, func(i, j int) bool {
			return copied[i].Key < copied[j].Key
		})
	case SortByStatus:
		sort.SliceStable(copied, func(i, j int) bool {
			return statusRank(copied[i].Status) < statusRank(copied[j].Status)
		})
	case SortByFile:
		sort.SliceStable(copied, func(i, j int) bool {
			fi := firstFile(copied[i])
			fj := firstFile(copied[j])
			if fi == fj {
				return copied[i].Key < copied[j].Key
			}
			return fi < fj
		})
	}

	return copied
}

func statusRank(status comparator.DiffStatus) int {
	switch status {
	case comparator.StatusMissing:
		return 0
	case comparator.StatusMismatch:
		return 1
	default:
		return 2
	}
}

func firstFile(r comparator.DiffResult) string {
	if len(r.PresentIn) > 0 {
		return r.PresentIn[0]
	}
	return ""
}
