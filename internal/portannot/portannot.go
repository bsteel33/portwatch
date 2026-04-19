// Package portannot provides free-text annotation storage for individual ports.
package portannot

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type key struct {
	Port  int
	Proto string
}

// Annotation holds a note attached to a port.
type Annotation struct {
	Port  int    `json:"port"`
	Proto string `json:"proto"`
	Note  string `json:"note"`
}

// Annotator stores and retrieves port annotations.
type Annotator struct {
	mu      sync.RWMutex
	path    string
	entries map[key]Annotation
}

// New loads annotations from path, creating an empty store if the file is missing.
func New(path string) (*Annotator, error) {
	a := &Annotator{path: path, entries: make(map[key]Annotation)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return a, nil
	}
	if err != nil {
		return nil, fmt.Errorf("portannot: read %s: %w", path, err)
	}
	var list []Annotation
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("portannot: parse %s: %w", path, err)
	}
	for _, ann := range list {
		a.entries[key{ann.Port, ann.Proto}] = ann
	}
	return a, nil
}

// Set adds or updates the annotation for a port.
func (a *Annotator) Set(port int, proto, note string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries[key{port, proto}] = Annotation{Port: port, Proto: proto, Note: note}
	return a.save()
}

// Get returns the annotation for a port, and whether it exists.
func (a *Annotator) Get(port int, proto string) (Annotation, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	ann, ok := a.entries[key{port, proto}]
	return ann, ok
}

// Remove deletes the annotation for a port.
func (a *Annotator) Remove(port int, proto string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.entries, key{port, proto})
	return a.save()
}

// All returns a copy of all annotations.
func (a *Annotator) All() []Annotation {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Annotation, 0, len(a.entries))
	for _, ann := range a.entries {
		out = append(out, ann)
	}
	return out
}

func (a *Annotator) save() error {
	list := make([]Annotation, 0, len(a.entries))
	for _, ann := range a.entries {
		list = append(list, ann)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.path, data, 0o644)
}
