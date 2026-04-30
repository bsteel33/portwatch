package portlookup

import (
	"bytes"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func stubResolver(port int, proto string) (string, bool) {
	if port == 80 && proto == "tcp" {
		return "http", true
	}
	if port == 443 && proto == "tcp" {
		return "https", true
	}
	return "", false
}

func TestResolve_KnownPort(t *testing.T) {
	l := New(stubResolver)
	r := l.Resolve(80, "tcp")
	if !r.Found {
		t.Fatal("expected Found=true for port 80/tcp")
	}
	if r.Service != "http" {
		t.Fatalf("expected service=http, got %q", r.Service)
	}
}

func TestResolve_UnknownPort(t *testing.T) {
	l := New(stubResolver)
	r := l.Resolve(9999, "tcp")
	if r.Found {
		t.Fatal("expected Found=false for unknown port")
	}
}

func TestResolve_CachesResult(t *testing.T) {
	calls := 0
	l := New(func(port int, proto string) (string, bool) {
		calls++
		return stubResolver(port, proto)
	})
	l.Resolve(80, "tcp")
	l.Resolve(80, "tcp")
	if calls != 1 {
		t.Fatalf("expected 1 resolver call, got %d", calls)
	}
}

func TestReset_ClearsCache(t *testing.T) {
	calls := 0
	l := New(func(port int, proto string) (string, bool) {
		calls++
		return stubResolver(port, proto)
	})
	l.Resolve(80, "tcp")
	l.Reset()
	l.Resolve(80, "tcp")
	if calls != 2 {
		t.Fatalf("expected 2 resolver calls after Reset, got %d", calls)
	}
}

func TestResolveAll_Length(t *testing.T) {
	l := New(stubResolver)
	ports := []scanner.Port{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
		{Port: 22, Proto: "tcp"},
	}
	results := l.ResolveAll(ports)
	if len(results) != len(ports) {
		t.Fatalf("expected %d results, got %d", len(ports), len(results))
	}
}

func TestFprint_NoResults(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, nil)
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output for nil results")
	}
}

func TestFprint_WithResults(t *testing.T) {
	l := New(stubResolver)
	results := l.ResolveAll([]scanner.Port{
		{Port: 80, Proto: "tcp"},
	})
	var buf bytes.Buffer
	Fprint(&buf, results)
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output")
	}
}
