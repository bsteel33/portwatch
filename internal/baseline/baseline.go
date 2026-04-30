// Package baseline manages a trusted baseline of open ports for comparison.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Baseline holds a trusted set of ports captured at a point in time.
type Baseline struct {
	CapturedAt time.Time      `json:"captured_at"`
	Ports      []scanner.Port `json:"ports"`
}

// Manager handles saving and loading baselines.
type Manager struct {
	path string
}

// New returns a Manager using the given file path.
func New(path string) *Manager {
	return &Manager{path: path}
}

// Save writes the baseline to disk.
func (m *Manager) Save(ports []scanner.Port) error {
	b := Baseline{
		CapturedAt: time.Now().UTC(),
		Ports:      ports,
	}
	f, err := os.Create(m.path)
	if err != nil {
		return fmt.Errorf("baseline: create %q: %w", m.path, err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(b); err != nil {
		return fmt.Errorf("baseline: encode: %w", err)
	}
	return nil
}

// Load reads the baseline from disk.
func (m *Manager) Load() (*Baseline, error) {
	f, err := os.Open(m.path)
	if err != nil {
		return nil, fmt.Errorf("baseline: open %q: %w", m.path, err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("baseline: decode %q: %w", m.path, err)
	}
	return &b, nil
}

// Exists reports whether a baseline file is present.
func (m *Manager) Exists() bool {
	_, err := os.Stat(m.path)
	return err == nil
}
