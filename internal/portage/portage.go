// Package portage tracks how long each port has been continuously open.
package portage

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records when a port was first seen open.
type Entry struct {
	Port     int       `json:"port"`
	Proto    string    `json:"proto"`
	FirstSeen time.Time `json:"first_seen"`
}

// Age returns how long the port has been open.
func (e Entry) Age(now time.Time) time.Duration {
	return now.Sub(e.FirstSeen)
}

// Tracker maintains first-seen timestamps for open ports.
type Tracker struct {
	path    string
	now     func() time.Time
	entries map[string]Entry
}

func key(port int, proto string) string {
	return proto + ":" + itoa(port)
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

// New loads or initialises a Tracker backed by path.
func New(path string) (*Tracker, error) {
	t := &Tracker{path: path, now: time.Now, entries: make(map[string]Entry)}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return t, nil
		}
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, e := range list {
		t.entries[key(e.Port, e.Proto)] = e
	}
	return t, nil
}

// Update records first-seen for any new ports and evicts ports no longer open.
func (t *Tracker) Update(ports []scanner.Port) {
	seen := make(map[string]bool, len(ports))
	now := t.now()
	for _, p := range ports {
		k := key(p.Port, p.Proto)
		seen[k] = true
		if _, ok := t.entries[k]; !ok {
			t.entries[k] = Entry{Port: p.Port, Proto: p.Proto, FirstSeen: now}
		}
	}
	for k := range t.entries {
		if !seen[k] {
			delete(t.entries, k)
		}
	}
}

// Get returns the Entry for a port, and whether it exists.
func (t *Tracker) Get(port int, proto string) (Entry, bool) {
	e, ok := t.entries[key(port, proto)]
	return e, ok
}

// Save persists the tracker state to disk.
func (t *Tracker) Save() error {
	list := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(t.path, data, 0o644)
}
