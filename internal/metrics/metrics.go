// Package metrics tracks runtime statistics for portwatch scans.
package metrics

import (
	"sync"
	"time"
)

// Stats holds accumulated scan metrics.
type Stats struct {
	mu           sync.Mutex
	ScanCount    int
	AlertCount   int
	LastScanTime time.Time
	LastScanDur  time.Duration
	OpenPortsHWM int // high-water mark of open ports seen
}

// Collector records metrics from scan cycles.
type Collector struct {
	stats Stats
}

// New returns a new Collector.
func New() *Collector {
	return &Collector{}
}

// RecordScan records the result of a completed scan.
func (c *Collector) RecordScan(dur time.Duration, openPorts int, alerted bool) {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.ScanCount++
	c.stats.LastScanTime = time.Now()
	c.stats.LastScanDur = dur
	if alerted {
		c.stats.AlertCount++
	}
	if openPorts > c.stats.OpenPortsHWM {
		c.stats.OpenPortsHWM = openPorts
	}
}

// Snapshot returns a copy of the current stats.
func (c *Collector) Snapshot() Stats {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	return Stats{
		ScanCount:    c.stats.ScanCount,
		AlertCount:   c.stats.AlertCount,
		LastScanTime: c.stats.LastScanTime,
		LastScanDur:  c.stats.LastScanDur,
		OpenPortsHWM: c.stats.OpenPortsHWM,
	}
}

// Reset clears all recorded metrics.
func (c *Collector) Reset() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats = Stats{}
}
