// Package portwatch ties together scanning, diffing, and alerting into a
// single reusable Watch call used by the CLI and daemon.
package portwatch

import (
	"time"

	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/snapshot"
)

// Result holds the outcome of a single watch cycle.
type Result struct {
	Ports    []scanner.Port
	Diff     snapshot.Diff
	ScannedAt time.Time
	Changed  bool
}

// Watcher performs a port scan and compares the result against a saved
// snapshot, returning a Result that callers can act on.
type Watcher struct {
	scanner  *scanner.Scanner
	snap     *snapshot.Snapshot
}

// Config holds tunable parameters for the Watcher.
type Config struct {
	SnapshotPath string
	Ports        []int
	Protocol     string
	Timeout      time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		SnapshotPath: "/var/lib/portwatch/snapshot.json",
		Protocol:     "tcp",
		Timeout:      2 * time.Second,
	}
}

// New creates a Watcher from cfg.
func New(cfg Config) (*Watcher, error) {
	s := scanner.New(scanner.Config{
		Protocol: cfg.Protocol,
		Timeout:  cfg.Timeout,
	})
	snap, err := snapshot.Load(cfg.SnapshotPath)
	if err != nil {
		snap = snapshot.New(cfg.SnapshotPath)
	}
	return &Watcher{scanner: s, snap: snap}, nil
}

// Run performs one scan cycle and returns the Result.
func (w *Watcher) Run(ports []int) (Result, error) {
	scanned, err := w.scanner.OpenPorts(ports)
	if err != nil {
		return Result{}, err
	}
	diff := snapshot.Compare(w.snap, scanned)
	changed := len(diff.Added) > 0 || len(diff.Removed) > 0
	if changed {
		if err := w.snap.Save(scanned); err != nil {
			return Result{}, err
		}
	}
	return Result{
		Ports:     scanned,
		Diff:      diff,
		ScannedAt: time.Now(),
		Changed:   changed,
	}, nil
}
