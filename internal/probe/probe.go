// Package probe performs banner grabbing and service fingerprinting on open ports.
package probe

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Result holds the outcome of probing a single port.
type Result struct {
	Port    int
	Proto   string
	Banner  string
	Latency time.Duration
}

// Config controls probe behaviour.
type Config struct {
	Timeout time.Duration
	MaxRead int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout: 2 * time.Second,
		MaxRead: 256,
	}
}

// Prober grabs banners from TCP ports.
type Prober struct {
	cfg Config
}

// New creates a Prober with the given config.
func New(cfg Config) *Prober {
	return &Prober{cfg: cfg}
}

// Probe connects to host:port, reads any initial banner, and returns a Result.
func (p *Prober) Probe(host string, port int) (Result, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", addr, p.cfg.Timeout)
	if err != nil {
		return Result{}, fmt.Errorf("probe %s: %w", addr, err)
	}
	defer conn.Close()

	latency := time.Since(start)
	_ = conn.SetReadDeadline(time.Now().Add(p.cfg.Timeout))

	buf := make([]byte, p.cfg.MaxRead)
	n, _ := conn.Read(buf)
	banner := strings.TrimSpace(string(buf[:n]))

	return Result{
		Port:    port,
		Proto:   "tcp",
		Banner:  banner,
		Latency: latency,
	}, nil
}
