// Package portpoll provides on-demand polling of a specific port's
// availability, returning latency and reachability status.
package portpoll

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single poll attempt.
type Result struct {
	Port    int
	Proto   string
	Addr    string
	Open    bool
	Latency time.Duration
	Error   string
}

// Poller polls individual ports on demand.
type Poller struct {
	cfg Config
}

// New returns a Poller using the provided Config.
func New(cfg Config) *Poller {
	return &Poller{cfg: cfg}
}

// Poll attempts to connect to the given host:port/proto and returns a Result.
func (p *Poller) Poll(host string, port int, proto string) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()

	conn, err := net.DialTimeout(proto, addr, p.cfg.Timeout)
	latency := time.Since(start)

	if err != nil {
		return Result{
			Port:    port,
			Proto:   proto,
			Addr:    addr,
			Open:    false,
			Latency: latency,
			Error:   err.Error(),
		}
	}
	conn.Close()

	return Result{
		Port:    port,
		Proto:   proto,
		Addr:    addr,
		Open:    true,
		Latency: latency,
	}
}

// PollAll polls each (port, proto) pair in the provided list and returns all results.
func (p *Poller) PollAll(host string, targets []Target) []Result {
	results := make([]Result, 0, len(targets))
	for _, t := range targets {
		results = append(results, p.Poll(host, t.Port, t.Proto))
	}
	return results
}

// Target represents a port/proto pair to poll.
type Target struct {
	Port  int
	Proto string
}
