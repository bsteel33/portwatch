package watch

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Watcher continuously scans ports and emits diffs when changes are detected.
type Watcher struct {
	scanner  *scanner.Scanner
	interval time.Duration
	OnChange func(diff snapshot.Diff)
}

// New creates a Watcher with the given scanner and poll interval.
func New(s *scanner.Scanner, interval time.Duration) *Watcher {
	return &Watcher{
		scanner:  s,
		interval: interval,
		OnChange: func(snapshot.Diff) {},
	}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context, snapshotPath string) error {
	prev, err := snapshot.Load(snapshotPath)
	if err != nil {
		log.Printf("watch: no previous snapshot, starting fresh: %v", err)
		prev = &snapshot.Snapshot{}
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			curr, err := w.scan(snapshotPath)
			if err != nil {
				log.Printf("watch: scan error: %v", err)
				continue
			}
			diff := snapshot.Compare(prev, curr)
			if diff.HasChanges() {
				w.OnChange(diff)
			}
			prev = curr
		}
	}
}

func (w *Watcher) scan(snapshotPath string) (*snapshot.Snapshot, error) {
	ports, err := w.scanner.OpenPorts()
	if err != nil {
		return nil, err
	}
	snap := snapshot.New(ports)
	if err := snap.Save(snapshotPath); err != nil {
		log.Printf("watch: failed to save snapshot: %v", err)
	}
	return snap, nil
}
