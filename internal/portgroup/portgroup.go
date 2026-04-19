// Package portgroup provides grouping of open ports by service category.
package portgroup

import "github.com/user/portwatch/internal/scanner"

// Group represents a named collection of ports.
type Group struct {
	Name  string
	Ports []scanner.Port
}

// Grouper categorizes ports into named groups based on port ranges or known services.
type Grouper struct {
	buckets []bucket
}

type bucket struct {
	name string
	low  int
	high int
}

// New returns a Grouper with default category buckets.
func New() *Grouper {
	return &Grouper{
		buckets: []bucket{
			{"well-known", 1, 1023},
			{"registered", 1024, 49151},
			{"dynamic", 49152, 65535},
		},
	}
}

// Apply groups the given ports into named buckets.
func (g *Grouper) Apply(ports []scanner.Port) []Group {
	index := make(map[string][]scanner.Port)
	for _, p := range ports {
		name := g.classify(p.Port)
		index[name] = append(index[name], p)
	}
	groups := make([]Group, 0, len(g.buckets))
	for _, b := range g.buckets {
		if ps, ok := index[b.name]; ok {
			groups = append(groups, Group{Name: b.name, Ports: ps})
		}
	}
	return groups
}

func (g *Grouper) classify(port int) string {
	for _, b := range g.buckets {
		if port >= b.low && port <= b.high {
			return b.name
		}
	}
	return "unknown"
}
