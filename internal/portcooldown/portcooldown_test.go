package portcooldown

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func newTestCooldown(window time.Duration) *Cooldown {
	c := New(window)
	t := fixedNow
	c.now = func() time.Time { return t }
	return c
}

func TestMark_And_IsCooling(t *testing.T) {
	c := newTestCooldown(10 * time.Second)
	c.Mark(80, "tcp")
	if !c.IsCooling(80, "tcp") {
		t.Fatal("expected port to be cooling after Mark")
	}
}

func TestIsCooling_NotMarked(t *testing.T) {
	c := newTestCooldown(10 * time.Second)
	if c.IsCooling(443, "tcp") {
		t.Fatal("expected port not to be cooling before Mark")
	}
}

func TestIsCooling_Expired(t *testing.T) {
	c := newTestCooldown(5 * time.Second)
	c.Mark(22, "tcp")
	// advance clock past window
	advanced := fixedNow.Add(10 * time.Second)
	c.now = func() time.Time { return advanced }
	if c.IsCooling(22, "tcp") {
		t.Fatal("expected cooldown to have expired")
	}
}

func TestReset_ClearsEntry(t *testing.T) {
	c := newTestCooldown(30 * time.Second)
	c.Mark(8080, "tcp")
	c.Reset(8080, "tcp")
	if c.IsCooling(8080, "tcp") {
		t.Fatal("expected cooldown to be cleared after Reset")
	}
}

func TestActive_ReturnsOnlyCooling(t *testing.T) {
	c := newTestCooldown(10 * time.Second)
	c.Mark(80, "tcp")
	c.Mark(443, "tcp")

	// expire port 80 by advancing clock for its check
	origNow := fixedNow
	c.now = func() time.Time { return origNow.Add(15 * time.Second) }
	// re-mark 443 so it's fresh at the advanced time
	c.Mark(443, "tcp")

	active := c.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active cooldown, got %d", len(active))
	}
	if active[0].Port != 443 {
		t.Errorf("expected port 443, got %d", active[0].Port)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Window != 5*time.Minute {
		t.Errorf("expected 5m window, got %v", cfg.Window)
	}
}

func TestApplyFlags_Override(t *testing.T) {
	dst := DefaultConfig()
	src := Config{Window: 2 * time.Minute}
	ApplyFlags(&dst, src)
	if dst.Window != 2*time.Minute {
		t.Errorf("expected 2m, got %v", dst.Window)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	dst := DefaultConfig()
	ApplyFlags(&dst, Config{})
	if dst.Window != 5*time.Minute {
		t.Errorf("expected default 5m, got %v", dst.Window)
	}
}
