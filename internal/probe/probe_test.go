package probe_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/probe"
)

func startBannerServer(t *testing.T, banner string) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		fmt.Fprint(conn, banner)
		conn.Close()
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestProbe_ReturnsBanner(t *testing.T) {
	port := startBannerServer(t, "SSH-2.0-OpenSSH_8.9")
	p := probe.New(probe.DefaultConfig())
	res, err := p.Probe("127.0.0.1", port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Banner != "SSH-2.0-OpenSSH_8.9" {
		t.Errorf("banner = %q, want SSH-2.0-OpenSSH_8.9", res.Banner)
	}
	if res.Port != port {
		t.Errorf("port = %d, want %d", res.Port, port)
	}
	if res.Proto != "tcp" {
		t.Errorf("proto = %q, want tcp", res.Proto)
	}
	if res.Latency <= 0 {
		t.Errorf("latency should be positive")
	}
}

func TestProbe_NoBanner(t *testing.T) {
	port := startBannerServer(t, "")
	p := probe.New(probe.DefaultConfig())
	res, err := p.Probe("127.0.0.1", port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Banner != "" {
		t.Errorf("expected empty banner, got %q", res.Banner)
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	cfg := probe.DefaultConfig()
	cfg.Timeout = 300 * time.Millisecond
	p := probe.New(cfg)
	_, err := p.Probe("127.0.0.1", 1)
	if err == nil {
		t.Error("expected error for closed port")
	}
}
