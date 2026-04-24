// Package portecho provides a lightweight echo/heartbeat checker that
// verifies a port responds correctly by sending a probe byte and expecting
// a non-empty reply within a configurable timeout.
package portecho

import (
	"fmt"
	"net"
	"time"
)

// Config holds configuration for the echo checker.
type Config struct {
	Timeout  time.Duration
	Probe    []byte
	MaxReply int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout:  2 * time.Second,
		Probe:    []byte("\n"),
		MaxReply: 256,
	}
}

// Result holds the outcome of a single echo check.
type Result struct {
	Port    int
	Proto   string
	Alive   bool
	Latency time.Duration
	Err     error
}

// Checker performs echo checks against ports.
type Checker struct {
	cfg Config
}

// New returns a new Checker using the given Config.
func New(cfg Config) *Checker {
	return &Checker{cfg: cfg}
}

// Check dials the given host:port over proto ("tcp") and sends the probe
// payload, returning a Result indicating whether the port echoed back data.
func (c *Checker) Check(host string, port int, proto string) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()

	conn, err := net.DialTimeout(proto, addr, c.cfg.Timeout)
	if err != nil {
		return Result{Port: port, Proto: proto, Alive: false, Err: err}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(c.cfg.Timeout))

	if len(c.cfg.Probe) > 0 {
		if _, err := conn.Write(c.cfg.Probe); err != nil {
			return Result{Port: port, Proto: proto, Alive: false, Err: err}
		}
	}

	buf := make([]byte, c.cfg.MaxReply)
	n, err := conn.Read(buf)
	latency := time.Since(start)

	if err != nil || n == 0 {
		return Result{Port: port, Proto: proto, Alive: false, Latency: latency, Err: err}
	}

	return Result{Port: port, Proto: proto, Alive: true, Latency: latency}
}
