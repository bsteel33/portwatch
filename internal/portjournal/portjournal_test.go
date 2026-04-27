package portjournal

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "journal.json")
}

func TestRecord_And_Entries(t *testing.T) {
	j, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	e := Entry{Port: 80, Proto: "tcp", Kind: EventOpened, Service: "http"}
	if err := j.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries := j.Entries()
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 80 || entries[0].Kind != EventOpened {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
	if entries[0].Time.IsZero() {
		t.Error("time should be set automatically")
	}
}

func TestRecord_Persistence(t *testing.T) {
	path := tempPath(t)
	j, _ := New(path)
	_ = j.Record(Entry{Port: 443, Proto: "tcp", Kind: EventOpened})
	_ = j.Record(Entry{Port: 22, Proto: "tcp", Kind: EventClosed})

	j2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := len(j2.Entries()); got != 2 {
		t.Fatalf("want 2 entries after reload, got %d", got)
	}
}

func TestLast_ReturnsRecent(t *testing.T) {
	j, _ := New(tempPath(t))
	for _, port := range []int{80, 443, 22, 8080} {
		_ = j.Record(Entry{Port: port, Proto: "tcp", Kind: EventOpened, Time: time.Now()})
	}
	last := j.Last(2)
	if len(last) != 2 {
		t.Fatalf("want 2, got %d", len(last))
	}
	if last[0].Port != 22 || last[1].Port != 8080 {
		t.Errorf("unexpected last entries: %+v", last)
	}
}

func TestClear_RemovesEntries(t *testing.T) {
	path := tempPath(t)
	j, _ := New(path)
	_ = j.Record(Entry{Port: 80, Proto: "tcp", Kind: EventOpened})
	if err := j.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if got := len(j.Entries()); got != 0 {
		t.Errorf("want 0 after clear, got %d", got)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "null\n" && string(data) != "[]" && string(data) != "null" {
		// accept any valid empty JSON representation
		if len(data) > 10 {
			t.Errorf("file should be empty JSON, got %q", string(data))
		}
	}
}

func TestNew_MissingFile(t *testing.T) {
	j, err := New(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(j.Entries()) != 0 {
		t.Error("expected empty journal for missing file")
	}
}
