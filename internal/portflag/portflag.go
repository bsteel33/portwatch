// Package portflag provides a mechanism to attach boolean flags to individual
// ports, enabling lightweight per-port state markers such as "watched",
// "ignored", or "reviewed".
package portflag

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Flagger stores named flags keyed by port+proto.
type Flagger struct {
	mu    sync.RWMutex
	flags map[string]map[string]bool // key -> flagName -> set
	path  string
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// New creates a Flagger, loading persisted state from path if it exists.
func New(path string) (*Flagger, error) {
	f := &Flagger{path: path, flags: make(map[string]map[string]bool)}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return f, nil
		}
		return nil, fmt.Errorf("portflag: read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &f.flags); err != nil {
		return nil, fmt.Errorf("portflag: parse %s: %w", path, err)
	}
	return f, nil
}

// Set adds flagName to the given port/proto.
func (f *Flagger) Set(port int, proto, flagName string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	k := key(port, proto)
	if f.flags[k] == nil {
		f.flags[k] = make(map[string]bool)
	}
	f.flags[k][flagName] = true
	return f.save()
}

// Unset removes flagName from the given port/proto.
func (f *Flagger) Unset(port int, proto, flagName string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	k := key(port, proto)
	delete(f.flags[k], flagName)
	if len(f.flags[k]) == 0 {
		delete(f.flags, k)
	}
	return f.save()
}

// Has reports whether flagName is set for the given port/proto.
func (f *Flagger) Has(port int, proto, flagName string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.flags[key(port, proto)][flagName]
}

// Flags returns all flag names set for the given port/proto.
func (f *Flagger) Flags(port int, proto string) []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	var out []string
	for name := range f.flags[key(port, proto)] {
		out = append(out, name)
	}
	return out
}

func (f *Flagger) save() error {
	data, err := json.MarshalIndent(f.flags, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, data, 0o644)
}
