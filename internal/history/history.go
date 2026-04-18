package history

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry represents a single historical record of a port scan diff.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Added     []snapshot.Port  `json:"added,omitempty"`
	Removed   []snapshot.Port  `json:"removed,omitempty"`
}

// History holds an ordered list of diff entries.
type History struct {
	Entries []Entry `json:"entries"`
	path    string
}

// New returns a History loaded from path, or empty if file does not exist.
func New(path string) (*History, error) {
	h := &History{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, h); err != nil {
		return nil, err
	}
	return h, nil
}

// Record appends a new entry and persists the history file.
func (h *History) Record(added, removed []snapshot.Port) error {
	if len(added) == 0 && len(removed) == 0 {
		return nil
	}
	h.Entries = append(h.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Added:     added,
		Removed:   removed,
	})
	return h.save()
}

// Last returns the most recent entry, or nil if history is empty.
func (h *History) Last() *Entry {
	if len(h.Entries) == 0 {
		return nil
	}
	e := h.Entries[len(h.Entries)-1]
	return &e
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0644)
}
