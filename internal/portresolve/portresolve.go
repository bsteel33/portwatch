// Package portresolve maps port numbers to human-readable service names
// with support for custom overrides and protocol-aware lookups.
package portresolve

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Resolver resolves port numbers to service names.
type Resolver struct {
	mu        sync.RWMutex
	overrides map[string]string
	cfg       Config
}

// Result holds the resolved name and its source.
type Result struct {
	Port    int
	Proto   string
	Name    string
	Source  string // "override", "builtin", or "unknown"
}

// New creates a Resolver with the given config.
func New(cfg Config) *Resolver {
	return &Resolver{
		cfg:       cfg,
		overrides: make(map[string]string),
	}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Override registers a custom name for a port/proto pair.
func (r *Resolver) Override(port int, proto, name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.overrides[key(port, proto)] = name
}

// Resolve returns the service name for the given port entry.
func (r *Resolver) Resolve(p scanner.Port) Result {
	r.mu.RLock()
	defer r.mu.RUnlock()

	k := key(p.Port, p.Proto)
	if name, ok := r.overrides[k]; ok {
		return Result{Port: p.Port, Proto: p.Proto, Name: name, Source: "override"}
	}

	if name := scanner.LookupService(p.Port, p.Proto); name != "" {
		return Result{Port: p.Port, Proto: p.Proto, Name: name, Source: "builtin"}
	}

	return Result{Port: p.Port, Proto: p.Proto, Name: fmt.Sprintf("port-%d", p.Port), Source: "unknown"}
}

// ResolveAll resolves a slice of ports and returns results in the same order.
func (r *Resolver) ResolveAll(ports []scanner.Port) []Result {
	out := make([]Result, len(ports))
	for i, p := range ports {
		out[i] = r.Resolve(p)
	}
	return out
}
