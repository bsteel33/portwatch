package portclone_test

import (
	"testing"

	"github.com/user/portwatch/internal/portclone"
	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Proto: "tcp", Service: "http"},
		{Port: 443, Proto: "tcp", Service: "https"},
		{Port: 53, Proto: "udp", Service: "dns"},
	}
}

func TestClone_ReturnsCopy(t *testing.T) {
	c := portclone.New()
	orig := samplePorts()
	got := c.Clone(orig)
	if len(got) != len(orig) {
		t.Fatalf("expected len %d, got %d", len(orig), len(got))
	}
	got[0].Port = 9999
	if orig[0].Port == 9999 {
		t.Fatal("Clone should not share memory with original")
	}
}

func TestClone_Nil(t *testing.T) {
	c := portclone.New()
	if got := c.Clone(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestCloneMap_ReturnsCopy(t *testing.T) {
	c := portclone.New()
	orig := map[string][]scanner.Port{
		"host1": samplePorts(),
	}
	got := c.CloneMap(orig)
	if len(got["host1"]) != len(orig["host1"]) {
		t.Fatalf("unexpected length")
	}
	got["host1"][0].Port = 1234
	if orig["host1"][0].Port == 1234 {
		t.Fatal("CloneMap should not share slice memory")
	}
}

func TestMerge_DeduplicatesAndOverrides(t *testing.T) {
	c := portclone.New()
	a := []scanner.Port{
		{Port: 80, Proto: "tcp", Service: "http"},
		{Port: 22, Proto: "tcp", Service: "ssh"},
	}
	b := []scanner.Port{
		{Port: 80, Proto: "tcp", Service: "http-override"},
		{Port: 8080, Proto: "tcp", Service: "http-alt"},
	}
	merged := c.Merge(a, b)
	if len(merged) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(merged))
	}
	for _, p := range merged {
		if p.Port == 80 && p.Service != "http-override" {
			t.Errorf("expected b to override a for port 80, got service %q", p.Service)
		}
	}
}

func TestMerge_EmptySlices(t *testing.T) {
	c := portclone.New()
	if got := c.Merge(nil, nil); len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}
