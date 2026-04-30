// Package portcluster groups ports into clusters based on shared attributes
// such as service family, port range, or protocol, enabling higher-level
// analysis and reporting of related port activity.
package portcluster

import (
	"fmt"
	"sort"

	"github.com/example/portwatch/internal/scanner"
)

// Cluster represents a named group of ports.
type Cluster struct {
	Name  string
	Ports []scanner.Port
}

// Clusterer groups ports into named clusters using configured rules.
type Clusterer struct {
	cfg Config
}

// New returns a new Clusterer with the given configuration.
func New(cfg Config) *Clusterer {
	return &Clusterer{cfg: cfg}
}

// Apply groups the given ports into clusters and returns them sorted by name.
func (c *Clusterer) Apply(ports []scanner.Port) []Cluster {
	index := make(map[string][]scanner.Port)

	for _, p := range ports {
		name := c.clusterName(p)
		index[name] = append(index[name], p)
	}

	clusters := make([]Cluster, 0, len(index))
	for name, ps := range index {
		clusters = append(clusters, Cluster{Name: name, Ports: ps})
	}

	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Name < clusters[j].Name
	})
	return clusters
}

// clusterName returns the cluster name for a port based on configured strategy.
func (c *Clusterer) clusterName(p scanner.Port) string {
	switch c.cfg.Strategy {
	case StrategyProto:
		return p.Proto
	case StrategyRange:
		return portRange(p.Port, c.cfg.RangeSize)
	default:
		return portRange(p.Port, 1024)
	}
}

func portRange(port int, size int) string {
	if size <= 0 {
		size = 1024
	}
	lo := (port / size) * size
	hi := lo + size - 1
	return fmt.Sprintf("%d-%d", lo, hi)
}
