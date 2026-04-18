package ratelimit

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_UnderLimit(t *testing.T) {
	cfg := Config{MaxEvents: 3, Window: time.Minute}
	l := New(cfg)
	base := time.Now()
	l.now = fixedClock(base)

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	cfg := Config{MaxEvents: 2, Window: time.Minute}
	l := New(cfg)
	base := time.Now()
	l.now = fixedClock(base)

	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected Allow() false when limit exceeded")
	}
}

func TestAllow_WindowSlides(t *testing.T) {
	cfg := Config{MaxEvents: 2, Window: time.Second}
	l := New(cfg)
	base := time.Now()
	l.now = fixedClock(base)

	l.Allow()
	l.Allow()

	// Advance past the window
	l.now = fixedClock(base.Add(2 * time.Second))
	if !l.Allow() {
		t.Fatal("expected Allow() true after window slides")
	}
}

func TestReset(t *testing.T) {
	cfg := Config{MaxEvents: 1, Window: time.Minute}
	l := New(cfg)
	l.Allow()
	l.Reset()
	if !l.Allow() {
		t.Fatal("expected Allow() true after Reset")
	}
}

func TestRemaining(t *testing.T) {
	cfg := Config{MaxEvents: 5, Window: time.Minute}
	l := New(cfg)
	base := time.Now()
	l.now = fixedClock(base)

	l.Allow()
	l.Allow()
	if got := l.Remaining(); got != 3 {
		t.Fatalf("expected Remaining()=3, got %d", got)
	}
}
