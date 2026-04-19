// Package portttl tracks time-to-live for open ports and marks them
// as expired when they exceed a configured maximum age.
package portttl

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry holds TTL metadata for a single port.
type Entry struct {
	Port     int       `json:"port"`
	Proto    string    `json:"proto"`
	FirstSeen time.Time `json:"first_seen"`
	TTL      time.Duration `json:"ttl"`
}

// Expired reports whether the entry has exceeded its TTL.
func (e Entry) Expired(now time.Time) bool {
	return now.Sub(e.FirstSeen) > e.TTL
}

// Tracker manages TTL entries persisted to disk.
type Tracker struct {
	path    string
	entries map[string]Entry
	now     func() time.Time
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// New loads or initialises a Tracker backed by path.
func New(path string, now func() time.Time) (*Tracker, error) {
	t := &Tracker{path: path, entries: make(map[string]Entry), now: now}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return t, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &t.entries); err != nil {
		return nil, err
	}
	return t, nil
}

// Track registers a port with the given TTL. If already tracked, TTL is updated.
func (t *Tracker) Track(port int, proto string, ttl time.Duration) {
	k := key(port, proto)
	if _, ok := t.entries[k]; !ok {
		t.entries[k] = Entry{Port: port, Proto: proto, FirstSeen: t.now(), TTL: ttl}
		return
	}
	e := t.entries[k]
	e.TTL = ttl
	t.entries[k] = e
}

// Expired returns all entries whose TTL has elapsed.
func (t *Tracker) Expired() []Entry {
	now := t.now()
	var out []Entry
	for _, e := range t.entries {
		if e.Expired(now) {
			out = append(out, e)
		}
	}
	return out
}

// Evict removes the entry for the given port.
func (t *Tracker) Evict(port int, proto string) {
	delete(t.entries, key(port, proto))
}

// Save persists the tracker state to disk.
func (t *Tracker) Save() error {
	data, err := json.MarshalIndent(t.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(t.path, data, 0o600)
}
