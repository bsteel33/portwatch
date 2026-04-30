package portlookup

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Result holds the resolved information for a single port.
type Result struct {
	Port    int
	Proto   string
	Service string
	Found   bool
}

// Lookup resolves service names for a slice of open ports.
type Lookup struct {
	mu      sync.RWMutex
	cache   map[string]Result
	resolve func(port int, proto string) (string, bool)
}

// New returns a Lookup using the provided resolver function.
// If resolver is nil, scanner.LookupService is used.
func New(resolver func(port int, proto string) (string, bool)) *Lookup {
	if resolver == nil {
		resolver = func(port int, proto string) (string, bool) {
			svc := scanner.LookupService(port, proto)
			if svc == fmt.Sprintf("%d/%s", port, proto) {
				return svc, false
			}
			return svc, true
		}
	}
	return &Lookup{
		cache:   make(map[string]Result),
		resolve: resolver,
	}
}

// Resolve returns a Result for the given port and protocol.
// Results are cached for the lifetime of the Lookup.
func (l *Lookup) Resolve(port int, proto string) Result {
	k := key(port, proto)

	l.mu.RLock()
	if r, ok := l.cache[k]; ok {
		l.mu.RUnlock()
		return r
	}
	l.mu.RUnlock()

	svc, found := l.resolve(port, proto)
	r := Result{Port: port, Proto: proto, Service: svc, Found: found}

	l.mu.Lock()
	l.cache[k] = r
	l.mu.Unlock()

	return r
}

// ResolveAll resolves service names for all provided ports.
func (l *Lookup) ResolveAll(ports []scanner.Port) []Result {
	out := make([]Result, len(ports))
	for i, p := range ports {
		out[i] = l.Resolve(p.Port, p.Proto)
	}
	return out
}

// Reset clears the internal cache.
func (l *Lookup) Reset() {
	l.mu.Lock()
	l.cache = make(map[string]Result)
	l.mu.Unlock()
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}
