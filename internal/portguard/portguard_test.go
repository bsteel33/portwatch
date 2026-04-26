package portguard

import (
	"testing"

	"github.com/example/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
		{Port: 8080, Proto: "tcp"},
	}
}

func TestCheck_AllAllowed(t *testing.T) {
	cfg := DefaultConfig()
	g := New(cfg)
	viols := g.Check(samplePorts())
	if len(viols) != 0 {
		t.Fatalf("expected no violations, got %d", len(viols))
	}
}

func TestCheck_DenylistViolation(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Denylist = []string{"22/tcp"}
	g := New(cfg)
	viols := g.Check(samplePorts())
	if len(viols) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(viols))
	}
	if viols[0].Port.Port != 22 {
		t.Errorf("expected port 22, got %d", viols[0].Port.Port)
	}
	if viols[0].Action != ActionDeny {
		t.Errorf("expected deny action, got %s", viols[0].Action)
	}
}

func TestCheck_DefaultDeny_AllowlistOnly(t *testing.T) {
	cfg := Config{
		Allowlist: []string{"80/tcp", "443/tcp"},
		Default:   ActionDeny,
	}
	g := New(cfg)
	viols := g.Check(samplePorts())
	// 22/tcp and 8080/tcp should be denied
	if len(viols) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(viols))
	}
	for _, v := range viols {
		if v.Action != ActionDeny {
			t.Errorf("expected deny, got %s for port %d", v.Action, v.Port.Port)
		}
	}
}

func TestEvaluate_ExplicitAllow(t *testing.T) {
	cfg := Config{
		Allowlist: []string{"443/tcp"},
		Default:   ActionDeny,
	}
	g := New(cfg)
	action, v := g.Evaluate(scanner.Port{Port: 443, Proto: "tcp"})
	if action != ActionAllow {
		t.Errorf("expected allow, got %s", action)
	}
	if v != nil {
		t.Errorf("expected no violation for explicitly allowed port")
	}
}

func TestEvaluate_DenylistBeatsAllowlist(t *testing.T) {
	cfg := Config{
		Allowlist: []string{"80/tcp"},
		Denylist:  []string{"80/tcp"},
		Default:   ActionAllow,
	}
	g := New(cfg)
	action, v := g.Evaluate(scanner.Port{Port: 80, Proto: "tcp"})
	if action != ActionDeny {
		t.Errorf("expected deny, got %s", action)
	}
	if v == nil {
		t.Error("expected violation, got nil")
	}
}

func TestCheck_NoViolations_EmptyPorts(t *testing.T) {
	cfg := Config{Default: ActionDeny}
	g := New(cfg)
	viols := g.Check(nil)
	if len(viols) != 0 {
		t.Errorf("expected no violations for empty ports, got %d", len(viols))
	}
}
