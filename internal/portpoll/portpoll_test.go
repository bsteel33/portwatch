package portpoll

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func startTCPListener(t *testing.T) (port int, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port, func() { ln.Close() }
}

func TestPoll_OpenPort(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	p := New(DefaultConfig())
	r := p.Poll("127.0.0.1", port, "tcp")

	if !r.Open {
		t.Errorf("expected port %d to be open, got error: %s", port, r.Error)
	}
	if r.Latency <= 0 {
		t.Errorf("expected positive latency, got %v", r.Latency)
	}
}

func TestPoll_ClosedPort(t *testing.T) {
	p := New(DefaultConfig())
	r := p.Poll("127.0.0.1", 1, "tcp") // port 1 is almost certainly closed

	if r.Open {
		t.Errorf("expected port 1 to be closed")
	}
	if r.Error == "" {
		t.Errorf("expected non-empty error for closed port")
	}
}

func TestPoll_Timeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Timeout = 50 * time.Millisecond
	p := New(cfg)

	start := time.Now()
	r := p.Poll("192.0.2.1", 9999, "tcp") // TEST-NET, should timeout
	elapsed := time.Since(start)

	if r.Open {
		t.Errorf("expected closed result for unroutable address")
	}
	if elapsed > 500*time.Millisecond {
		t.Errorf("poll took too long: %v", elapsed)
	}
}

func TestPollAll(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	p := New(DefaultConfig())
	targets := []Target{
		{Port: port, Proto: "tcp"},
		{Port: 1, Proto: "tcp"},
	}
	results := p.PollAll("127.0.0.1", targets)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Open {
		t.Errorf("expected first target (port %d) to be open", port)
	}
	if results[1].Open {
		t.Errorf("expected second target (port 1) to be closed")
	}
	expectedAddr := fmt.Sprintf("127.0.0.1:%d", port)
	if results[0].Addr != expectedAddr {
		t.Errorf("expected addr %s, got %s", expectedAddr, results[0].Addr)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Timeout != 2*time.Second {
		t.Errorf("unexpected default timeout: %v", cfg.Timeout)
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("unexpected default host: %s", cfg.Host)
	}
	if cfg.Proto != "tcp" {
		t.Errorf("unexpected default proto: %s", cfg.Proto)
	}
}
