package tags_test

import (
	"os"
	"path/filepath"
	"testing"

	"portwatch/internal/tags"
)

func TestResolve_Found(t *testing.T) {
	tr := &tags.Tagger{}
	tr = newTagger(t, []tags.Tag{
		{Port: 22, Proto: "tcp", Label: "SSH"},
		{Port: 80, Proto: "tcp", Label: "HTTP"},
	})
	if got := tr.Resolve(22, "tcp"); got != "SSH" {
		t.Fatalf("expected SSH, got %q", got)
	}
	if got := tr.Resolve(80, "tcp"); got != "HTTP" {
		t.Fatalf("expected HTTP, got %q", got)
	}
}

func TestResolve_NotFound(t *testing.T) {
	tr := newTagger(t, []tags.Tag{{Port: 22, Proto: "tcp", Label: "SSH"}})
	if got := tr.Resolve(443, "tcp"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestResolve_ProtoMismatch(t *testing.T) {
	tr := newTagger(t, []tags.Tag{{Port: 53, Proto: "udp", Label: "DNS"}})
	if got := tr.Resolve(53, "tcp"); got != "" {
		t.Fatalf("expected empty for proto mismatch, got %q", got)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")

	tr, err := tags.New("")
	if err != nil {
		t.Fatal(err)
	}
	tr.Add(8080, "tcp", "Dev HTTP")
	tr.Add(9090, "tcp", "Metrics")

	if err := tr.Save(path); err != nil {
		t.Fatal(err)
	}

	loaded, err := tags.New(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := loaded.Resolve(8080, "tcp"); got != "Dev HTTP" {
		t.Fatalf("expected 'Dev HTTP', got %q", got)
	}
}

func TestNew_MissingFile(t *testing.T) {
	tr, err := tags.New("/nonexistent/tags.json")
	if err != nil {
		t.Fatal("expected no error for missing file, got:", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tagger")
	}
}

func TestKey(t *testing.T) {
	if got := tags.Key(443, "TCP"); got != "443/tcp" {
		t.Fatalf("expected 443/tcp, got %q", got)
	}
}

// newTagger builds a Tagger from a slice of Tag via save+load roundtrip.
func newTagger(t *testing.T, ts []tags.Tag) *tags.Tagger {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")
	tr, _ := tags.New("")
	for _, tag := range ts {
		tr.Add(tag.Port, tag.Proto, tag.Label)
	}
	if err := tr.Save(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := tags.New(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(path)
	return loaded
}
