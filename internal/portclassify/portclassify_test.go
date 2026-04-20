package portclassify

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 23, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
		{Port: 6379, Proto: "tcp"},
		{Port: 9999, Proto: "tcp"},
	}
}

func TestClassify_KnownPorts(t *testing.T) {
	cl := New(DefaultConfig())
	results := cl.Classify(samplePorts())

	expected := map[int]Class{
		22:   ClassMonitor,
		23:   ClassDangerous,
		80:   ClassSafe,
		6379: ClassDangerous,
		9999: ClassUnknown,
	}

	for _, r := range results {
		want, ok := expected[r.Port.Port]
		if !ok {
			t.Errorf("unexpected port %d in results", r.Port.Port)
			continue
		}
		if r.Class != want {
			t.Errorf("port %d: got class %s, want %s", r.Port.Port, r.Class, want)
		}
	}
}

func TestClassify_ProtoDistinct(t *testing.T) {
	cfg := DefaultConfig()
	cl := New(cfg)

	// port 53 is classified as udp; tcp/53 should be unknown
	ports := []scanner.Port{
		{Port: 53, Proto: "tcp"},
		{Port: 53, Proto: "udp"},
	}
	results := cl.Classify(ports)

	if results[0].Class != ClassUnknown {
		t.Errorf("tcp/53: expected unknown, got %s", results[0].Class)
	}
	if results[1].Class != ClassMonitor {
		t.Errorf("udp/53: expected monitor, got %s", results[1].Class)
	}
}

func TestClassify_EmptyPorts(t *testing.T) {
	cl := New(DefaultConfig())
	results := cl.Classify(nil)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestClassify_CustomRule(t *testing.T) {
	cfg := Config{
		Rules: []Rule{
			{Port: 1234, Proto: "tcp", Class: ClassDangerous, Reason: "custom rule"},
		},
	}
	cl := New(cfg)
	results := cl.Classify([]scanner.Port{{Port: 1234, Proto: "tcp"}})
	if results[0].Class != ClassDangerous {
		t.Errorf("expected dangerous, got %s", results[0].Class)
	}
	if results[0].Reason != "custom rule" {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestFprint_ContainsClass(t *testing.T) {
	cl := New(DefaultConfig())
	results := cl.Classify([]scanner.Port{{Port: 23, Proto: "tcp"}})

	var buf bytes.Buffer
	Fprint(&buf, results)
	out := buf.String()

	if !strings.Contains(out, "dangerous") {
		t.Errorf("expected 'dangerous' in output, got: %s", out)
	}
}
