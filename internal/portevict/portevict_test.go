package portevict

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "evict.json")
}

func TestRecord_And_Load(t *testing.T) {
	p := tempPath(t)
	l, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	if err := l.Record(8080, "tcp", now.Add(-10*time.Minute), now); err != nil {
		t.Fatalf("Record: %v", err)
	}

	// reload from disk
	l2, err := New(p)
	if err != nil {
		t.Fatalf("New reload: %v", err)
	}
	entries := l2.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Port != 8080 || e.Proto != "tcp" {
		t.Errorf("unexpected entry: %+v", e)
	}
	if e.Duration != "10m0s" {
		t.Errorf("unexpected duration: %s", e.Duration)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	p := tempPath(t)
	l, _ := New(p)

	now := time.Now().UTC()
	_ = l.Record(22, "tcp", now.Add(-1*time.Hour), now)
	_ = l.Record(443, "tcp", now.Add(-30*time.Minute), now)

	if got := len(l.Entries()); got != 2 {
		t.Errorf("expected 2 entries, got %d", got)
	}
}

func TestNew_MissingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
	l, err := New(p)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(l.Entries()) != 0 {
		t.Error("expected empty entries for missing file")
	}
}

func TestNew_CorruptFile(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0o644)
	_, err := New(p)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
