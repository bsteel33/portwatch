// Package portjournal records a running journal of port lifecycle events,
// persisting entries to disk for later review or export.
package portjournal

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// EventKind describes the type of journal entry.
type EventKind string

const (
	EventOpened EventKind = "opened"
	EventClosed EventKind = "closed"
	EventChanged EventKind = "changed"
)

// Entry is a single journal record.
type Entry struct {
	Time    time.Time `json:"time"`
	Port    int       `json:"port"`
	Proto   string    `json:"proto"`
	Service string    `json:"service,omitempty"`
	Kind    EventKind `json:"kind"`
	Note    string    `json:"note,omitempty"`
}

// Journal holds an ordered list of port lifecycle entries.
type Journal struct {
	mu      sync.Mutex
	path    string
	entries []Entry
}

// New loads an existing journal from path, or starts a new one.
func New(path string) (*Journal, error) {
	j := &Journal{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return j, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &j.entries); err != nil {
		return nil, err
	}
	return j, nil
}

// Record appends a new entry and persists the journal.
func (j *Journal) Record(e Entry) error {
	if e.Time.IsZero() {
		e.Time = time.Now()
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	j.entries = append(j.entries, e)
	return j.save()
}

// Entries returns a copy of all journal entries.
func (j *Journal) Entries() []Entry {
	j.mu.Lock()
	defer j.mu.Unlock()
	out := make([]Entry, len(j.entries))
	copy(out, j.entries)
	return out
}

// Last returns up to n most recent entries.
func (j *Journal) Last(n int) []Entry {
	j.mu.Lock()
	defer j.mu.Unlock()
	if n >= len(j.entries) {
		out := make([]Entry, len(j.entries))
		copy(out, j.entries)
		return out
	}
	out := make([]Entry, n)
	copy(out, j.entries[len(j.entries)-n:])
	return out
}

// Clear removes all entries and persists the empty journal.
func (j *Journal) Clear() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.entries = nil
	return j.save()
}

func (j *Journal) save() error {
	data, err := json.MarshalIndent(j.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(j.path, data, 0o644)
}
