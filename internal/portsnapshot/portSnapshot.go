// Package portSnapshot provides periodic snapshotting of port state
// with change detection and retention management.
package portSnapshot

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds a single snapshot record.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Ports     []scanner.Port `json:"ports"`
}

// Snapshotter manages a rolling window of port snapshots.
type Snapshotter struct {
	mu      sync.Mutex
	path    string
	entries []Entry
	cfg     Config
}

// New creates a Snapshotter, loading any existing data from path.
func New(path string, cfg Config) (*Snapshotter, error) {
	s := &Snapshotter{path: path, cfg: cfg}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Record appends a new snapshot entry and prunes old ones.
func (s *Snapshotter) Record(ports []scanner.Port) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = append(s.entries, Entry{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	})
	s.prune()
	return s.save()
}

// Last returns the most recent snapshot entry, or false if none exist.
func (s *Snapshotter) Last() (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.entries) == 0 {
		return Entry{}, false
	}
	return s.entries[len(s.entries)-1], true
}

// All returns all retained snapshot entries.
func (s *Snapshotter) All() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

func (s *Snapshotter) prune() {
	cutoff := time.Now().UTC().Add(-s.cfg.Retention)
	start := 0
	for start < len(s.entries) && s.entries[start].Timestamp.Before(cutoff) {
		start++
	}
	s.entries = s.entries[start:]
	if s.cfg.MaxEntries > 0 && len(s.entries) > s.cfg.MaxEntries {
		s.entries = s.entries[len(s.entries)-s.cfg.MaxEntries:]
	}
}

func (s *Snapshotter) save() error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s.entries)
}

func (s *Snapshotter) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.entries)
}
