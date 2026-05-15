package exporter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/exporter"
)

func makeResults() []comparator.Result {
	return []comparator.Result{
		{Key: "DB_HOST", Status: comparator.StatusMissing, MissingIn: []string{"production.env"}},
		{Key: "API_KEY", Status: comparator.StatusMismatch},
		{Key: "PORT", Status: comparator.StatusMatch},
	}
}

func TestExport_EnvTemplate_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "template.env")

	err := exporter.Export(makeResults(), exporter.FormatEnvTemplate, out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=") {
		t.Error("expected DB_HOST= in env template")
	}
	if !strings.Contains(content, "API_KEY=") {
		t.Error("expected API_KEY= in env template")
	}
	if !strings.Contains(content, "MISSING") {
		t.Error("expected MISSING comment in env template")
	}
}

func TestExport_Markdown_ContainsTable(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "report.md")

	err := exporter.Export(makeResults(), exporter.FormatMarkdown, out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "| Key |") {
		t.Error("expected markdown table header")
	}
	if !strings.Contains(content, "DB_HOST") {
		t.Error("expected DB_HOST row in markdown")
	}
	if !strings.Contains(content, "Missing in:") {
		t.Error("expected 'Missing in:' detail in markdown")
	}
}

func TestExport_UnsupportedFormat_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.txt")

	err := exporter.Export(makeResults(), exporter.Format("xml"), out)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported export format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExport_InvalidPath_ReturnsError(t *testing.T) {
	err := exporter.Export(makeResults(), exporter.FormatMarkdown, "/nonexistent/path/report.md")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
