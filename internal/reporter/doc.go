// Package reporter provides functionality for rendering envdiff comparison
// results in multiple output formats.
//
// Supported formats:
//
//   - text: Human-readable plain text output with missing and mismatched keys
//     clearly labeled and indented for readability.
//
//   - json: Structured JSON output suitable for machine consumption or
//     integration with other tooling.
//
// Usage:
//
//	err := reporter.Report(os.Stdout, result, reporter.FormatText)
//
// The Report function accepts any io.Writer, making it easy to redirect
// output to files, buffers, or standard streams.
package reporter
