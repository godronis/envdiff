package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/reporter"
)

func TestReportText_NoDifferences(t *testing.T) {
	result := comparator.Result{}
	var buf bytes.Buffer
	if err := reporter.Report(&buf, result, reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestReportText_MissingKey(t *testing.T) {
	result := comparator.Result{
		Missing: []comparator.MissingKey{{File: ".env.prod", Key: "SECRET"}},
	}
	var buf bytes.Buffer
	if err := reporter.Report(&buf, result, reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Missing keys") || !strings.Contains(out, "SECRET") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestReportText_MismatchedValue(t *testing.T) {
	result := comparator.Result{
		Mismatched: []comparator.MismatchedKey{
			{Key: "DB_HOST", Values: map[string]string{".env.dev": "localhost", ".env.prod": "db.prod"}},
		},
	}
	var buf bytes.Buffer
	if err := reporter.Report(&buf, result, reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Mismatched values") || !strings.Contains(out, "DB_HOST") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestReportJSON_ValidOutput(t *testing.T) {
	result := comparator.Result{
		Missing: []comparator.MissingKey{{File: ".env.staging", Key: "API_KEY"}},
	}
	var buf bytes.Buffer
	if err := reporter.Report(&buf, result, reporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") || !strings.Contains(out, "{") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestReport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := reporter.Report(&buf, comparator.Result{}, reporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
