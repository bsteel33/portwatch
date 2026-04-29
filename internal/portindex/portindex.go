// Package portindex maintains an in-memory index of open ports keyed by
// port number and protocol, enabling fast lookups and set operations across
// multiple scan results.
package portindex

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Index holds a fast-lookup structure for a set of open ports.
type Index struct {
	mu      sync.RWMutex
	entries map[string]scanner.Port
}

// New returns an empty Index.
func New() *Index {
	return &Index{
		entries: make(map[string]scanner.Port),
	}
}

// key returns the canonical map key for a port.
func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Add inserts or replaces a port in the index.
func (idx *Index) Add(p scanner.Port) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.entries[key(p.Port, p.Proto)] = p
}

// Remove deletes a port from the index. It is a no-op if the port is not present.
func (idx *Index) Remove(port int, proto string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.entries, key(port, proto))
}

// Get returns the Port for the given port number and protocol, and whether it
// was found.
func (idx *Index) Get(port int, proto string) (scanner.Port, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	p, ok := idx.entries[key(port, proto)]
	return p, ok
}

// Has reports whether the port/proto pair is present in the index.
func (idx *Index) Has(port int, proto string) bool {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	_, ok := idx.entries[key(port, proto)]
	return ok
}

// All returns a snapshot of all ports currently in the index.
func (idx *Index) All() []scanner.Port {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	out := make([]scanner.Port, 0, len(idx.entries))
	for _, p := range idx.entries {
		out = append(out, p)
	}
	return out
}

// Len returns the number of ports in the index.
func (idx *Index) Len() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.entries)
}

// Reset removes all entries from the index.
func (idx *Index) Reset() {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.entries = make(map[string]scanner.Port)
}

// Rebuild replaces the entire index contents with the provided port list.
func (idx *Index) Rebuild(ports []scanner.Port) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.entries = make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		idx.entries[key(p.Port, p.Proto)] = p
	}
}

// Diff returns ports that are in other but not in idx (added) and ports that
// are in idx but not in other (removed).
func (idx *Index) Diff(other *Index) (added, removed []scanner.Port) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()

	for k, p := range other.entries {
		if _, exists := idx.entries[k]; !exists {
			added = append(added, p)
		}
	}
	for k, p := range idx.entries {
		if _, exists := other.entries[k]; !exists {
			removed = append(removed, p)
		}
	}
	return added, removed
}
