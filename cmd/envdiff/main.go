package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <file1> <file2> [file3...]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	envs := make(map[string]map[string]string, len(args))
	for _, path := range args {
		parsed, err := parser.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			os.Exit(1)
		}
		envs[path] = parsed
	}

	result := comparator.Compare(envs)

	if err := reporter.Report(os.Stdout, result, reporter.Format(*format)); err != nil {
		fmt.Fprintf(os.Stderr, "error generating report: %v\n", err)
		os.Exit(1)
	}

	if len(result.Missing) > 0 || len(result.Mismatched) > 0 {
		os.Exit(2)
	}
}
