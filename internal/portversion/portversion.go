// Package portversion tracks version/banner strings observed on open ports
// and detects when they change between scans.
package portversion

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Entry holds the last observed version string for a port.
type Entry struct {
	Port    int    `json:"port"`
	Proto   string `json:"proto"`
	Version string `json:"version"`
}

// Change describes a version string that changed between two scans.
type Change struct {
	Port    int
	Proto   string
	OldVersion string
	NewVersion string
}

// Tracker stores per-port version strings and detects changes.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]Entry
	path    string
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// New loads a Tracker from path, or returns an empty one if the file is missing.
func New(path string) (*Tracker, error) {
	t := &Tracker{path: path, entries: make(map[string]Entry)}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return t, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &t.entries); err != nil {
		return nil, err
	}
	return t, nil
}

// Update records the version for a port and returns a Change if it differs
// from the previously stored value. Returns nil if unchanged or first seen.
func (t *Tracker) Update(port int, proto, version string) *Change {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := key(port, proto)
	prev, exists := t.entries[k]
	t.entries[k] = Entry{Port: port, Proto: proto, Version: version}
	if exists && prev.Version != version {
		return &Change{Port: port, Proto: proto, OldVersion: prev.Version, NewVersion: version}
	}
	return nil
}

// Get returns the stored entry for a port, and whether it exists.
func (t *Tracker) Get(port int, proto string) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key(port, proto)]
	return e, ok
}

// Save persists the tracker state to disk.
func (t *Tracker) Save() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	data, err := json.MarshalIndent(t.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(t.path, data, 0o644)
}

// Reset clears all stored entries.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]Entry)
}
