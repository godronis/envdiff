package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of a parsed .env file at a point in time.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	FilePath  string            `json:"file_path"`
	Env       map[string]string `json:"env"`
}

// Save writes the snapshot to the given destination path as JSON.
func Save(s Snapshot, destPath string) error {
	if s.Label == "" {
		return fmt.Errorf("snapshot label must not be empty")
	}
	if len(s.Env) == 0 {
		return fmt.Errorf("snapshot env must not be empty")
	}
	s.CreatedAt = time.Now().UTC()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot to %s: %w", destPath, err)
	}
	return nil
}

// Load reads a snapshot from the given path.
func Load(srcPath string) (Snapshot, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return Snapshot{}, fmt.Errorf("failed to read snapshot from %s: %w", srcPath, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("failed to parse snapshot: %w", err)
	}
	return s, nil
}

// Diff compares two snapshots and returns keys that were added, removed, or changed.
func Diff(base, current Snapshot) DiffResult {
	result := DiffResult{
		BaseLabel:    base.Label,
		CurrentLabel: current.Label,
	}
	for k, v := range current.Env {
		if bv, ok := base.Env[k]; !ok {
			result.Added = append(result.Added, k)
		} else if bv != v {
			result.Changed = append(result.Changed, k)
		}
	}
	for k := range base.Env {
		if _, ok := current.Env[k]; !ok {
			result.Removed = append(result.Removed, k)
		}
	}
	return result
}

// DiffResult holds the keys that changed between two snapshots.
type DiffResult struct {
	BaseLabel    string
	CurrentLabel string
	Added        []string
	Removed      []string
	Changed      []string
}
