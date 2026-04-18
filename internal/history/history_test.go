package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func samplePorts() ([]snapshot.Port, []snapshot.Port) {
	added := []snapshot.Port{{Port: 8080, Proto: "tcp", Service: "http-alt"}}
	removed := []snapshot.Port{{Port: 22, Proto: "tcp", Service: "ssh"}}
	return added, removed
}

func TestRecord_And_Load(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	added, removed := samplePorts()
	if err := h.Record(added, removed); err != nil {
		t.Fatalf("Record: %v", err)
	}

	h2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(h2.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h2.Entries))
	}
	if len(h2.Entries[0].Added) != 1 || h2.Entries[0].Added[0].Port != 8080 {
		t.Errorf("unexpected added ports: %v", h2.Entries[0].Added)
	}
}

func TestRecord_NoChanges_SkipsEntry(t *testing.T) {
	dir := t.TempDir()
	h, _ := New(filepath.Join(dir, "history.json"))
	if err := h.Record(nil, nil); err != nil {
		t.Fatal(err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(h.Entries))
	}
}

func TestLast_Empty(t *testing.T) {
	dir := t.TempDir()
	h, _ := New(filepath.Join(dir, "history.json"))
	if h.Last() != nil {
		t.Error("expected nil for empty history")
	}
}

func TestNew_MissingFile(t *testing.T) {
	h, err := New(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil || len(h.Entries) != 0 {
		t.Error("expected empty history")
	}
}

func TestNew_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := New(path)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
