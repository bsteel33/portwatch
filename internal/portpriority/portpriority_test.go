package portpriority_test

import (
	"bytes"
	"flag"
	"testing"

	"github.com/user/portwatch/internal/portpriority"
	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp", Service: "ssh"},
		{Port: 80, Proto: "tcp", Service: "http"},
		{Port: 443, Proto: "tcp", Service: "https"},
		{Port: 9999, Proto: "udp", Service: ""},
	}
}

func TestAssign_MatchesRule(t *testing.T) {
	rules := []portpriority.Rule{
		{Port: 22, Proto: "tcp", Level: portpriority.Critical},
		{Port: 80, Proto: "tcp", Level: portpriority.High},
	}
	p := portpriority.New(rules, portpriority.Low)
	if got := p.Assign(scanner.Port{Port: 22, Proto: "tcp"}); got != portpriority.Critical {
		t.Fatalf("expected Critical, got %s", got)
	}
	if got := p.Assign(scanner.Port{Port: 80, Proto: "tcp"}); got != portpriority.High {
		t.Fatalf("expected High, got %s", got)
	}
}

func TestAssign_DefaultLevel(t *testing.T) {
	p := portpriority.New(nil, portpriority.Medium)
	if got := p.Assign(scanner.Port{Port: 9999, Proto: "udp"}); got != portpriority.Medium {
		t.Fatalf("expected Medium, got %s", got)
	}
}

func TestAssign_ProtoMismatch_UsesDefault(t *testing.T) {
	rules := []portpriority.Rule{
		{Port: 22, Proto: "tcp", Level: portpriority.Critical},
	}
	p := portpriority.New(rules, portpriority.Low)
	// same port, wrong proto
	if got := p.Assign(scanner.Port{Port: 22, Proto: "udp"}); got != portpriority.Low {
		t.Fatalf("expected Low, got %s", got)
	}
}

func TestAssignAll_ReturnsMap(t *testing.T) {
	rules := []portpriority.Rule{
		{Port: 443, Proto: "tcp", Level: portpriority.High},
	}
	p := portpriority.New(rules, portpriority.Low)
	m := p.AssignAll(samplePorts())
	if m["443/tcp"] != portpriority.High {
		t.Fatalf("expected High for 443/tcp")
	}
	if m["9999/udp"] != portpriority.Low {
		t.Fatalf("expected Low for 9999/udp")
	}
}

func TestParseRule_Valid(t *testing.T) {
	r, err := portpriority.ParseRule("22/tcp=critical")
	if err != nil {
		t.Fatal(err)
	}
	if r.Port != 22 || r.Proto != "tcp" || r.Level != portpriority.Critical {
		t.Fatalf("unexpected rule: %+v", r)
	}
}

func TestParseRule_Invalid(t *testing.T) {
	if _, err := portpriority.ParseRule("badformat"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := portpriority.ParseRule("abc/tcp=high"); err == nil {
		t.Fatal("expected error for non-numeric port")
	}
	if _, err := portpriority.ParseRule("22/tcp=unknown"); err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestBuild_FromConfig(t *testing.T) {
	cfg := portpriority.DefaultConfig()
	cfg.Rules = []string{"22/tcp=critical", "80/tcp=high"}
	p, err := portpriority.Build(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Assign(scanner.Port{Port: 22, Proto: "tcp"}); got != portpriority.Critical {
		t.Fatalf("expected Critical, got %s", got)
	}
}

func TestRegisterFlags(t *testing.T) {
	cfg := portpriority.DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	portpriority.RegisterFlags(fs, &cfg)
	if err := fs.Parse([]string{"--priority.default=high", "--priority.rule=443/tcp=critical"}); err != nil {
		t.Fatal(err)
	}
	if cfg.DefaultLevel != "high" {
		t.Fatalf("expected high, got %s", cfg.DefaultLevel)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0] != "443/tcp=critical" {
		t.Fatalf("unexpected rules: %v", cfg.Rules)
	}
}

func TestFprint_SortedByCriticalFirst(t *testing.T) {
	rules := []portpriority.Rule{
		{Port: 22, Proto: "tcp", Level: portpriority.Critical},
	}
	p := portpriority.New(rules, portpriority.Low)
	var buf bytes.Buffer
	portpriority.Fprint(&buf, samplePorts(), p)
	out := buf.String()
	if len(out) == 0 {
		t.Fatal("expected output")
	}
	// critical should appear before low
	critIdx := bytes.Index(buf.Bytes(), []byte("critical"))
	lowIdx := bytes.Index(buf.Bytes(), []byte("low"))
	if critIdx == -1 || lowIdx == -1 || critIdx > lowIdx {
		t.Fatalf("expected critical before low in output:\n%s", out)
	}
}
