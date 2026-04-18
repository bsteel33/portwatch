package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/snapshot"
)

func makeDiff(added, removed []snapshot.Port) snapshot.Diff {
	return snapshot.Diff{Added: added, Removed: removed}
}

func TestRecord_And_Load(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	l := audit.New(path)
	diff := makeDiff(
		[]snapshot.Port{{Port: 8080, Proto: "tcp"}},
		nil,
	)
	if err := l.Record("alert", diff); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, err := audit.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Event != "alert" {
		t.Errorf("expected event=alert, got %s", entries[0].Event)
	}
	if len(entries[0].Added) != 1 {
		t.Errorf("expected 1 added port")
	}
}

func TestRecord_DefaultEvent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	l := audit.New(path)
	if err := l.Record("", makeDiff(nil, nil)); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, _ := audit.Load(path)
	if entries[0].Event != "scan" {
		t.Errorf("expected default event=scan, got %s", entries[0].Event)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	entries, err := audit.Load("/nonexistent/audit.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing file")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	l := audit.New(path)
	for i := 0; i < 3; i++ {
		if err := l.Record("scan", makeDiff(nil, nil)); err != nil {
			t.Fatalf("Record %d: %v", i, err)
		}
	}
	entries, _ := audit.Load(path)
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	os.WriteFile(path, []byte{}, 0o644)
	entries, err := audit.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries")
	}
}
