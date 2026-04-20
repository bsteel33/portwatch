package portcache

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func fakePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Proto: "tcp", Service: "http"},
		{Port: 443, Proto: "tcp", Service: "https"},
	}
}

func TestGet_EmptyCache(t *testing.T) {
	c := New(DefaultConfig())
	_, ok := c.Get()
	if ok {
		t.Fatal("expected cache miss on empty cache")
	}
}

func TestSet_And_Get(t *testing.T) {
	c := New(DefaultConfig())
	ports := fakePorts()
	c.Set(ports)
	got, ok := c.Get()
	if !ok {
		t.Fatal("expected cache hit after Set")
	}
	if len(got) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestGet_Expired(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TTL = 50 * time.Millisecond
	c := New(cfg)
	c.Set(fakePorts())

	time.Sleep(80 * time.Millisecond)

	_, ok := c.Get()
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate(t *testing.T) {
	c := New(DefaultConfig())
	c.Set(fakePorts())
	c.Invalidate()
	_, ok := c.Get()
	if ok {
		t.Fatal("expected cache miss after Invalidate")
	}
}

func TestSet_Overwrites(t *testing.T) {
	c := New(DefaultConfig())
	c.Set(fakePorts())
	newPorts := []scanner.Port{{Port: 22, Proto: "tcp", Service: "ssh"}}
	c.Set(newPorts)
	got, ok := c.Get()
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got) != 1 || got[0].Port != 22 {
		t.Fatalf("expected overwritten value, got %+v", got)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TTL <= 0 {
		t.Fatalf("expected positive default TTL, got %v", cfg.TTL)
	}
}
