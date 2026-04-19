package portping_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portping"
)

func startEchoServer(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func TestPing_Reachable(t *testing.T) {
	port := startEchoServer(t)
	p := portping.New(portping.DefaultConfig())
	res := p.Ping("127.0.0.1", port, "tcp")
	if !res.Reachable {
		t.Fatal("expected port to be reachable")
	}
	if res.Latency <= 0 {
		t.Errorf("expected positive latency, got %v", res.Latency)
	}
}

func TestPing_Unreachable(t *testing.T) {
	cfg := portping.DefaultConfig()
	cfg.Timeout = 100 * time.Millisecond
	cfg.Attempts = 1
	p := portping.New(cfg)
	res := p.Ping("127.0.0.1", 1, "tcp")
	if res.Reachable {
		t.Fatal("expected port to be unreachable")
	}
	if res.Latency != 0 {
		t.Errorf("expected zero latency for unreachable port, got %v", res.Latency)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := portping.DefaultConfig()
	if cfg.Timeout <= 0 {
		t.Error("expected positive timeout")
	}
	if cfg.Attempts <= 0 {
		t.Error("expected positive attempts")
	}
}
