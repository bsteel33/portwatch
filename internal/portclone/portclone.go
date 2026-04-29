// Package portclone provides utilities for deep-copying port snapshots
// and scan results, allowing safe mutation without affecting the original.
package portclone

import "github.com/user/portwatch/internal/scanner"

// Cloner copies port slices and maps without sharing underlying memory.
type Cloner struct{}

// New returns a new Cloner.
func New() *Cloner {
	return &Cloner{}
}

// Clone returns a deep copy of the given port slice.
func (c *Cloner) Clone(ports []scanner.Port) []scanner.Port {
	if ports == nil {
		return nil
	}
	out := make([]scanner.Port, len(ports))
	copy(out, ports)
	return out
}

// CloneMap returns a deep copy of a map from string keys to port slices.
func (c *Cloner) CloneMap(m map[string][]scanner.Port) map[string][]scanner.Port {
	if m == nil {
		return nil
	}
	out := make(map[string][]scanner.Port, len(m))
	for k, v := range m {
		out[k] = c.Clone(v)
	}
	return out
}

// Merge combines two port slices, deduplicating by (Port, Proto) pair.
// Entries from b override entries from a when there is a conflict.
func (c *Cloner) Merge(a, b []scanner.Port) []scanner.Port {
	type key struct {
		port  int
		proto string
	}
	idx := make(map[key]scanner.Port, len(a)+len(b))
	for _, p := range a {
		idx[key{p.Port, p.Proto}] = p
	}
	for _, p := range b {
		idx[key{p.Port, p.Proto}] = p
	}
	out := make([]scanner.Port, 0, len(idx))
	for _, p := range idx {
		out = append(out, p)
	}
	return out
}
