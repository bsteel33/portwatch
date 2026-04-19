package portquota

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
		{Port: 53, Proto: "udp"},
	}
}

func TestCheck_NoViolations(t *testing.T) {
	cfg := DefaultConfig()
	q := New(cfg)
	violations := q.Check(samplePorts())
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_TotalExceeded(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TotalLimit = 2
	q := New(cfg)
	v := q.Check(samplePorts())
	if len(v) != 1 {
		t.Fatalf("expected 1 (v))
	}
	if v[0].Proto != "any" || v[0].Actual != 4 {
		t.Errorf("unexpected violation: %+v", v[0])
	}
}

func TestCheck_ProtoExceeded(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ProtoLimits = map[string]int{"tcp": 2}
	q := New(cfg)
	v := q.Check(samplePorts())
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Proto != "tcp" || v[0].Limit != 2 || v[0].Actual != 3 {
		t.Errorf("unexpected violation: %+v", v[0])
	}
}

func TestCheck_MultipleViolations(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TotalLimit = 1
	cfg.ProtoLimits = map[string]int{"tcp": 1, "udp": 0}
	q := New(cfg)
	v := q.Check(samplePorts())
	if len(v) < 2 {
		t.Fatalf("expected multiple violations, got %d", len(v))
	}
}

func TestTotals_ReflectsLastCheck(t *testing.T) {
	q := New(DefaultConfig())
	q.Check(samplePorts())
	totals := q.Totals()
	if totals["tcp"] != 3 {
		t.Errorf("expected tcp=3, got %d", totals["tcp"])
	}
	if totals["udp"] != 1 {
		t.Errorf("expected udp=1, got %d", totals["udp"])
	}
}
