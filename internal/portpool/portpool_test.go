package portpool

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int, proto string) scanner.Port {
	return scanner.Port{Port: port, Proto: proto, Service: ""}
}

func TestAdd_And_Contains(t *testing.T) {
	p := New("test", 0)
	port := makePort(80, "tcp")
	if err := p.Add(port); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.Contains(port) {
		t.Fatal("expected pool to contain port 80/tcp")
	}
}

func TestAdd_Duplicate_IsNoop(t *testing.T) {
	p := New("test", 0)
	port := makePort(443, "tcp")
	_ = p.Add(port)
	_ = p.Add(port)
	if p.Len() != 1 {
		t.Fatalf("expected len 1, got %d", p.Len())
	}
}

func TestAdd_CapacityExceeded(t *testing.T) {
	p := New("limited", 2)
	_ = p.Add(makePort(80, "tcp"))
	_ = p.Add(makePort(443, "tcp"))
	err := p.Add(makePort(8080, "tcp"))
	if err == nil {
		t.Fatal("expected capacity error, got nil")
	}
	if p.Len() != 2 {
		t.Fatalf("expected len 2, got %d", p.Len())
	}
}

func TestRemove(t *testing.T) {
	p := New("test", 0)
	port := makePort(22, "tcp")
	_ = p.Add(port)
	p.Remove(port)
	if p.Contains(port) {
		t.Fatal("expected port to be removed")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	p := New("test", 0)
	_ = p.Add(makePort(80, "tcp"))
	_ = p.Add(makePort(443, "tcp"))
	all := p.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(all))
	}
}

func TestFprint_ContainsName(t *testing.T) {
	p := New("mypool", 5)
	_ = p.Add(scanner.Port{Port: 80, Proto: "tcp", Service: "http"})
	var buf bytes.Buffer
	Fprint(&buf, p)
	out := buf.String()
	if !strings.Contains(out, "mypool") {
		t.Errorf("expected pool name in output, got: %s", out)
	}
	if !strings.Contains(out, "http") {
		t.Errorf("expected service name in output, got: %s", out)
	}
}
