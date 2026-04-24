package portdigest_test

import (
	"testing"

	"github.com/user/portwatch/internal/portdigest"
	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp", Service: "ssh"},
		{Port: 80, Proto: "tcp", Service: "http"},
	}
}

func TestCompute_Deterministic(t *testing.T) {
	ports := samplePorts()
	d1 := portdigest.Compute(ports)
	d2 := portdigest.Compute(ports)
	if d1 != d2 {
		t.Fatalf("expected identical digests, got %q and %q", d1, d2)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := []scanner.Port{{Port: 22, Proto: "tcp"}, {Port: 80, Proto: "tcp"}}
	b := []scanner.Port{{Port: 80, Proto: "tcp"}, {Port: 22, Proto: "tcp"}}
	if portdigest.Compute(a) != portdigest.Compute(b) {
		t.Fatal("digest should be order-independent")
	}
}

func TestCompute_DifferentPorts(t *testing.T) {
	a := samplePorts()
	b := []scanner.Port{{Port: 443, Proto: "tcp"}}
	if portdigest.Compute(a) == portdigest.Compute(b) {
		t.Fatal("expected different digests for different port sets")
	}
}

func TestUpdate_FirstCall_NoChange(t *testing.T) {
	tr := portdigest.New()
	changed := tr.Update(samplePorts())
	if changed {
		t.Fatal("first update should not report a change")
	}
	if tr.Current().Digest == "" {
		t.Fatal("current digest should be set after first update")
	}
}

func TestUpdate_SamePorts_NoChange(t *testing.T) {
	tr := portdigest.New()
	tr.Update(samplePorts())
	changed := tr.Update(samplePorts())
	if changed {
		t.Fatal("same ports should not report a change")
	}
}

func TestUpdate_DifferentPorts_Changed(t *testing.T) {
	tr := portdigest.New()
	tr.Update(samplePorts())
	newPorts := []scanner.Port{{Port: 443, Proto: "tcp"}}
	changed := tr.Update(newPorts)
	if !changed {
		t.Fatal("different ports should report a change")
	}
}

func TestPrevious_AfterTwoUpdates(t *testing.T) {
	tr := portdigest.New()
	tr.Update(samplePorts())
	first := tr.Current().Digest
	tr.Update([]scanner.Port{{Port: 443, Proto: "tcp"}})
	if tr.Previous().Digest != first {
		t.Fatalf("expected previous digest %q, got %q", first, tr.Previous().Digest)
	}
}

func TestReset_ClearsState(t *testing.T) {
	tr := portdigest.New()
	tr.Update(samplePorts())
	tr.Reset()
	if tr.Current().Digest != "" {
		t.Fatal("expected empty digest after reset")
	}
}
