package portreport_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portreport"
	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
		{Port: 53, Proto: "udp"},
	}
}

func TestNew_Sorted(t *testing.T) {
	r := portreport.New(samplePorts(), nil)
	if r.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", r.Len())
	}
	ports := []int{}
	for _, e := range r.Entries() {
		ports = append(ports, e.Port)
	}
	for i := 1; i < len(ports); i++ {
		if ports[i] < ports[i-1] {
			t.Errorf("entries not sorted: %v", ports)
		}
	}
}

func TestNew_ServiceResolved(t *testing.T) {
	svc := func(port int, proto string) string {
		if port == 22 && proto == "tcp" {
			return "ssh"
		}
		return ""
	}
	r := portreport.New(samplePorts(), svc)
	for _, e := range r.Entries() {
		if e.Port == 22 && e.Service != "ssh" {
			t.Errorf("expected ssh, got %q", e.Service)
		}
	}
}

func TestFprint_Empty(t *testing.T) {
	r := portreport.New(nil, nil)
	var buf bytes.Buffer
	portreport.Fprint(&buf, r)
	if !strings.Contains(buf.String(), "no open ports") {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}

func TestFprint_ContainsHeaders(t *testing.T) {
	r := portreport.New(samplePorts(), nil)
	var buf bytes.Buffer
	portreport.Fprint(&buf, r)
	for _, col := range []string{"PORT", "PROTO", "SERVICE", "STATE"} {
		if !strings.Contains(buf.String(), col) {
			t.Errorf("missing column header %q", col)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := portreport.DefaultConfig()
	if !cfg.ShowUnknown {
		t.Error("expected ShowUnknown to be true by default")
	}
	if cfg.SortByScore {
		t.Error("expected SortByScore to be false by default")
	}
}
