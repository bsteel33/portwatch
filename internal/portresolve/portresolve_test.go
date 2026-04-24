package portresolve_test

import (
	"testing"

	"github.com/user/portwatch/internal/portresolve"
	"github.com/user/portwatch/internal/scanner"
)

func TestResolve_BuiltinHTTP(t *testing.T) {
	r := portresolve.New(portresolve.DefaultConfig())
	p := scanner.Port{Port: 80, Proto: "tcp"}
	res := r.Resolve(p)
	if res.Source != "builtin" {
		t.Fatalf("expected source=builtin, got %q", res.Source)
	}
	if res.Name == "" {
		t.Fatal("expected non-empty name for port 80/tcp")
	}
}

func TestResolve_Override(t *testing.T) {
	r := portresolve.New(portresolve.DefaultConfig())
	r.Override(9090, "tcp", "my-service")
	p := scanner.Port{Port: 9090, Proto: "tcp"}
	res := r.Resolve(p)
	if res.Source != "override" {
		t.Fatalf("expected source=override, got %q", res.Source)
	}
	if res.Name != "my-service" {
		t.Fatalf("expected name=my-service, got %q", res.Name)
	}
}

func TestResolve_Unknown(t *testing.T) {
	r := portresolve.New(portresolve.DefaultConfig())
	p := scanner.Port{Port: 59999, Proto: "tcp"}
	res := r.Resolve(p)
	if res.Source != "unknown" {
		t.Fatalf("expected source=unknown, got %q", res.Source)
	}
	if res.Name != "port-59999" {
		t.Fatalf("expected name=port-59999, got %q", res.Name)
	}
}

func TestResolveAll_ReturnsMatchingLength(t *testing.T) {
	r := portresolve.New(portresolve.DefaultConfig())
	ports := []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
		{Port: 12345, Proto: "udp"},
	}
	results := r.ResolveAll(ports)
	if len(results) != len(ports) {
		t.Fatalf("expected %d results, got %d", len(ports), len(results))
	}
	for i, res := range results {
		if res.Port != ports[i].Port || res.Proto != ports[i].Proto {
			t.Errorf("result[%d] port/proto mismatch", i)
		}
	}
}

func TestResolve_ProtoDistinct(t *testing.T) {
	r := portresolve.New(portresolve.DefaultConfig())
	r.Override(53, "tcp", "dns-tcp")
	tcpRes := r.Resolve(scanner.Port{Port: 53, Proto: "tcp"})
	udpRes := r.Resolve(scanner.Port{Port: 53, Proto: "udp"})
	if tcpRes.Name != "dns-tcp" {
		t.Errorf("expected dns-tcp, got %q", tcpRes.Name)
	}
	if udpRes.Source == "override" {
		t.Error("udp should not use tcp override")
	}
}
