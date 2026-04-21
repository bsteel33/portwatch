package portdebounce

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func newTestDebouncer(window time.Duration) *Debouncer {
	d := New(Config{Window: window})
	d.now = func() time.Time { return fixedNow }
	return d
}

func TestObserve_FirstCall_ReturnsFalse(t *testing.T) {
	d := newTestDebouncer(5 * time.Second)
	if d.Observe(80, "tcp") {
		t.Fatal("expected false on first observation")
	}
	if d.Pending() != 1 {
		t.Fatalf("expected 1 pending, got %d", d.Pending())
	}
}

func TestObserve_WithinWindow_ReturnsFalse(t *testing.T) {
	d := New(Config{Window: 10 * time.Second})
	t0 := fixedNow
	d.now = func() time.Time { return t0 }
	d.Observe(443, "tcp")

	d.now = func() time.Time { return t0.Add(5 * time.Second) }
	if d.Observe(443, "tcp") {
		t.Fatal("expected false within window")
	}
}

func TestObserve_AfterWindow_ReturnsTrue(t *testing.T) {
	d := New(Config{Window: 10 * time.Second})
	t0 := fixedNow
	d.now = func() time.Time { return t0 }
	d.Observe(22, "tcp")

	d.now = func() time.Time { return t0.Add(10 * time.Second) }
	if !d.Observe(22, "tcp") {
		t.Fatal("expected true after window elapses")
	}
	if d.Pending() != 0 {
		t.Fatalf("expected 0 pending after confirmation, got %d", d.Pending())
	}
}

func TestDismiss_RemovesPending(t *testing.T) {
	d := newTestDebouncer(5 * time.Second)
	d.Observe(8080, "tcp")
	d.Dismiss(8080, "tcp")
	if d.Pending() != 0 {
		t.Fatalf("expected 0 pending after dismiss, got %d", d.Pending())
	}
}

func TestReset_ClearsAll(t *testing.T) {
	d := newTestDebouncer(5 * time.Second)
	d.Observe(80, "tcp")
	d.Observe(443, "tcp")
	d.Reset()
	if d.Pending() != 0 {
		t.Fatalf("expected 0 after reset, got %d", d.Pending())
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Window <= 0 {
		t.Fatal("expected positive default window")
	}
}
