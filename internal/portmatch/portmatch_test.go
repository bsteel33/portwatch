package portmatch_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmatch"
	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
		{Port: 53, Proto: "udp"},
		{Port: 8080, Proto: "tcp"},
	}
}

func TestParseRule_Valid(t *testing.T) {
	cases := []struct{ in, port, proto string }{
		{"80", "80", "*"},
		{"80/tcp", "80", "tcp"},
		{"*/udp", "*", "udp"},
		{"*", "*", "*"},
	}
	for _, c := range cases {
		r, err := portmatch.ParseRule(c.in)
		if err != nil {
			t.Fatalf("ParseRule(%q) error: %v", c.in, err)
		}
		if r.Port != c.port || r.Proto != c.proto {
			t.Errorf("ParseRule(%q) = {%s %s}, want {%s %s}", c.in, r.Port, r.Proto, c.port, c.proto)
		}
	}
}

func TestParseRule_Invalid(t *testing.T) {
	invalid := []string{"abc", "80/sctp", "abc/tcp"}
	for _, s := range invalid {
		if _, err := portmatch.ParseRule(s); err == nil {
			t.Errorf("ParseRule(%q) expected error, got nil", s)
		}
	}
}

func TestMatch_ByPort(t *testing.T) {
	m, _ := portmatch.New([]string{"80/tcp"})
	if !m.Match(scanner.Port{Port: 80, Proto: "tcp"}) {
		t.Error("expected match for 80/tcp")
	}
	if m.Match(scanner.Port{Port: 443, Proto: "tcp"}) {
		t.Error("unexpected match for 443/tcp")
	}
}

func TestMatch_Wildcard(t *testing.T) {
	m, _ := portmatch.New([]string{"*/udp"})
	if !m.Match(scanner.Port{Port: 53, Proto: "udp"}) {
		t.Error("expected match for 53/udp")
	}
	if m.Match(scanner.Port{Port: 53, Proto: "tcp"}) {
		t.Error("unexpected match for 53/tcp with */udp rule")
	}
}

func TestFilter_ReturnsMatching(t *testing.T) {
	m, _ := portmatch.New([]string{"80/tcp", "53/udp"})
	got := m.Filter(samplePorts())
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestFilter_NoRules_ReturnsAll(t *testing.T) {
	m, _ := portmatch.New(nil)
	got := m.Filter(samplePorts())
	if len(got) != len(samplePorts()) {
		t.Errorf("expected all ports returned, got %d", len(got))
	}
}
