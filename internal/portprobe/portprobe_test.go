package portprobe

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRecord_And_Get(t *testing.T) {
	tr := New()
	tr.Record(80, "tcp", true, 5*time.Millisecond)
	r, ok := tr.Get(80, "tcp")
	if !ok {
		t.Fatal("expected result to exist")
	}
	if !r.Reachable {
		t.Error("expected reachable=true")
	}
	if r.Latency != 5*time.Millisecond {
		t.Errorf("unexpected latency: %v", r.Latency)
	}
}

func TestGet_Missing(t *testing.T) {
	tr := New()
	_, ok := tr.Get(9999, "tcp")
	if ok {
		t.Error("expected no result for unknown port")
	}
}

func TestUnreachable_ReturnsOnlyFailed(t *testing.T) {
	tr := New()
	tr.Record(80, "tcp", true, 2*time.Millisecond)
	tr.Record(443, "tcp", false, 0)
	tr.Record(22, "tcp", false, 0)
	unreachable := tr.Unreachable()
	if len(unreachable) != 2 {
		t.Errorf("expected 2 unreachable, got %d", len(unreachable))
	}
	for _, r := range unreachable {
		if r.Reachable {
			t.Errorf("port %d should not be reachable", r.Port)
		}
	}
}

func TestAll_ReturnsAllResults(t *testing.T) {
	tr := New()
	tr.Record(80, "tcp", true, 1*time.Millisecond)
	tr.Record(53, "udp", true, 2*time.Millisecond)
	all := tr.All()
	if len(all) != 2 {
		t.Errorf("expected 2 results, got %d", len(all))
	}
}

func TestReset_ClearsResults(t *testing.T) {
	tr := New()
	tr.Record(80, "tcp", true, 1*time.Millisecond)
	tr.Reset()
	if len(tr.All()) != 0 {
		t.Error("expected empty results after reset")
	}
}

func TestFprint_ContainsHeaders(t *testing.T) {
	tr := New()
	tr.Record(80, "tcp", true, 3*time.Millisecond)
	var buf bytes.Buffer
	Fprint(&buf, tr.All())
	out := buf.String()
	for _, hdr := range []string{"PORT", "PROTO", "REACHABLE", "LATENCY"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFprint_Empty(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, nil)
	if !strings.Contains(buf.String(), "no probe results") {
		t.Error("expected empty message")
	}
}
