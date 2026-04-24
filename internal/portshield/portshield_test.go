package portshield

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp", Service: "ssh"},
		{Port: 80, Proto: "tcp", Service: "http"},
		{Port: 443, Proto: "tcp", Service: "https"},
		{Port: 8080, Proto: "tcp", Service: "http-alt"},
	}
}

func TestEvaluate_DefaultAllow(t *testing.T) {
	s := New(Allow)
	p := scanner.Port{Port: 9999, Proto: "tcp"}
	if got := s.Evaluate(p); got != Allow {
		t.Fatalf("expected Allow, got %s", got)
	}
}

func TestEvaluate_DefaultBlock(t *testing.T) {
	s := New(Block)
	p := scanner.Port{Port: 9999, Proto: "tcp"}
	if got := s.Evaluate(p); got != Block {
		t.Fatalf("expected Block, got %s", got)
	}
}

func TestEvaluate_ExplicitRule(t *testing.T) {
	s := New(Allow)
	s.Add(22, "tcp", Block)
	p := scanner.Port{Port: 22, Proto: "tcp"}
	if got := s.Evaluate(p); got != Block {
		t.Fatalf("expected Block for port 22/tcp, got %s", got)
	}
}

func TestFilter_RemovesBlocked(t *testing.T) {
	s := New(Allow)
	s.Add(22, "tcp", Block)
	s.Add(8080, "tcp", Block)

	result := s.Filter(samplePorts())
	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(result))
	}
	for _, p := range result {
		if p.Port == 22 || p.Port == 8080 {
			t.Errorf("blocked port %d/tcp should not be in result", p.Port)
		}
	}
}

func TestBuild_AllowAndBlock(t *testing.T) {
	cfg := Config{
		AllowPorts:    "80/tcp,443/tcp",
		BlockPorts:    "22/tcp",
		DefaultAction: "allow",
	}
	s, err := Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Evaluate(scanner.Port{Port: 22, Proto: "tcp"}); got != Block {
		t.Errorf("expected Block for 22/tcp, got %s", got)
	}
	if got := s.Evaluate(scanner.Port{Port: 80, Proto: "tcp"}); got != Allow {
		t.Errorf("expected Allow for 80/tcp, got %s", got)
	}
}

func TestBuild_InvalidRule(t *testing.T) {
	cfg := Config{AllowPorts: "notaport", DefaultAction: "allow"}
	_, err := Build(cfg)
	if err == nil {
		t.Fatal("expected error for invalid rule, got nil")
	}
}

func TestAction_String(t *testing.T) {
	if Allow.String() != "allow" {
		t.Errorf("expected 'allow', got %s", Allow.String())
	}
	if Block.String() != "block" {
		t.Errorf("expected 'block', got %s", Block.String())
	}
}
