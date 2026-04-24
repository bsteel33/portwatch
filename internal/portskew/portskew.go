// Package portskew detects statistical skew in port distributions,
// flagging hosts whose open-port profile deviates significantly from
// a rolling baseline mean.
package portskew

import (
	"math"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Result holds the skew analysis for a single scan.
type Result struct {
	Mean      float64
	StdDev    float64
	Current   int
	ZScore    float64
	Skewed    bool
}

// Detector tracks historical port counts and computes z-scores.
type Detector struct {
	mu        sync.Mutex
	threshold float64
	samples   []float64
}

// Config holds Detector configuration.
type Config struct {
	Threshold float64 // z-score threshold above which a scan is flagged
	MinSamples int    // minimum samples before flagging
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Threshold:  2.0,
		MinSamples: 5,
	}
}

// New creates a Detector with the given config.
func New(cfg Config) *Detector {
	return &Detector{threshold: cfg.Threshold}
}

// Analyze records the port count from ports and returns a Result.
func (d *Detector) Analyze(ports []scanner.Port, cfg Config) Result {
	d.mu.Lock()
	defer d.mu.Unlock()

	current := float64(len(ports))
	d.samples = append(d.samples, current)

	if len(d.samples) < cfg.MinSamples {
		return Result{Current: int(current)}
	}

	mean := d.mean()
	std := d.stddev(mean)

	var z float64
	if std > 0 {
		z = (current - mean) / std
	}

	return Result{
		Mean:    mean,
		StdDev:  std,
		Current: int(current),
		ZScore:  z,
		Skewed:  math.Abs(z) >= cfg.Threshold,
	}
}

// Reset clears all recorded samples.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.samples = nil
}

func (d *Detector) mean() float64 {
	var sum float64
	for _, s := range d.samples {
		sum += s
	}
	return sum / float64(len(d.samples))
}

func (d *Detector) stddev(mean float64) float64 {
	var variance float64
	for _, s := range d.samples {
		diff := s - mean
		variance += diff * diff
	}
	variance /= float64(len(d.samples))
	return math.Sqrt(variance)
}
