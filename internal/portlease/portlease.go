// Package portlease tracks temporary ownership leases on ports,
// allowing a port to be "claimed" for a duration and released.
package portlease

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Lease represents a single port lease entry.
type Lease struct {
	Owner     string    `json:"owner"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Leaser manages port leases.
type Leaser struct {
	mu      sync.Mutex
	leases  map[string]Lease
	path    string
	now     func() time.Time
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// New creates a new Leaser backed by the given file path.
func New(path string) (*Leaser, error) {
	l := &Leaser{
		leases: make(map[string]Lease),
		path:   path,
		now:    time.Now,
	}
	if err := l.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return l, nil
}

// Claim assigns a lease on port/proto to owner for the given duration.
func (l *Leaser) Claim(port int, proto, owner string, ttl time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.leases[key(port, proto)] = Lease{
		Owner:     owner,
		ExpiresAt: l.now().Add(ttl),
	}
	return l.save()
}

// Release removes the lease on port/proto.
func (l *Leaser) Release(port int, proto string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.leases, key(port, proto))
	return l.save()
}

// Get returns the active lease for port/proto, if any.
func (l *Leaser) Get(port int, proto string) (Lease, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	lease, ok := l.leases[key(port, proto)]
	if !ok || l.now().After(lease.ExpiresAt) {
		return Lease{}, false
	}
	return lease, true
}

// Active returns all non-expired leases.
func (l *Leaser) Active() map[string]Lease {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make(map[string]Lease)
	for k, v := range l.leases {
		if !l.now().After(v.ExpiresAt) {
			out[k] = v
		}
	}
	return out
}

func (l *Leaser) save() error {
	data, err := json.Marshal(l.leases)
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}

func (l *Leaser) load() error {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &l.leases)
}
