package portmemo

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "memos.json")
}

func TestSet_And_Get(t *testing.T) {
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set(80, "tcp", "owner", "webteam"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	v, ok := s.Get(80, "tcp", "owner")
	if !ok || v != "webteam" {
		t.Errorf("Get = %q, %v; want \"webteam\", true", v, ok)
	}
}

func TestGet_Missing(t *testing.T) {
	s, _ := New(tempPath(t))
	v, ok := s.Get(443, "tcp", "owner")
	if ok || v != "" {
		t.Errorf("expected missing, got %q %v", v, ok)
	}
}

func TestRemove(t *testing.T) {
	p := tempPath(t)
	s, _ := New(p)
	_ = s.Set(22, "tcp", "note", "ssh")
	_ = s.Remove(22, "tcp", "note")
	_, ok := s.Get(22, "tcp", "note")
	if ok {
		t.Error("expected note to be removed")
	}
}

func TestPersistence(t *testing.T) {
	p := tempPath(t)
	s1, _ := New(p)
	_ = s1.Set(8080, "tcp", "env", "staging")

	s2, err := New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	v, ok := s2.Get(8080, "tcp", "env")
	if !ok || v != "staging" {
		t.Errorf("after reload: Get = %q, %v", v, ok)
	}
}

func TestNew_MissingFile(t *testing.T) {
	s, err := New(filepath.Join(t.TempDir(), "notexist.json"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.All()) != 0 {
		t.Error("expected empty store")
	}
}

func TestFprint_Output(t *testing.T) {
	s, _ := New(tempPath(t))
	_ = s.Set(80, "tcp", "team", "ops")
	_ = s.Set(443, "tcp", "tier", "prod")

	var buf bytes.Buffer
	Fprint(&buf, s)
	out := buf.String()
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
	for _, want := range []string{"80", "tcp", "ops", "443", "prod"} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestFprint_Empty(t *testing.T) {
	s, _ := New(tempPath(t))
	var buf bytes.Buffer
	Fprint(&buf, s)
	if !bytes.Contains(buf.Bytes(), []byte("no memos")) {
		t.Errorf("expected 'no memos' message, got: %s", buf.String())
	}
}

func TestRemove_EntireEntry(t *testing.T) {
	p := tempPath(t)
	s, _ := New(p)
	_ = s.Set(9090, "udp", "k", "v")
	_ = s.Remove(9090, "udp", "k")
	if len(s.All()) != 0 {
		t.Error("expected store to be empty after removing last note")
	}
	// verify file was updated
	s2, _ := New(p)
	if len(s2.All()) != 0 {
		t.Error("expected persisted store to be empty")
	}
}

func init() {
	// ensure os package is used
	_ = os.DevNull
}
