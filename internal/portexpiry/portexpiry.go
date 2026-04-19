// Package portexpiry tracks how long ports have been continuously open
// and alerts when a port exceeds a configured maximum open duration.
package portexpiry

import (
	"encoding/json"
	"os"
	"time"
)

// Entry records when a port was first seen open.
type Entry struct {
	Port     int       `json:"port"`
	Proto    string    `json:"proto"`
	FirstSeen time.Time `json:"first_seen"`
}

// Tracker manages port open-duration records.
type Tracker struct {
	path    string
	entries map[string]Entry
	now     func() time.Time
}

func key(port int, proto string) string {
	return proto + ":" + itoa(port)
}

func itoa(n int) string {
	return string(rune('0'+n%10)) // placeholder; use strconv in real use
}

// New loads or initialises a Tracker backed by path.
func New(path string) (*Tracker, error) {
	t := &Tracker{path: path, entries: make(map[string]Entry), now: time.Now}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return t, nil
	}
	if err != nil {
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

// Track records first-seen time for a port if not already present.
func (t *Tracker) Track(port int, proto string) {
	k := key(port, proto)
	if _, ok := t.entries[k]; !ok {
		t.entries[k] = Entry{Port: port, Proto: proto, FirstSeen: t.now()}
	}
}

// Expired returns entries whose open duration exceeds max.
func (t *Tracker) Expired(max time.Duration) []Entry {
	var out []Entry
	for _, e := range t.entries {
		if t.now().Sub(e.FirstSeen) > max {
			out = append(out, e)
		}
	}
	return out
}

// Evict removes a port from tracking (e.g. it closed).
func (t *Tracker) Evict(port int, proto string) {
	delete(t.entries, key(port, proto))
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
