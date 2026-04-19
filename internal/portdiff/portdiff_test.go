package portdiff_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portdiff"
	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap(ports []snapshot.Port) *snapshot.Snapshot {
	return &snapshot.Snapshot{Ports: ports, At: time.Now()}
}

var base = []snapshot.Port{
	{Port: 22, Proto: "tcp", Name: "ssh"},
	{Port: 80, Proto: "tcp", Name: "http"},
}

func TestBuild_Added(t *testing.T) {
	prev := makeSnap(base)
	curr := makeSnap(append(base, snapshot.Port{Port: 443, Proto: "tcp", Name: "https"}))
	s := portdiff.Build(prev, curr)
	if len(s.Added) != 1 || s.Added[0].Port != 443 {
		t.Fatalf("expected 1 added port 443, got %+v", s.Added)
	}
	if len(s.Removed) != 0 {
		t.Fatalf("expected no removed ports")
	}
}

func TestBuild_Removed(t *testing.T) {
	prev := makeSnap(base)
	curr := makeSnap(base[:1])
	s := portdiff.Build(prev, curr)
	if len(s.Removed) != 1 || s.Removed[0].Port != 80 {
		t.Fatalf("expected 1 removed port 80, got %+v", s.Removed)
	}
}

func TestBuild_NoChanges(t *testing.T) {
	prev := makeSnap(base)
	curr := makeSnap(base)
	s := portdiff.Build(prev, curr)
	if s.HasChanges() {
		t.Fatal("expected no changes")
	}
}

func TestFprint_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	portdiff.Fprint(&buf, portdiff.Summary{})
	if buf.String() != "no changes detected\n" {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestFprint_WithChanges(t *testing.T) {
	s := portdiff.Summary{
		Added:   []portdiff.Entry{{Port: 8080, Proto: "tcp", Name: "http-alt", Kind: "added"}},
		Removed: []portdiff.Entry{{Port: 23, Proto: "tcp", Name: "", Kind: "removed"}},
	}
	var buf bytes.Buffer
	portdiff.Fprint(&buf, s)
	out := buf.String()
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	for _, want := range []string{"+", "8080", "-", "23", "(unknown)"} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}
