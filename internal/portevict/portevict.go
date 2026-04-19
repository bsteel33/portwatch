// Package portevict tracks ports that have been evicted (closed after being
// open for a sustained period) and records how long they were active.
package portevict

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records a single eviction event.
type Entry struct {
	Port     int       `json:"port"`
	Proto    string    `json:"proto"`
	OpenedAt time.Time `json:"opened_at"`
	ClosedAt time.Time `json:"closed_at"`
	Duration string    `json:"duration"`
}

// Log holds a list of eviction entries.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	path    string
}

// New creates a new eviction log backed by path. Existing entries are loaded
// if the file exists.
func New(path string) (*Log, error) {
	l := &Log{path: path}
	if err := l.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return l, nil
}

// Record appends an eviction entry and persists the log.
func (l *Log) Record(port int, proto string, openedAt, closedAt time.Time) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := Entry{
		Port:     port,
		Proto:    proto,
		OpenedAt: openedAt,
		ClosedAt: closedAt,
		Duration: closedAt.Sub(openedAt).Round(time.Second).String(),
	}
	l.entries = append(l.entries, e)
	return l.save()
}

// Entries returns a copy of all recorded evictions.
func (l *Log) Entries() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

func (l *Log) load() error {
	f, err := os.Open(l.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&l.entries)
}

func (l *Log) save() error {
	f, err := os.Create(l.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(l.entries)
}
