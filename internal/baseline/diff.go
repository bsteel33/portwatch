package baseline

import (
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Deviation describes ports that deviate from the baseline.
type Deviation struct {
	Added   []scanner.Port // present now but not in baseline
	Removed []scanner.Port // in baseline but no longer present
}

// HasChanges reports whether any deviation exists.
func (d Deviation) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// Compare returns the deviation between the baseline and current ports.
func Compare(base *Baseline, current []scanner.Port) Deviation {
	diff := snapshot.Compare(
		&snapshot.Snapshot{Ports: base.Ports},
		&snapshot.Snapshot{Ports: current},
	)
	return Deviation{
		Added:   diff.Added,
		Removed: diff.Removed,
	}
}
