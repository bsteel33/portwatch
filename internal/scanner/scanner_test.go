package scanner

import (
	"net"
	"testing"
	"time"
)

// startTestListener opens a TCP listener on a random port and returns the port and a stop func.
func startTestListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScan_OpenPort(t *testing.T) {
	port, stop := startTestListener(t)
	defer stop()

	s := New("127.0.0.1", []int{port}, time.Second)
	result := s.Scan()

	if result.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", result.Host)
	}
	if len(result.Ports) != 1 {
		t.Fatalf("expected 1 port result, got %d", len(result.Ports))
	}
	if !result.Ports[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	s := New("127.0.0.1", []int{1}, time.Second)
	result := s.Scan()

	if len(result.Ports) != 1 {
		t.Fatalf("expected 1 port result, got %d", len(result.Ports))
	}
	if result.Ports[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestOpenPorts(t *testing.T) {
	result := ScanResult{
		Ports: []PortState{
			{Port: 80, Open: true},
			{Port: 81, Open: false},
			{Port: 443, Open: true},
		},
	}
	open := OpenPorts(result)
	if len(open) != 2 {
		t.Errorf("expected 2 open ports, got %d", len(open))
	}
}
