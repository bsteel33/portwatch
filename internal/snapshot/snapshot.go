package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// PortInfo holds information about a single open port.
type PortInfo struct {
	Port    int    `json:"port"`
	Service string `json:"service"`
	Proto   string `json:"proto"`
}

// Snapshot represents the state of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time  `json:"timestamp"`
	Ports     []PortInfo `json:"ports"`
}

// New creates a new Snapshot with the current timestamp.
func New(ports []PortInfo) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now(),
		Ports:     ports,
	}
}

// Save writes the snapshot to a JSON file at the given path.
func (s *Snapshot) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}
