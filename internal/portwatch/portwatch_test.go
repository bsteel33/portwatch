package portwatch

import (
	"bytes"
	"testing"
	"time"

	"github.com/example/portwatch/internal/scanner"
)

func tempSnap(t *testing.T) string {
	t.Helper()
	return t.TempDir() + "/snap.json"
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.SnapshotPath == "" {
		t.Fatal("expected non-empty snapshot path")
	}
	if cfg.Timeout <= 0 {
		t.Fatal("expected positive timeout")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SnapshotPath = tempSnap(t)
	if err := Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_BadProto(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Proto = "icmp"
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for bad proto")
	}
}

func TestApplyFlags_Override(t *testing.T) {
	dst := DefaultConfig()
	src := Config{SnapshotPath: "/tmp/x.json", Timeout: 10 * time.Second, Proto: "udp"}
	ApplyFlags(&dst, src)
	if dst.SnapshotPath != "/tmp/x.json" {
		t.Errorf("got %q", dst.SnapshotPath)
	}
	if dst.Proto != "udp" {
		t.Errorf("got %q", dst.Proto)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	dst := DefaultConfig()
	origPath := dst.SnapshotPath
	ApplyFlags(&dst, Config{})
	if dst.SnapshotPath != origPath {
		t.Errorf("snapshot path should not change")
	}
}

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
		{Port: 53, Proto: "udp"},
	}
}

func TestEvaluate_MatchesRule(t *testing.T) {
	w := New(DefaultConfig())
	w.AddRule(Rule{Name: "ssh", Port: 22, Proto: "tcp"})
	events := w.Evaluate(samplePorts())
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Port.Port != 22 {
		t.Errorf("expected port 22, got %d", events[0].Port.Port)
	}
}

func TestEvaluate_ProtoMismatch(t *testing.T) {
	w := New(DefaultConfig())
	w.AddRule(Rule{Name: "ssh-udp", Port: 22, Proto: "udp"})
	if events := w.Evaluate(samplePorts()); len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestEvaluate_NoRules(t *testing.T) {
	w := New(DefaultConfig())
	if events := w.Evaluate(samplePorts()); len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestFprint_NoEvents(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, nil)
	if buf.Len() == 0 {
		t.Fatal("expected output for empty events")
	}
}

func TestFprint_WithEvents(t *testing.T) {
	w := New(DefaultConfig())
	w.AddRule(Rule{Name: "http", Port: 80, Proto: "tcp"})
	events := w.Evaluate(samplePorts())
	var buf bytes.Buffer
	Fprint(&buf, events)
	if buf.Len() == 0 {
		t.Fatal("expected output")
	}
}
