// Package portmap maintains a human-readable mapping of ports to host-defined
// service names, allowing operators to annotate ports with custom labels.
package portmap

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Entry associates a port/proto pair with a custom service name.
type Entry struct {
	Port  int    `json:"port"`
	Proto string `json:"proto"`
	Name  string `json:"name"`
}

// Map holds the in-memory port-to-name mapping.
type Map struct {
	mu      sync.RWMutex
	entries map[string]Entry
	path    string
}

func key(port int, proto string) string {
	return proto + ":" + itoa(port)
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

// New loads a Map from path, or returns an empty Map if the file does not exist.
func New(path string) (*Map, error) {
	m := &Map{path: path, entries: make(map[string]Entry)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return m, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, e := range list {
		m.entries[key(e.Port, e.Proto)] = e
	}
	return m, nil
}

// Set adds or updates a mapping.
func (m *Map) Set(port int, proto, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[key(port, proto)] = Entry{Port: port, Proto: proto, Name: name}
}

// Get returns the custom name for a port, and whether it was found.
func (m *Map) Get(port int, proto string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[key(port, proto)]
	return e.Name, ok
}

// Remove deletes a mapping.
func (m *Map) Remove(port int, proto string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, key(port, proto))
}

// Save persists the map to disk.
func (m *Map) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0o644)
}
