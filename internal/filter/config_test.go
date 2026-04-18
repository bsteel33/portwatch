package filter

import (
	"flag"
	"testing"
)

func TestParseRule_PortOnly(t *testing.T) {
	r, err := parseRule("8080")
	if err != nil {
		t.Fatal(err)
	}
	if r.Port != 8080 || r.Protocol != "" {
		t.Fatalf("unexpected rule: %+v", r)
	}
}

func TestParseRule_PortAndProto(t *testing.T) {
	r, err := parseRule("443/tcp")
	if err != nil {
		t.Fatal(err)
	}
	if r.Port != 443 || r.Protocol != "tcp" {
		t.Fatalf("unexpected rule: %+v", r)
	}
}

func TestParseRule_Invalid(t *testing.T) {
	_, err := parseRule("notaport")
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestRegisterFlags_Include(t *testing.T) {
	cfg := DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs, &cfg)
	if err := fs.Parse([]string{"-include-port", "80/tcp", "-include-port", "443"}); err != nil {
		t.Fatal(err)
	}
	if len(cfg.Include) != 2 {
		t.Fatalf("expected 2 include rules, got %d", len(cfg.Include))
	}
}

func TestRegisterFlags_Exclude(t *testing.T) {
	cfg := DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs, &cfg)
	if err := fs.Parse([]string{"-exclude-port", "22/tcp"}); err != nil {
		t.Fatal(err)
	}
	if len(cfg.Exclude) != 1 || cfg.Exclude[0].Port != 22 {
		t.Fatalf("unexpected exclude rules: %+v", cfg.Exclude)
	}
}
