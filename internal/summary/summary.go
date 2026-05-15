package summary

import (
	"fmt"
	"strings"

	"github.com/yourusername/envdiff/internal/comparator"
)

// Stats holds aggregated counts from a comparison result set.
type Stats struct {
	Total      int
	Missing    int
	Mismatched int
	Matched    int
}

// Compute calculates statistics from a slice of comparator results.
func Compute(results []comparator.Result) Stats {
	s := Stats{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case comparator.StatusMissing:
			s.Missing++
		case comparator.StatusMismatch:
			s.Mismatched++
		case comparator.StatusMatch:
			s.Matched++
		}
	}
	return s
}

// FormatText returns a human-readable summary string.
func FormatText(s Stats) string {
	var sb strings.Builder
	sb.WriteString("--- Summary ---\n")
	sb.WriteString(fmt.Sprintf("Total keys : %d\n", s.Total))
	sb.WriteString(fmt.Sprintf("Matched    : %d\n", s.Matched))
	sb.WriteString(fmt.Sprintf("Missing    : %d\n", s.Missing))
	sb.WriteString(fmt.Sprintf("Mismatched : %d\n", s.Mismatched))
	return sb.String()
}

// FormatJSON returns a JSON-encoded summary string.
func FormatJSON(s Stats) string {
	return fmt.Sprintf(
		`{"total":%d,"matched":%d,"missing":%d,"mismatched":%d}`,
		s.Total, s.Matched, s.Missing, s.Mismatched,
	)
}
