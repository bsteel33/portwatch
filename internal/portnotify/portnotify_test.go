package portnotify

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp", Service: "ssh"},
		{Port: 80, Proto: "tcp", Service: "http"},
		{Port: 3306, Proto: "tcp", Service: "mysql"},
		{Port: 53, Proto: "udp", Service: "dns"},
	}
}

func TestCheck_MatchesPort(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rules = []Rule{{Port: 22, Proto: "tcp", Label: "ssh-watch"}}
	n := New(cfg)
	events := n.Check(samplePorts())
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Rule.Label != "ssh-watch" {
		t.Errorf("unexpected label: %s", events[0].Rule.Label)
	}
}

func TestCheck_ProtoMismatch(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rules = []Rule{{Port: 53, Proto: "tcp", Label: "dns-tcp"}}
	n := New(cfg)
	events := n.Check(samplePorts())
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestCheck_NoRules(t *testing.T) {
	n := New(DefaultConfig())
	if got := n.Check(samplePorts()); len(got) != 0 {
		t.Errorf("expected no events, got %d", len(got))
	}
}

func TestCheck_MultipleRules(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rules = []Rule{
		{Port: 22, Proto: "tcp", Label: "ssh"},
		{Port: 3306, Proto: "tcp", Label: "mysql"},
	}
	n := New(cfg)
	if got := n.Check(samplePorts()); len(got) != 2 {
		t.Errorf("expected 2 events, got %d", len(got))
	}
}

func TestParseRule_Valid(t *testing.T) {
	r, err := parseRule("443/tcp=https-watch")
	if err != nil {
		t.Fatal(err)
	}
	if r.Port != 443 || r.Proto != "tcp" || r.Label != "https-watch" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRule_Invalid(t *testing.T) {
	if _, err := parseRule("notaport"); err == nil {
		t.Error("expected error for invalid rule")
	}
}

func TestFprint_Events(t *testing.T) {
	events := []Event{
		{Port: scanner.Port{Port: 22, Proto: "tcp", Service: "ssh"}, Rule: Rule{Label: "ssh-watch"}},
	}
	var buf bytes.Buffer
	Fprint(&buf, events)
	if !strings.Contains(buf.String(), "ssh-watch") {
		t.Errorf("output missing label: %s", buf.String())
	}
}

func TestFprint_NoEvents(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, nil)
	if !strings.Contains(buf.String(), "no matching") {
		t.Errorf("expected no-events message, got: %s", buf.String())
	}
}
