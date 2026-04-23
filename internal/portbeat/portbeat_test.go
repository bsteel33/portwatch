package portbeat

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

var fixedNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func newTracker(stale time.Duration, now time.Time) *Tracker {
	t := New(stale)
	t.now = func() time.Time { return now }
	return t
}

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
	}
}

func TestPulse_RecordsLastSeen(t *testing.T) {
	tr := newTracker(5*time.Minute, fixedNow)
	tr.Pulse(samplePorts())

	b, ok := tr.Get(22, "tcp")
	if !ok {
		t.Fatal("expected beat for port 22")
	}
	if !b.LastSeen.Equal(fixedNow) {
		t.Errorf("expected LastSeen=%v, got %v", fixedNow, b.LastSeen)
	}
}

func TestGet_Missing(t *testing.T) {
	tr := newTracker(5*time.Minute, fixedNow)
	_, ok := tr.Get(9999, "tcp")
	if ok {
		t.Error("expected missing port to return false")
	}
}

func TestStale_NoneExpired(t *testing.T) {
	tr := newTracker(5*time.Minute, fixedNow)
	tr.Pulse(samplePorts())
	// advance time by 1 minute — still within window
	tr.now = func() time.Time { return fixedNow.Add(1 * time.Minute) }

	stale := tr.Stale()
	if len(stale) != 0 {
		t.Errorf("expected 0 stale ports, got %d", len(stale))
	}
}

func TestStale_AfterWindow(t *testing.T) {
	tr := newTracker(5*time.Minute, fixedNow)
	tr.Pulse(samplePorts())
	// advance time past the stale window
	tr.now = func() time.Time { return fixedNow.Add(10 * time.Minute) }

	stale := tr.Stale()
	if len(stale) != 2 {
		t.Errorf("expected 2 stale ports, got %d", len(stale))
	}
}

func TestReset_ClearsBeats(t *testing.T) {
	tr := newTracker(5*time.Minute, fixedNow)
	tr.Pulse(samplePorts())
	tr.Reset()

	_, ok := tr.Get(22, "tcp")
	if ok {
		t.Error("expected tracker to be empty after reset")
	}
	if s := tr.Stale(); len(s) != 0 {
		t.Errorf("expected no stale entries after reset, got %d", len(s))
	}
}
