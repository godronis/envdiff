// Package snapshot provides functionality for saving and loading point-in-time
// captures of parsed .env files, and computing the difference between two
// such captures.
//
// A Snapshot records the key-value pairs of an environment file along with
// a human-readable label and the UTC timestamp of creation. Snapshots are
// persisted as JSON files and can be compared using Diff to identify keys
// that were added, removed, or changed between two points in time.
//
// Typical usage:
//
//	snap := snapshot.Snapshot{Label: "pre-deploy", FilePath: ".env", Env: env}
//	_ = snapshot.Save(snap, "snapshots/pre-deploy.json")
//
//	loaded, _ := snapshot.Load("snapshots/pre-deploy.json")
//	result := snapshot.Diff(loaded, currentSnap)
package snapshot
