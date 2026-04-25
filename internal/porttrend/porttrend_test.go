package porttrend

import (
	"bytes"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func newTracker(window time.Duration, threshold int) *Tracker {
	cfg := Config{Window: window, Threshold: threshold}
	tr := New(cfg)
	return tr
}

func TestRecord_And_Analyze_Stable(t *testing.T) {
	tr := newTracker(time.Minute, 2)
	base := time.Now()
	tr.clock = fixedClock(base)
	tr.Record(5)
	tr.clock = fixedClock(base.Add(10 * time.Second))
	tr.Record(6) // delta 1, below threshold 2

	r := tr.Analyze()
	if r.Direction != DirectionStable {
		t.Fatalf("expected stable, got %s", r.Direction)
	}
	if r.Delta != 1 {
		t.Fatalf("expected delta 1, got %d", r.Delta)
	}
}

func TestRecord_TrendUp(t *testing.T) {
	tr := newTracker(time.Minute, 2)
	base := time.Now()
	tr.clock = fixedClock(base)
	tr.Record(3)
	tr.clock = fixedClock(base.Add(5 * time.Second))
	tr.Record(8) // delta 5 > threshold 2

	r := tr.Analyze()
	if r.Direction != DirectionUp {
		t.Fatalf("expected up, got %s", r.Direction)
	}
}

func TestRecord_TrendDown(t *testing.T) {
	tr := newTracker(time.Minute, 2)
	base := time.Now()
	tr.clock = fixedClock(base)
	tr.Record(10)
	tr.clock = fixedClock(base.Add(5 * time.Second))
	tr.Record(4) // delta -6

	r := tr.Analyze()
	if r.Direction != DirectionDown {
		t.Fatalf("expected down, got %s", r.Direction)
	}
	if r.Delta != -6 {
		t.Fatalf("expected delta -6, got %d", r.Delta)
	}
}

func TestPrune_RemovesOldSamples(t *testing.T) {
	tr := newTracker(30*time.Second, 1)
	base := time.Now()
	tr.clock = fixedClock(base)
	tr.Record(5)
	// advance past window
	tr.clock = fixedClock(base.Add(60 * time.Second))
	tr.Record(5)

	r := tr.Analyze()
	if r.Samples != 1 {
		t.Fatalf("expected 1 sample after pruning, got %d", r.Samples)
	}
}

func TestReset_ClearsSamples(t *testing.T) {
	tr := newTracker(time.Minute, 1)
	tr.Record(5)
	tr.Record(9)
	tr.Reset()
	r := tr.Analyze()
	if r.Samples != 0 {
		t.Fatalf("expected 0 samples after reset, got %d", r.Samples)
	}
}

func TestFprint_Output(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, Result{Direction: DirectionUp, Delta: 4, Samples: 3})
	out := buf.String()
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	for _, want := range []string{"up", "+4", "3"} {
		if !bytes.Contains([]byte(out), []byte(want)) {
			t.Errorf("output %q missing %q", out, want)
		}
	}
}
