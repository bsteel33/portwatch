package portversion

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "portversion.json")
}

func TestUpdate_FirstCall_NoChange(t *testing.T) {
	tr, _ := New(tempPath(t))
	ch := tr.Update(80, "tcp", "nginx/1.24")
	if ch != nil {
		t.Fatalf("expected nil change on first update, got %+v", ch)
	}
}

func TestUpdate_SameVersion_NoChange(t *testing.T) {
	tr, _ := New(tempPath(t))
	tr.Update(80, "tcp", "nginx/1.24")
	ch := tr.Update(80, "tcp", "nginx/1.24")
	if ch != nil {
		t.Fatalf("expected nil change for same version, got %+v", ch)
	}
}

func TestUpdate_VersionChanged(t *testing.T) {
	tr, _ := New(tempPath(t))
	tr.Update(80, "tcp", "nginx/1.24")
	ch := tr.Update(80, "tcp", "nginx/1.25")
	if ch == nil {
		t.Fatal("expected a Change, got nil")
	}
	if ch.OldVersion != "nginx/1.24" || ch.NewVersion != "nginx/1.25" {
		t.Fatalf("unexpected change values: %+v", ch)
	}
}

func TestGet_Missing(t *testing.T) {
	tr, _ := New(tempPath(t))
	_, ok := tr.Get(443, "tcp")
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestSaveAndLoad(t *testing.T) {
	p := tempPath(t)
	tr, _ := New(p)
	tr.Update(22, "tcp", "OpenSSH_9.3")
	if err := tr.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	tr2, err := New(p)
	if err != nil {
		t.Fatalf("New after save: %v", err)
	}
	e, ok := tr2.Get(22, "tcp")
	if !ok {
		t.Fatal("entry not found after reload")
	}
	if e.Version != "OpenSSH_9.3" {
		t.Fatalf("expected OpenSSH_9.3, got %q", e.Version)
	}
}

func TestNew_MissingFile(t *testing.T) {
	tr, err := New(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestNew_CorruptFile(t *testing.T) {
	p := tempPath(t)
	os.WriteFile(p, []byte("not json{"), 0o644)
	_, err := New(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	tr, _ := New(tempPath(t))
	tr.Update(80, "tcp", "apache")
	tr.Reset()
	_, ok := tr.Get(80, "tcp")
	if ok {
		t.Fatal("expected entry to be cleared after Reset")
	}
}

func TestSave_WritesValidJSON(t *testing.T) {
	p := tempPath(t)
	tr, _ := New(p)
	tr.Update(8080, "tcp", "Jetty/11")
	tr.Save()

	raw, _ := os.ReadFile(p)
	var out map[string]Entry
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
}
