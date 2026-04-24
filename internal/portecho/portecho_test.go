package portecho

import (
	"net"
	"strconv"
	"testing"
	"time"
)

// startEchoServer opens a TCP listener that echoes back whatever it receives.
func startEchoServer(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.Error())
	// Use the TCPAddr directly.
	port = ln.Addr().(*net.TCPAddr).Port
	t.Cleanup(func() { ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				buf := make([]byte, 64)
				n, _ := c.Read(buf)
				if n > 0 {
					_, _ = c.Write(buf[:n])
				}
			}(conn)
		}
	}()

	return port
}

func TestCheck_AlivePort(t *testing.T) {
	port := startEchoServer(t)
	checker := New(DefaultConfig())
	res := checker.Check("127.0.0.1", port, "tcp")
	if !res.Alive {
		t.Fatalf("expected alive=true, got err=%v", res.Err)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestCheck_ClosedPort(t *testing.T) {
	checker := New(DefaultConfig())
	res := checker.Check("127.0.0.1", 1, "tcp")
	if res.Alive {
		t.Fatal("expected alive=false for closed port")
	}
	if res.Err == nil {
		t.Error("expected non-nil error for closed port")
	}
}

func TestCheck_Timeout(t *testing.T) {
	// Use a very short timeout against a black-hole address.
	cfg := DefaultConfig()
	cfg.Timeout = 50 * time.Millisecond
	checker := New(cfg)
	start := time.Now()
	res := checker.Check("192.0.2.1", 9999, "tcp") // TEST-NET, should time out
	elapsed := time.Since(start)
	if res.Alive {
		t.Fatal("expected alive=false on timeout")
	}
	if elapsed > 500*time.Millisecond {
		t.Errorf("timeout took too long: %v", elapsed)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Timeout != 2*time.Second {
		t.Errorf("unexpected timeout: %v", cfg.Timeout)
	}
	if len(cfg.Probe) == 0 {
		t.Error("expected non-empty default probe")
	}
	if cfg.MaxReply <= 0 {
		t.Error("expected positive MaxReply")
	}
}
