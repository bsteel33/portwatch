// Package portlock tracks ports that have been explicitly locked (allowlisted)
// and flags any scan result containing unlocked ports.
package portlock

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Lock represents a single locked port entry.
type Lock struct {
	Port  int    `json:"port"`
	Proto string `json:"proto"`
	Note  string `json:"note,omitempty"`
}

// Locker holds the set of locked ports.
type Locker struct {
	mu    sync.RWMutex
	locks map[string]Lock
	path  string
}

func key(port int, proto string) string {
	return proto + ":" + itoa(port)
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

// New loads a Locker from path, creating an empty one if missing.
func New(path string) (*Locker, error) {
	l := &Locker{path: path, locks: make(map[string]Lock)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Lock
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	for _, e := range entries {
		l.locks[key(e.Port, e.Proto)] = e
	}
	return l, nil
}

// Add adds a lock entry and persists.
func (l *Locker) Add(port int, proto, note string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.locks[key(port, proto)] = Lock{Port: port, Proto: proto, Note: note}
	return l.save()
}

// Remove removes a lock entry and persists.
func (l *Locker) Remove(port int, proto string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.locks, key(port, proto))
	return l.save()
}

// Unlocked returns ports from ps that are not locked.
func (l *Locker) Unlocked(ps []scanner.Port) []scanner.Port {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []scanner.Port
	for _, p := range ps {
		if _, ok := l.locks[key(p.Port, p.Proto)]; !ok {
			out = append(out, p)
		}
	}
	return out
}

func (l *Locker) save() error {
	entries := make([]Lock, 0, len(l.locks))
	for _, v := range l.locks {
		entries = append(entries, v)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}
