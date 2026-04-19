// Package portrank ranks open ports by risk level based on well-known
// service associations and configurable weights.
package portrank

import (
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents a risk tier.
type Level int

const (
	Low Level = iota
	Medium
	High
	Critical
)

func (l Level) String() string {
	switch l {
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	case Critical:
		return "critical"
	default:
		return "unknown"
	}
}

// Entry pairs a port with its computed risk level.
type Entry struct {
	Port  scanner.Port
	Level Level
	Score int
}

// Ranker assigns risk levels to ports.
type Ranker struct {
	cfg Config
}

// New returns a Ranker with the given config.
func New(cfg Config) *Ranker {
	return &Ranker{cfg: cfg}
}

// Rank scores and sorts ports from highest to lowest risk.
func (r *Ranker) Rank(ports []scanner.Port) []Entry {
	entries := make([]Entry, 0, len(ports))
	for _, p := range ports {
		score := r.score(p)
		entries = append(entries, Entry{
			Port:  p,
			Level: levelFromScore(score),
			Score: score,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})
	return entries
}

func (r *Ranker) score(p scanner.Port) int {
	s := 0
	if w, ok := r.cfg.Weights[p.Port]; ok {
		s += w
	}
	if p.Port < 1024 {
		s += 10
	}
	if p.Proto == "udp" {
		s += 5
	}
	return s
}

func levelFromScore(score int) Level {
	switch {
	case score >= 80:
		return Critical
	case score >= 50:
		return High
	case score >= 20:
		return Medium
	default:
		return Low
	}
}
