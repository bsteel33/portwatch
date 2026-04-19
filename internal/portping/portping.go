// Package portping measures round-trip latency to open ports.
package portping

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single ping attempt.
type Result struct {
	Port     int
	Proto    string
	Latency  time.Duration
	Reachable bool
}

// Config holds portping settings.
type Config struct {
	Timeout  time.Duration
	Attempts int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout:  2 * time.Second,
		Attempts: 3,
	}
}

// Pinger pings ports and returns latency results.
type Pinger struct {
	cfg Config
}

// New creates a new Pinger with the given config.
func New(cfg Config) *Pinger {
	return &Pinger{cfg: cfg}
}

// Ping attempts to connect to the given port and returns the average latency.
func (p *Pinger) Ping(host string, port int, proto string) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	var total time.Duration
	success := 0
	for i := 0; i < p.cfg.Attempts; i++ {
		start := time.Now()
		conn, err := net.DialTimeout(proto, addr, p.cfg.Timeout)
		elapsed := time.Since(start)
		if err != nil {
			continue
		}
		conn.Close()
		total += elapsed
		success++
	}
	if success == 0 {
		return Result{Port: port, Proto: proto, Reachable: false}
	}
	return Result{
		Port:      port,
		Proto:     proto,
		Latency:   total / time.Duration(success),
		Reachable: true,
	}
}
