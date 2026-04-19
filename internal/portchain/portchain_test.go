package portchain_test

import (
	"testing"

	"github.com/user/portwatch/internal/portchain"
	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp", Service: "ssh"},
		{Port: 80, Proto: "tcp", Service: "http"},
		{Port: 8080, Proto: "tcp", Service: "http-alt"},
		{Port: 53, Proto: "udp", Service: "dns"},
	}
}

func TestRun_NoStages(t *testing.T) {
	c := portchain.New(false)
	ports := samplePorts()
	out := c.Run(ports)
	if len(out) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(out))
	}
}

func TestRun_SingleStage(t *testing.T) {
	c := portchain.New(false)
	c.Add("tcp-only", func(ports []scanner.Port) []scanner.Port {
		var out []scanner.Port
		for _, p := range ports {
			if p.Proto == "tcp" {
				out = append(out, p)
			}
		}
		return out
	})
	out := c.Run(samplePorts())
	if len(out) != 3 {
		t.Fatalf("expected 3 tcp ports, got %d", len(out))
	}
}

func TestRun_MultipleStages(t *testing.T) {
	c := portchain.New(false)
	c.Add("tcp-only", func(ports []scanner.Port) []scanner.Port {
		var out []scanner.Port
		for _, p := range ports {
			if p.Proto == "tcp" {
				out = append(out, p)
			}
		}
		return out
	})
	c.Add("high-port", func(ports []scanner.Port) []scanner.Port {
		var out []scanner.Port
		for _, p := range ports {
			if p.Port >= 1024 {
				out = append(out, p)
			}
		}
		return out
	})
	out := c.Run(samplePorts())
	if len(out) != 1 {
		t.Fatalf("expected 1 port, got %d", len(out))
	}
	if out[0].Port != 8080 {
		t.Fatalf("expected port 8080, got %d", out[0].Port)
	}
}

func TestLen(t *testing.T) {
	c := portchain.New(false)
	if c.Len() != 0 {
		t.Fatal("expected 0 stages")
	}
	c.Add("a", func(p []scanner.Port) []scanner.Port { return p })
	c.Add("b", func(p []scanner.Port) []scanner.Port { return p })
	if c.Len() != 2 {
		t.Fatalf("expected 2 stages, got %d", c.Len())
	}
}
