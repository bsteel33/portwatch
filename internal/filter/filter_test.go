package filter

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Protocol: "tcp", Service: "ssh"},
		{Port: 80, Protocol: "tcp", Service: "http"},
		{Port: 443, Protocol: "tcp", Service: "https"},
		{Port: 53, Protocol: "udp", Service: "dns"},
	}
}

func TestApply_NoRules(t *testing.T) {
	f := New(DefaultConfig())
	ports := samplePorts()
	got := f.Apply(ports)
	if len(got) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestApply_IncludePort(t *testing.T) {
	cfg := Config{Include: []Rule{{Port: 80, Protocol: "tcp"}}}
	f := New(cfg)
	got := f.Apply(samplePorts())
	if len(got) != 1 || got[0].Port != 80 {
		t.Fatalf("expected only port 80, got %v", got)
	}
}

func TestApply_ExcludePort(t *testing.T) {
	cfg := Config{Exclude: []Rule{{Port: 22}}}
	f := New(cfg)
	got := f.Apply(samplePorts())
	for _, p := range got {
		if p.Port == 22 {
			t.Fatal("port 22 should have been excluded")
		}
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(got))
	}
}

func TestApply_ExcludeProto(t *testing.T) {
	cfg := Config{Exclude: []Rule{{Protocol: "udp"}}}
	f := New(cfg)
	got := f.Apply(samplePorts())
	for _, p := range got {
		if p.Protocol == "udp" {
			t.Fatalf("udp port should be excluded: %v", p)
		}
	}
}

func TestApply_IncludeAndExclude(t *testing.T) {
	cfg := Config{
		Include: []Rule{{Protocol: "tcp"}},
		Exclude: []Rule{{Port: 443}},
	}
	f := New(cfg)
	got := f.Apply(samplePorts())
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d: %v", len(got), got)
	}
}
