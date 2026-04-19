package portrank

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 23, Proto: "tcp"},   // telnet — critical weight
		{Port: 8080, Proto: "tcp"}, // no special weight, registered range
		{Port: 22, Proto: "tcp"},   // ssh — high weight
		{Port: 9999, Proto: "udp"}, // unknown + udp bonus
	}
}

func TestRank_OrderedByScore(t *testing.T) {
	r := New(DefaultConfig())
	entries := r.Rank(samplePorts())
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Score > entries[i-1].Score {
			t.Errorf("entries not sorted: index %d score %d > index %d score %d",
				i, entries[i].Score, i-1, entries[i-1].Score)
		}
	}
}

func TestRank_TelnetIsCritical(t *testing.T) {
	r := New(DefaultConfig())
	entries := r.Rank([]scanner.Port{{Port: 23, Proto: "tcp"}})
	if entries[0].Level != Critical {
		t.Errorf("expected critical, got %s", entries[0].Level)
	}
}

func TestRank_UnknownPortIsLow(t *testing.T) {
	cfg := DefaultConfig()
	r := New(cfg)
	entries := r.Rank([]scanner.Port{{Port: 49200, Proto: "tcp"}})
	if entries[0].Level != Low {
		t.Errorf("expected low, got %s", entries[0].Level)
	}
}

func TestRank_EmptyPorts(t *testing.T) {
	r := New(DefaultConfig())
	entries := r.Rank(nil)
	if len(entries) != 0 {
		t.Errorf("expected empty slice")
	}
}

func TestFprint_ContainsHeaders(t *testing.T) {
	r := New(DefaultConfig())
	entries := r.Rank(samplePorts())
	var buf bytes.Buffer
	Fprint(&buf, entries)
	out := buf.String()
	for _, h := range []string{"PORT", "PROTO", "RISK", "SCORE"} {
		if !strings.Contains(out, h) {
			t.Errorf("output missing header %q", h)
		}
	}
}

func TestFprint_Empty(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, nil)
	if !strings.Contains(buf.String(), "no ports") {
		t.Errorf("expected empty message")
	}
}
