// Package portdigest computes and tracks a rolling digest of open port sets,
// allowing callers to detect when the port landscape has changed between scans.
package portdigest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds a digest value and the time it was recorded.
type Entry struct {
	Digest    string
	RecordedAt time.Time
}

// Tracker maintains the most recent digest and detects changes.
type Tracker struct {
	mu      sync.Mutex
	current Entry
	prev    Entry
}

// New returns an empty Tracker.
func New() *Tracker {
	return &Tracker{}
}

// Compute returns a stable hex digest for the given port list.
func Compute(ports []scanner.Port) string {
	keys := make([]string, 0, len(ports))
	for _, p := range ports {
		keys = append(keys, fmt.Sprintf("%s/%d", p.Proto, p.Port))
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		h.Write([]byte(k))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Update records a new digest for the given ports and returns whether it
// differs from the previously recorded digest.
func (t *Tracker) Update(ports []scanner.Port) (changed bool) {
	d := Compute(ports)
	now := time.Now()

	t.mu.Lock()
	defer t.mu.Unlock()

	changed = t.current.Digest != "" && t.current.Digest != d
	t.prev = t.current
	t.current = Entry{Digest: d, RecordedAt: now}
	return changed
}

// Current returns the most recently recorded entry.
func (t *Tracker) Current() Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.current
}

// Previous returns the entry recorded before the last Update.
func (t *Tracker) Previous() Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.prev
}

// Reset clears all recorded state.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current = Entry{}
	t.prev = Entry{}
}
