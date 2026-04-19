// Package portscore aggregates risk signals for an open port into a
// single normalised score in the range [0, 100].
package portscore

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Score holds the computed risk score and contributing factors for one port.
type Score struct {
	Port    scanner.Port
	Total   int
	Factors map[string]int
}

// Scorer computes risk scores for a slice of ports.
type Scorer struct {
	cfg Config
}

// New returns a Scorer using the provided Config.
func New(cfg Config) *Scorer {
	return &Scorer{cfg: cfg}
}

// Score computes a Score for every port in ports.
func (s *Scorer) Score(ports []scanner.Port) []Score {
	results := make([]Score, 0, len(ports))
	for _, p := range ports {
		results = append(results, s.scoreOne(p))
	}
	return results
}

func (s *Scorer) scoreOne(p scanner.Port) Score {
	factors := make(map[string]int)

	// Well-known port bonus.
	if p.Port < 1024 {
		factors["well_known"] = s.cfg.WellKnownWeight
	}

	// Privileged protocol weight.
	if p.Proto == "tcp" {
		factors["tcp"] = s.cfg.TCPWeight
	}

	// Explicitly flagged high-risk ports.
	key := fmt.Sprintf("%d/%s", p.Port, p.Proto)
	if w, ok := s.cfg.HighRiskPorts[key]; ok {
		factors["high_risk"] = w
	}

	total := 0
	for _, v := range factors {
		total += v
	}
	if total > 100 {
		total = 100
	}

	return Score{Port: p, Total: total, Factors: factors}
}
