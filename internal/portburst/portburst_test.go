package portburst

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_BelowThreshold(t *testing.T) {
	cfg := Config{Threshold: 5, Window: 10 * time.Second}
	d := New(cfg)
	ev := d.Record(3)
	if ev != nil {
		t.Fatalf("expected nil event below threshold, got %+v", ev)
	}
}

func TestRecord_ExceedsThreshold(t *testing.T) {
	cfg := Config{Threshold: 5, Window: 10 * time.Second}
	d := New(cfg)
	ev := d.Record(6)
	if ev == nil {
		t.Fatal("expected burst event, got nil")
	}
	if ev.Count != 6 {
		t.Errorf("expected count 6, got %d", ev.Count)
	}
	if ev.Threshold != 5 {
		t.Errorf("expected threshold 5, got %d", ev.Threshold)
	}
}

func TestRecord_WindowSlides(t *testing.T) {
	now := time.Now()
	cfg := Config{Threshold: 3, Window: 5 * time.Second}
	d := New(cfg)
	d.clock = fixedClock(now)
	d.Record(2)

	// advance past window
	d.clock = fixedClock(now.Add(6 * time.Second))
	ev := d.Record(2)
	if ev != nil {
		t.Fatalf("old events should have been pruned; got event %+v", ev)
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	cfg := Config{Threshold: 3, Window: 10 * time.Second}
	d := New(cfg)
	d.Record(2)
	d.Reset()
	ev := d.Record(2)
	if ev != nil {
		t.Fatalf("expected nil after reset, got %+v", ev)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Threshold <= 0 {
		t.Error("default threshold should be positive")
	}
	if cfg.Window <= 0 {
		t.Error("default window should be positive")
	}
}
