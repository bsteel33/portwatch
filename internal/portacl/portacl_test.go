package portacl

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 23, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
		{Port: 53, Proto: "udp"},
	}
}

func TestEvaluate_Allow(t *testing.T) {
	acl := New([]Rule{{Port: 23, Proto: "tcp", Action: Deny}}, Allow)
	if got := acl.Evaluate(scanner.Port{Port: 80, Proto: "tcp"}); got != Allow {
		t.Fatalf("expected allow, got %s", got)
	}
}

func TestEvaluate_Deny(t *testing.T) {
	acl := New([]Rule{{Port: 23, Proto: "tcp", Action: Deny}}, Allow)
	if got := acl.Evaluate(scanner.Port{Port: 23, Proto: "tcp"}); got != Deny {
		t.Fatalf("expected deny, got %s", got)
	}
}

func TestFilter_RemovesDenied(t *testing.T) {
	acl := New([]Rule{{Port: 23, Proto: "tcp", Action: Deny}}, Allow)
	filtered := acl.Filter(samplePorts())
	for _, p := range filtered {
		if p.Port == 23 {
			t.Fatal("port 23 should have been filtered out")
		}
	}
	if len(filtered) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(filtered))
	}
}

func TestParseRule_Valid(t *testing.T) {
	r, err := ParseRule("deny:23/tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Port != 23 || r.Proto != "tcp" || r.Action != Deny {
		t.Fatalf("unexpected rule: %+v", r)
	}
}

func TestParseRule_Invalid(t *testing.T) {
	if _, err := ParseRule("badformat"); err == nil {
		t.Fatal("expected error for bad format")
	}
	if _, err := ParseRule("unknown:80"); err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestBuild_FromConfig(t *testing.T) {
	cfg := Config{
		Rules:         []string{"deny:23/tcp", "deny:21/tcp"},
		DefaultAction: Allow,
	}
	acl, err := Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acl.Evaluate(scanner.Port{Port: 23, Proto: "tcp"}) != Deny {
		t.Fatal("expected port 23 to be denied")
	}
	if acl.Evaluate(scanner.Port{Port: 80, Proto: "tcp"}) != Allow {
		t.Fatal("expected port 80 to be allowed")
	}
}
