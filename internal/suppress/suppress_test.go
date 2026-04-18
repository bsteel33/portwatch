package suppress

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func newTestSuppressor(ttl time.Duration) *Suppressor {
	s := New(Config{TTL: ttl})
	return s
}

func TestSuppress_And_IsSuppressed(t *testing.T) {
	s := newTestSuppressor(5 * time.Minute)
	base := time.Now()
	s.now = fixedNow(base)

	if s.IsSuppressed("tcp:80") {
		t.Fatal("expected not suppressed before Suppress call")
	}
	s.Suppress("tcp:80")
	if !s.IsSuppressed("tcp:80") {
		t.Fatal("expected suppressed after Suppress call")
	}
}

func TestSuppress_Expiry(t *testing.T) {
	s := newTestSuppressor(5 * time.Minute)
	base := time.Now()
	s.now = fixedNow(base)
	s.Suppress("tcp:443")

	// advance past TTL
	s.now = fixedNow(base.Add(6 * time.Minute))
	if s.IsSuppressed("tcp:443") {
		t.Fatal("expected suppression to have expired")
	}
}

func TestReset(t *testing.T) {
	s := newTestSuppressor(5 * time.Minute)
	s.Suppress("tcp:22")
	s.Reset("tcp:22")
	if s.IsSuppressed("tcp:22") {
		t.Fatal("expected key to be cleared after Reset")
	}
}

func TestLen(t *testing.T) {
	s := newTestSuppressor(5 * time.Minute)
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	s.Suppress("tcp:80")
	s.Suppress("udp:53")
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TTL != 10*time.Minute {
		t.Fatalf("unexpected default TTL: %v", cfg.TTL)
	}
}
