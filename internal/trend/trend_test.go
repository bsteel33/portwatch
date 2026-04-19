package trend

import (
	"testing"
	"time"
)

func fixedClock(base time.Time, step time.Duration) func() time.Time {
	t := base
	return func() time.Time {
		now := t
		t = t.Add(step)
		return now
	}
}

func TestRecord_And_Points(t *testing.T) {
	tr := New(Config{Window: time.Minute})
	tr.now = fixedClock(time.Now(), time.Second)
	tr.Record(3)
	tr.Record(5)
	pts := tr.Points()
	if len(pts) != 2 {
		t.Fatalf("expected 2 points, got %d", len(pts))
	}
	if pts[1].Count != 5 {
		t.Errorf("expected last count 5, got %d", pts[1].Count)
	}
}

func TestDelta(t *testing.T) {
	tr := New(Config{Window: time.Minute})
	tr.now = fixedClock(time.Now(), time.Second)
	tr.Record(2)
	tr.Record(7)
	if d := tr.Delta(); d != 5 {
		t.Errorf("expected delta 5, got %d", d)
	}
}

func TestDelta_SinglePoint(t *testing.T) {
	tr := New(Config{Window: time.Minute})
	tr.now = fixedClock(time.Now(), time.Second)
	tr.Record(4)
	if d := tr.Delta(); d != 0 {
		t.Errorf("expected delta 0 for single point, got %d", d)
	}
}

func TestPrune_RemovesOldPoints(t *testing.T) {
	base := time.Now()
	tr := New(Config{Window: 5 * time.Second})
	calls := []time.Time{base, base.Add(2 * time.Second), base.Add(10 * time.Second)}
	i := 0
	tr.now = func() time.Time { v := calls[i]; i++; return v }
	tr.Record(1)
	tr.Record(2)
	tr.Record(3) // prunes first two as they are outside 5s window from now=base+10s
	pts := tr.Points()
	if len(pts) != 1 {
		t.Errorf("expected 1 point after prune, got %d", len(pts))
	}
}

func TestReset(t *testing.T) {
	tr := New(Config{Window: time.Minute})
	tr.now = fixedClock(time.Now(), time.Second)
	tr.Record(5)
	tr.Reset()
	if pts := tr.Points(); len(pts) != 0 {
		t.Errorf("expected no points after reset, got %d", len(pts))
	}
}
