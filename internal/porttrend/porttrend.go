// Package porttrend tracks directional change trends for open port counts
// over a sliding window, reporting whether the host is trending toward
// more open ports, fewer, or remaining stable.
package porttrend

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Direction represents the trend direction.
type Direction string

const (
	DirectionUp     Direction = "up"
	DirectionDown   Direction = "down"
	DirectionStable Direction = "stable"
)

// Sample is a single observation of open port count at a point in time.
type Sample struct {
	At    time.Time
	Count int
}

// Result summarises the current trend.
type Result struct {
	Direction Direction
	Delta     int
	Samples   int
}

// Tracker records port count samples and derives a trend.
type Tracker struct {
	cfg     Config
	clock   func() time.Time
	samples []Sample
}

// New returns a new Tracker using cfg.
func New(cfg Config) *Tracker {
	return &Tracker{cfg: cfg, clock: time.Now}
}

// Record adds a new sample for the given port count.
func (t *Tracker) Record(count int) {
	now := t.clock()
	t.samples = append(t.samples, Sample{At: now, Count: count})
	t.prune(now)
}

// Analyze returns the current trend over the retained window.
func (t *Tracker) Analyze() Result {
	n := len(t.samples)
	if n < 2 {
		return Result{Direction: DirectionStable, Samples: n}
	}
	first := t.samples[0].Count
	last := t.samples[n-1].Count
	delta := last - first
	dir := DirectionStable
	switch {
	case delta > t.cfg.Threshold:
		dir = DirectionUp
	case delta < -t.cfg.Threshold:
		dir = DirectionDown
	}
	return Result{Direction: dir, Delta: delta, Samples: n}
}

// Reset clears all recorded samples.
func (t *Tracker) Reset() {
	t.samples = nil
}

func (t *Tracker) prune(now time.Time) {
	cutoff := now.Add(-t.cfg.Window)
	i := 0
	for i < len(t.samples) && t.samples[i].At.Before(cutoff) {
		i++
	}
	t.samples = t.samples[i:]
}

// Fprint writes a human-readable trend summary to w.
func Fprint(w io.Writer, r Result) {
	fmt.Fprintf(w, "trend: %s (delta %+d, samples %d)\n", r.Direction, r.Delta, r.Samples)
}

// Print writes a trend summary to stdout.
func Print(r Result) { Fprint(os.Stdout, r) }
