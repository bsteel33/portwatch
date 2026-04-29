// Package portrelay forwards port scan results to one or more downstream
// destinations (e.g. HTTP endpoints, files, or channels).
package portrelay

import (
	"fmt"
	"io"
	"sync"

	"github.com/user/portwatch/internal/snapshot"
)

// Destination is a sink that receives forwarded snapshots.
type Destination interface {
	Send(snap *snapshot.Snapshot) error
	Name() string
}

// Relay fans out a snapshot to multiple destinations.
type Relay struct {
	mu   sync.RWMutex
	dsts []Destination
	out  io.Writer
}

// New creates a new Relay writing errors to out.
func New(out io.Writer) *Relay {
	return &Relay{out: out}
}

// Register adds a destination to the relay.
func (r *Relay) Register(d Destination) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.dsts = append(r.dsts, d)
}

// Forward sends snap to every registered destination.
// Errors are printed to the relay's writer; all destinations are attempted.
func (r *Relay) Forward(snap *snapshot.Snapshot) {
	r.mu.RLock()
	dsts := make([]Destination, len(r.dsts))
	copy(dsts, r.dsts)
	r.mu.RUnlock()

	for _, d := range dsts {
		if err := d.Send(snap); err != nil {
			fmt.Fprintf(r.out, "portrelay: destination %q error: %v\n", d.Name(), err)
		}
	}
}

// Len returns the number of registered destinations.
func (r *Relay) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.dsts)
}
