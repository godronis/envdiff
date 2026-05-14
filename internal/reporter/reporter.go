package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes the comparison result to the given writer in the specified format.
func Report(w io.Writer, result comparator.Result, format Format) error {
	switch format {
	case FormatJSON:
		return reportJSON(w, result)
	case FormatText:
		return reportText(w, result)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func reportText(w io.Writer, result comparator.Result) error {
	if len(result.Missing) == 0 && len(result.Mismatched) == 0 {
		_, err := fmt.Fprintln(w, "✓ No differences found.")
		return err
	}

	var sb strings.Builder

	if len(result.Missing) > 0 {
		sb.WriteString("Missing keys:\n")
		for _, m := range result.Missing {
			sb.WriteString(fmt.Sprintf("  [%s] %s\n", m.File, m.Key))
		}
	}

	if len(result.Mismatched) > 0 {
		sb.WriteString("Mismatched values:\n")
		for _, m := range result.Mismatched {
			sb.WriteString(fmt.Sprintf("  %s:\n", m.Key))
			for file, val := range m.Values {
				sb.WriteString(fmt.Sprintf("    %s = %q\n", file, val))
			}
		}
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}

func reportJSON(w io.Writer, result comparator.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
