package portfence_test

import (
	"net"
	"testing"

	"portwatch/internal/portfence"
)

func mustParseCIDR(s string) *net.IPNet {
	_, cidr, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return cidr
}

func TestEvaluate_AllowMatchingCIDR(t *testing.T) {
	f := portfence.New([]portfence.Rule{
		{CIDR: mustParseCIDR("192.168.1.0/24"), Action: portfence.ActionAllow},
	}, portfence.ActionBlock)

	ip := net.ParseIP("192.168.1.50")
	if got := f.Evaluate(ip); got != portfence.ActionAllow {
		t.Errorf("expected allow, got %s", got)
	}
}

func TestEvaluate_BlockMatchingCIDR(t *testing.T) {
	f := portfence.New([]portfence.Rule{
		{CIDR: mustParseCIDR("10.0.0.0/8"), Action: portfence.ActionBlock},
	}, portfence.ActionAllow)

	ip := net.ParseIP("10.5.6.7")
	if got := f.Evaluate(ip); got != portfence.ActionBlock {
		t.Errorf("expected block, got %s", got)
	}
}

func TestEvaluate_DefaultActionOnNoMatch(t *testing.T) {
	f := portfence.New([]portfence.Rule{
		{CIDR: mustParseCIDR("172.16.0.0/12"), Action: portfence.ActionBlock},
	}, portfence.ActionAllow)

	ip := net.ParseIP("8.8.8.8")
	if got := f.Evaluate(ip); got != portfence.ActionAllow {
		t.Errorf("expected default allow, got %s", got)
	}
}

func TestAllowed_ReturnsTrue(t *testing.T) {
	f := portfence.New([]portfence.Rule{
		{CIDR: mustParseCIDR("192.168.0.0/16"), Action: portfence.ActionAllow},
	}, portfence.ActionBlock)

	if !f.Allowed(net.ParseIP("192.168.2.1")) {
		t.Error("expected Allowed to return true")
	}
}

func TestParseRule_Valid(t *testing.T) {
	r, err := portfence.ParseRule("10.0.0.0/8:block")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Action != portfence.ActionBlock {
		t.Errorf("expected block, got %s", r.Action)
	}
	if !r.CIDR.Contains(net.ParseIP("10.1.2.3")) {
		t.Error("CIDR should contain 10.1.2.3")
	}
}

func TestParseRule_InvalidFormat(t *testing.T) {
	if _, err := portfence.ParseRule("nocolon"); err == nil {
		t.Error("expected error for missing colon")
	}
}

func TestParseRule_InvalidCIDR(t *testing.T) {
	if _, err := portfence.ParseRule("notacidr:allow"); err == nil {
		t.Error("expected error for bad CIDR")
	}
}

func TestParseRule_InvalidAction(t *testing.T) {
	if _, err := portfence.ParseRule("192.168.0.0/24:permit"); err == nil {
		t.Error("expected error for unknown action")
	}
}
