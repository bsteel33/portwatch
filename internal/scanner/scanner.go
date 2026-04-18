package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a scanned port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
	Service  string
}

// ScanResult holds the results of a full host scan.
type ScanResult struct {
	Host      string
	ScannedAt time.Time
	Ports     []PortState
}

// Scanner defines configuration for port scanning.
type Scanner struct {
	Host    string
	Ports   []int
	Timeout time.Duration
}

// New creates a new Scanner with the given host, port list, and timeout.
func New(host string, ports []int, timeout time.Duration) *Scanner {
	return &Scanner{
		Host:    host,
		Ports:   ports,
		Timeout: timeout,
	}
}

// Scan performs a TCP scan on all configured ports and returns a ScanResult.
func (s *Scanner) Scan() ScanResult {
	result := ScanResult{
		Host:      s.Host,
		ScannedAt: time.Now(),
	}

	for _, port := range s.Ports {
		address := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", address, s.Timeout)
		state := PortState{
			Port:     port,
			Protocol: "tcp",
			Open:     err == nil,
		}
		if err == nil {
			conn.Close()
		}
		result.Ports = append(result.Ports, state)
	}

	return result
}

// OpenPorts returns only the open ports from a ScanResult.
func OpenPorts(result ScanResult) []PortState {
	var open []PortState
	for _, p := range result.Ports {
		if p.Open {
			open = append(open, p)
		}
	}
	return open
}
