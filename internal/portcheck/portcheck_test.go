package portcheck

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
		{Port: 53, Proto: "udp"},
	}
}

func TestEvaluate_OpenPass(t *testing.T) {
	ch := New([]Condition{{Port: 80, Proto: "tcp", MustBeOpen: true}})
	results := ch.Evaluate(samplePorts())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Passed {
		t.Errorf("expected pass, got reason: %s", results[0].Reason)
	}
}

func TestEvaluate_OpenFail(t *testing.T) {
	ch := New([]Condition{{Port: 8080, Proto: "tcp", MustBeOpen: true}})
	results := ch.Evaluate(samplePorts())
	if results[0].Passed {
		t.Error("expected fail for closed port")
	}
	if !strings.Contains(results[0].Reason, "expected open") {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestEvaluate_ClosedPass(t *testing.T) {
	ch := New([]Condition{{Port: 22, Proto: "tcp", MustBeOpen: false}})
	results := ch.Evaluate(samplePorts())
	if !results[0].Passed {
		t.Errorf("expected pass, got reason: %s", results[0].Reason)
	}
}

func TestEvaluate_ClosedFail(t *testing.T) {
	ch := New([]Condition{{Port: 80, Proto: "tcp", MustBeOpen: false}})
	results := ch.Evaluate(samplePorts())
	if results[0].Passed {
		t.Error("expected fail for open port that should be closed")
	}
	if !strings.Contains(results[0].Reason, "expected closed") {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestAnyFailed(t *testing.T) {
	passing := []Result{{Passed: true}, {Passed: true}}
	if AnyFailed(passing) {
		t.Error("expected no failures")
	}
	mixed := []Result{{Passed: true}, {Passed: false}}
	if !AnyFailed(mixed) {
		t.Error("expected failure detected")
	}
}

func TestParseRule_Valid(t *testing.T) {
	cases := []struct {
		input      string
		port       int
		proto      string
		mustBeOpen bool
	}{
		{"80/tcp:open", 80, "tcp", true},
		{"53/udp:closed", 53, "udp", false},
	}
	for _, tc := range cases {
		c, err := parseRule(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if c.Port != tc.port || c.Proto != tc.proto || c.MustBeOpen != tc.mustBeOpen {
			t.Errorf("parseRule(%q) = %+v, want port=%d proto=%s open=%v", tc.input, c, tc.port, tc.proto, tc.mustBeOpen)
		}
	}
}

func TestParseRule_Invalid(t *testing.T) {
	invalid := []string{"80tcp:open", "80/tcp", "80/tcp:maybe", "99999/tcp:open", "80/icmp:open"}
	for _, s := range invalid {
		if _, err := parseRule(s); err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}

func TestFprint_Output(t *testing.T) {
	results := []Result{
		{Condition: Condition{Port: 80, Proto: "tcp", MustBeOpen: true}, Passed: true},
		{Condition: Condition{Port: 22, Proto: "tcp", MustBeOpen: false}, Passed: false, Reason: "port 22/tcp expected closed but is open"},
	}
	var buf bytes.Buffer
	Fprint(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Error("expected OK in output")
	}
	if !strings.Contains(out, "[FAIL]") {
		t.Error("expected FAIL in output")
	}
}
