package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeReport(format Format) (*Report, *bytes.Buffer) {
	ports := []scanner.Port{
		{Port: 80, Proto: "tcp", Service: "http"},
		{Port: 443, Proto: "tcp", Service: "https"},
	}
	snap := &snapshot.Snapshot{Ports: ports, CreatedAt: time.Now()}
	diff := &snapshot.Diff{
		Added:   []scanner.Port{{Port: 443, Proto: "tcp", Service: "https"}},
		Removed: []scanner.Port{},
	}
	buf := &bytes.Buffer{}
	r := New(snap, diff)
	r.Format = format
	r.Writer = buf
	return r, buf
}

func TestRender_Text(t *testing.T) {
	r, buf := makeReport(FormatText)
	if err := r.Render(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "portwatch report") {
		t.Error("expected report header in text output")
	}
	if !strings.Contains(out, "443") {
		t.Error("expected port 443 in text output")
	}
}

func TestRender_JSON(t *testing.T) {
	r, buf := makeReport(FormatJSON)
	if err := r.Render(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"timestamp\"") {
		t.Error("expected timestamp field in JSON output")
	}
	if !strings.Contains(out, "443") {
		t.Error("expected port 443 in JSON output")
	}
}

func TestRender_NoChanges(t *testing.T) {
	snap := &snapshot.Snapshot{Ports: []scanner.Port{}, CreatedAt: time.Now()}
	diff := &snapshot.Diff{Added: []scanner.Port{}, Removed: []scanner.Port{}}
	buf := &bytes.Buffer{}
	r := New(snap, diff)
	r.Writer = buf
	if err := r.Render(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes detected") {
		t.Error("expected no-changes message")
	}
}
