// Package portmemo stores arbitrary key-value memos attached to a port.
package portmemo

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Memo holds a collection of key-value notes for a single port entry.
type Memo struct {
	Port  int               `json:"port"`
	Proto string            `json:"proto"`
	Notes map[string]string `json:"notes"`
}

// Store manages memos persisted to a JSON file.
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]*Memo // key: "port/proto"
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// New loads an existing memo store from path, or creates an empty one.
func New(path string) (*Store, error) {
	s := &Store{path: path, data: make(map[string]*Memo)}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	var entries []*Memo
	if err := json.Unmarshal(b, &entries); err != nil {
		return nil, err
	}
	for _, m := range entries {
		s.data[key(m.Port, m.Proto)] = m
	}
	return s, nil
}

// Set adds or updates a note for the given port/proto.
func (s *Store) Set(port int, proto, k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[key(port, proto)]
	if !ok {
		entry = &Memo{Port: port, Proto: proto, Notes: make(map[string]string)}
		s.data[key(port, proto)] = entry
	}
	entry.Notes[k] = v
	return s.save()
}

// Get retrieves a note value; returns "", false if absent.
func (s *Store) Get(port int, proto, k string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.data[key(port, proto)]
	if !ok {
		return "", false
	}
	v, found := entry.Notes[k]
	return v, found
}

// Remove deletes a single note key from the given port/proto.
func (s *Store) Remove(port int, proto, k string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[key(port, proto)]
	if !ok {
		return nil
	}
	delete(entry.Notes, k)
	if len(entry.Notes) == 0 {
		delete(s.data, key(port, proto))
	}
	return s.save()
}

// All returns a copy of all stored memos.
func (s *Store) All() []*Memo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Memo, 0, len(s.data))
	for _, m := range s.data {
		out = append(out, m)
	}
	return out
}

func (s *Store) save() error {
	entries := make([]*Memo, 0, len(s.data))
	for _, m := range s.data {
		entries = append(entries, m)
	}
	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}
