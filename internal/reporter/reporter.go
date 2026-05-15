package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/envdiff/internal/comparator"
	"github.com/yourusername/envdiff/internal/summary"
)

// Report writes comparison results to w in the given format ("text" or "json").
// It appends a summary section after the per-key output.
func Report(w io.Writer, results []comparator.Result, format string) error {
	switch strings.ToLower(format) {
	case "text", "":
		return reportText(w, results)
	case "json":
		return reportJSON(w, results)
	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func reportText(w io.Writer, results []comparator.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}
	for _, r := range results {
		switch r.Status {
		case comparator.StatusMissing:
			if _, err := fmt.Fprintf(w, "[MISSING]  %s (absent in: %s)\n", r.Key, strings.Join(r.MissingIn, ", ")); err != nil {
				return err
			}
		case comparator.StatusMismatch:
			if _, err := fmt.Fprintf(w, "[MISMATCH] %s\n", r.Key); err != nil {
				return err
			}
			for file, val := range r.Values {
				if _, err := fmt.Fprintf(w, "             %s = %q\n", file, val); err != nil {
					return err
				}
			}
		}
	}
	// Append summary
	stats := summary.Compute(results)
	_, err := fmt.Fprint(w, summary.FormatText(stats))
	return err
}

type jsonOutput struct {
	Results []comparator.Result `json:"results"`
	Summary summary.Stats       `json:"summary"`
}

func reportJSON(w io.Writer, results []comparator.Result) error {
	out := jsonOutput{
		Results: results,
		Summary: summary.Compute(results),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
