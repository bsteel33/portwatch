package portpause

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func newTestPauser(t time.Time) *Pauser {
	p := New()
	p.now = fixedNow(t)
	return p
}

func TestPause_And_IsPaused(t *testing.T) {
	p := newTestPauser(epoch)
	p.Pause(80, "tcp", 10*time.Minute)
	if !p.IsPaused(80, "tcp") {
		t.Fatal("expected port to be paused")
	}
}

func TestIsPaused_Expired(t *testing.T) {
	p := newTestPauser(epoch)
	p.Pause(80, "tcp", 5*time.Minute)
	// advance clock past expiry
	p.now = fixedNow(epoch.Add(10 * time.Minute))
	if p.IsPaused(80, "tcp") {
		t.Fatal("expected pause to have expired")
	}
}

func TestResume_ClearsPause(t *testing.T) {
	p := newTestPauser(epoch)
	p.Pause(443, "tcp", 1*time.Hour)
	p.Resume(443, "tcp")
	if p.IsPaused(443, "tcp") {
		t.Fatal("expected port to be resumed")
	}
}

func TestActive_ReturnsNonExpired(t *testing.T) {
	p := newTestPauser(epoch)
	p.Pause(22, "tcp", 10*time.Minute)
	p.Pause(53, "udp", 10*time.Minute)
	p.Pause(8080, "tcp", 1*time.Minute)

	// advance so 8080 expires
	p.now = fixedNow(epoch.Add(5 * time.Minute))

	actives := p.Active()
	if len(actives) != 2 {
		t.Fatalf("expected 2 active pauses, got %d", len(actives))
	}
}

func TestReset_ClearsAll(t *testing.T) {
	p := newTestPauser(epoch)
	p.Pause(80, "tcp", 1*time.Hour)
	p.Pause(443, "tcp", 1*time.Hour)
	p.Reset()
	if len(p.entries) != 0 {
		t.Fatal("expected all entries cleared after reset")
	}
}

func TestIsPaused_ProtoDistinct(t *testing.T) {
	p := newTestPauser(epoch)
	p.Pause(53, "tcp", 1*time.Hour)
	if p.IsPaused(53, "udp") {
		t.Fatal("udp should not be paused when only tcp is paused")
	}
}
