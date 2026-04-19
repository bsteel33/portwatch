// Package trend tracks port count changes over time and detects growth patterns.
package trend

import (
	"sync"
	"time"
)

// Point represents a single observation.
type Point struct {
	At    time.Time
	Count int
}

// Trend holds a sliding window of port count observations.
type Trend struct {
	mu     sync.Mutex
	points []Point
	window time.Duration
	now    func() time.Time
}

// New returns a Trend that retains points within the given window.
func New(cfg Config) *Trend {
	return &Trend{
		window: cfg.Window,
		now:    time.Now,
	}
}

// Record adds a new observation.
func (t *Trend) Record(count int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	at := t.now()
	t.points = append(t.points, Point{At: at, Count: count})
	t.prune(at)
}

// Points returns a copy of the current window of observations.
func (t *Trend) Points() []Point {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Point, len(t.points))
	copy(out, t.points)
	return out
}

// Delta returns the difference between the latest and earliest count in the window.
// Returns 0 if fewer than two points exist.
func (t *Trend) Delta() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.points) < 2 {
		return 0
	}
	return t.points[len(t.points)-1].Count - t.points[0].Count
}

// Reset clears all recorded points.
func (t *Trend) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.points = nil
}

func (t *Trend) prune(now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(t.points) && t.points[i].At.Before(cutoff) {
		i++
	}
	t.points = t.points[i:]
}
