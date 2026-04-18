package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Number: 22, Proto: "tcp", Service: "ssh"},
		{Number: 80, Proto: "tcp", Service: "http"},
		{Number: 443, Proto: "tcp", Service: "https"},
	}
}

func TestCompute_Deterministic(t *testing.T) {
	ports := samplePorts()
	a := fingerprint.Compute(ports)
	b := fingerprint.Compute(ports)
	if !fingerprint.Equal(a, b) {
		t.Fatalf("expected equal fingerprints, got %s vs %s", a, b)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	forward := samplePorts()
	reversed := []scanner.Port{forward[2], forward[1], forward[0]}

	a := fingerprint.Compute(forward)
	b := fingerprint.Compute(reversed)
	if !fingerprint.Equal(a, b) {
		t.Fatalf("order should not affect fingerprint: %s vs %s", a, b)
	}
}

func TestCompute_DifferentPorts(t *testing.T) {
	a := fingerprint.Compute(samplePorts())
	other := []scanner.Port{
		{Number: 8080, Proto: "tcp", Service: "http-alt"},
	}
	b := fingerprint.Compute(other)
	if fingerprint.Equal(a, b) {
		t.Fatal("expected different fingerprints for different port sets")
	}
}

func TestCompute_EmptyPorts(t *testing.T) {
	a := fingerprint.Compute(nil)
	b := fingerprint.Compute([]scanner.Port{})
	if !fingerprint.Equal(a, b) {
		t.Fatalf("nil and empty should produce same fingerprint: %s vs %s", a, b)
	}
}

func TestCompute_ProtoDistinct(t *testing.T) {
	tcp := []scanner.Port{{Number: 53, Proto: "tcp"}}
	udp := []scanner.Port{{Number: 53, Proto: "udp"}}
	if fingerprint.Equal(fingerprint.Compute(tcp), fingerprint.Compute(udp)) {
		t.Fatal("tcp and udp on same port should differ")
	}
}
