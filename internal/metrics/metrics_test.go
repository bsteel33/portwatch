package metrics

import (
	"testing"
	"time"
)

func TestRecordScan_IncrementsScanCount(t *testing.T) {
	c := New()
	c.RecordScan(10*time.Millisecond, 5, false)
	c.RecordScan(20*time.Millisecond, 3, false)
	s := c.Snapshot()
	if s.ScanCount != 2 {
		t.Fatalf("expected ScanCount 2, got %d", s.ScanCount)
	}
}

func TestRecordScan_AlertCount(t *testing.T) {
	c := New()
	c.RecordScan(5*time.Millisecond, 2, true)
	c.RecordScan(5*time.Millisecond, 2, false)
	c.RecordScan(5*time.Millisecond, 2, true)
	s := c.Snapshot()
	if s.AlertCount != 2 {
		t.Fatalf("expected AlertCount 2, got %d", s.AlertCount)
	}
}

func TestRecordScan_OpenPortsHWM(t *testing.T) {
	c := New()
	c.RecordScan(1*time.Millisecond, 4, false)
	c.RecordScan(1*time.Millisecond, 9, false)
	c.RecordScan(1*time.Millisecond, 2, false)
	s := c.Snapshot()
	if s.OpenPortsHWM != 9 {
		t.Fatalf("expected OpenPortsHWM 9, got %d", s.OpenPortsHWM)
	}
}

func TestRecordScan_LastScanDur(t *testing.T) {
	c := New()
	c.RecordScan(42*time.Millisecond, 1, false)
	s := c.Snapshot()
	if s.LastScanDur != 42*time.Millisecond {
		t.Fatalf("expected 42ms, got %v", s.LastScanDur)
	}
}

func TestReset(t *testing.T) {
	c := New()
	c.RecordScan(10*time.Millisecond, 5, true)
	c.Reset()
	s := c.Snapshot()
	if s.ScanCount != 0 || s.AlertCount != 0 || s.OpenPortsHWM != 0 {
		t.Fatalf("expected zeroed stats after Reset, got %+v", s)
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	c := New()
	c.RecordScan(1*time.Millisecond, 3, false)
	s1 := c.Snapshot()
	c.RecordScan(1*time.Millisecond, 3, false)
	s2 := c.Snapshot()
	if s1.ScanCount == s2.ScanCount {
		t.Fatal("expected snapshot to be independent copy")
	}
}
