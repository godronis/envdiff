package snapshot

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// FormatText returns a human-readable summary of a DiffResult.
func FormatText(d DiffResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Snapshot diff: %s → %s\n", d.BaseLabel, d.CurrentLabel)

	sort.Strings(d.Added)
	sort.Strings(d.Removed)
	sort.Strings(d.Changed)

	if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 {
		sb.WriteString("  No changes detected.\n")
		return sb.String()
	}
	for _, k := range d.Added {
		fmt.Fprintf(&sb, "  + %s (added)\n", k)
	}
	for _, k := range d.Removed {
		fmt.Fprintf(&sb, "  - %s (removed)\n", k)
	}
	for _, k := range d.Changed {
		fmt.Fprintf(&sb, "  ~ %s (changed)\n", k)
	}
	return sb.String()
}

// FormatJSON returns a JSON-encoded summary of a DiffResult.
func FormatJSON(d DiffResult) (string, error) {
	type payload struct {
		Base    string   `json:"base"`
		Current string   `json:"current"`
		Added   []string `json:"added"`
		Removed []string `json:"removed"`
		Changed []string `json:"changed"`
	}
	sort.Strings(d.Added)
	sort.Strings(d.Removed)
	sort.Strings(d.Changed)

	p := payload{
		Base:    d.BaseLabel,
		Current: d.CurrentLabel,
		Added:   orEmpty(d.Added),
		Removed: orEmpty(d.Removed),
		Changed: orEmpty(d.Changed),
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", fmt.Errorf("snapshot: failed to marshal diff: %w", err)
	}
	return string(data), nil
}

func orEmpty(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
