// Package audit logs scan events and alerts to a persistent audit trail.
package audit

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Event     string           `json:"event"`
	Added     []snapshot.Port  `json:"added,omitempty"`
	Removed   []snapshot.Port  `json:"removed,omitempty"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
}

// New returns a Logger that appends to path.
func New(path string) *Logger {
	return &Logger{path: path}
}

// Record appends an entry to the audit log.
func (l *Logger) Record(event string, diff snapshot.Diff) error {
	if event == "" {
		event = "scan"
	}
	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Added:     diff.Added,
		Removed:   diff.Removed,
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(e)
}

// Load reads all audit entries from path.
func Load(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
