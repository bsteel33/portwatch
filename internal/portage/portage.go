// Package portage tracks when ports were first and last seen.
package portage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records the first and last observed time for a port.
type Entry struct {
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

// Tracker maintains age information for observed ports.
type Tracker struct {
	path    string
	now     func() time.Time
	entries map[string]*Entry
}

// New creates a Tracker, loading persisted state from path if it exists.
func New(path string, now func() time.Time) (*Tracker, error) {
	if now == nil {
		now = time.Now
	}
	t := &Tracker{path: path, now: now, entries: make(map[string]*Entry)}
	if err := t.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("portage: load: %w", err)
	}
	return t, nil
}

func key(p scanner.Port) string { return itoa(p.Port) + "/" + p.Proto }
func itoa(n int) string         { return fmt.Sprintf("%d", n) }

// Update records the current time as LastSeen for each port, and sets
// FirstSeen on the initial observation. Closed ports are evicted.
func (t *Tracker) Update(open []scanner.Port) error {
	now := t.now()
	seen := make(map[string]struct{}, len(open))
	for _, p := range open {
		k := key(p)
		seen[k] = struct{}{}
		if e, ok := t.entries[k]; ok {
			e.LastSeen = now
		} else {
			t.entries[k] = &Entry{FirstSeen: now, LastSeen: now}
		}
	}
	for k := range t.entries {
		if _, ok := seen[k]; !ok {
			delete(t.entries, k)
		}
	}
	return t.save()
}

// Get returns the age entry for a port, or nil if not tracked.
func (t *Tracker) Get(p scanner.Port) *Entry {
	e, ok := t.entries[key(p)]
	if !ok {
		return nil
	}
	return e
}

// Age returns how long a port has been open since first seen.
func (t *Tracker) Age(p scanner.Port) (time.Duration, bool) {
	e := t.Get(p)
	if e == nil {
		return 0, false
	}
	return t.now().Sub(e.FirstSeen), true
}

func (t *Tracker) save() error {
	data, err := json.Marshal(t.entries)
	if err != nil {
		return err
	}
	return os.WriteFile(t.path, data, 0o644)
}

func (t *Tracker) load() error {
	data, err := os.ReadFile(t.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &t.entries)
}
