// Package portchain provides a composable pipeline for processing scanned ports
// through a sequence of named stages (filter, label, rank, etc.).
package portchain

import (
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/scanner"
)

// Stage is a single processing step that transforms a slice of ports.
type Stage struct {
	Name string
	Fn   func([]scanner.Port) []scanner.Port
}

// Chain runs ports through an ordered list of stages.
type Chain struct {
	stages  []Stage
	verbose bool
	out     io.Writer
}

// New returns an empty Chain.
func New(verbose bool) *Chain {
	return &Chain{verbose: verbose, out: os.Stderr}
}

// Add appends a stage to the chain.
func (c *Chain) Add(name string, fn func([]scanner.Port) []scanner.Port) *Chain {
	c.stages = append(c.stages, Stage{Name: name, Fn: fn})
	return c
}

// Run executes all stages in order and returns the final port list.
func (c *Chain) Run(ports []scanner.Port) []scanner.Port {
	current := ports
	for _, s := range c.stages {
		before := len(current)
		current = s.Fn(current)
		if c.verbose {
			fmt.Fprintf(c.out, "[chain] stage=%s before=%d after=%d\n", s.Name, before, len(current))
		}
	}
	return current
}

// Len returns the number of registered stages.
func (c *Chain) Len() int { return len(c.stages) }
