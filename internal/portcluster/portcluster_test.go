package portcluster

import (
	"testing"

	"github.com/example/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
		{Port: 53, Proto: "udp"},
		{Port: 8080, Proto: "tcp"},
	}
}

func TestApply_StrategyProto(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Strategy = StrategyProto
	c := New(cfg)

	clusters := c.Apply(samplePorts())

	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters (tcp, udp), got %d", len(clusters))
	}
	if clusters[0].Name != "tcp" {
		t.Errorf("expected first cluster 'tcp', got %q", clusters[0].Name)
	}
	if clusters[1].Name != "udp" {
		t.Errorf("expected second cluster 'udp', got %q", clusters[1].Name)
	}
	if len(clusters[0].Ports) != 4 {
		t.Errorf("expected 4 tcp ports, got %d", len(clusters[0].Ports))
	}
}

func TestApply_StrategyRange(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Strategy = StrategyRange
	cfg.RangeSize = 1024
	c := New(cfg)

	clusters := c.Apply(samplePorts())

	// ports 22,53,80,443 → 0-1023; port 8080 → 8192-9215
	if len(clusters) != 2 {
		t.Fatalf("expected 2 range clusters, got %d", len(clusters))
	}
	if clusters[0].Name != "0-1023" {
		t.Errorf("expected '0-1023', got %q", clusters[0].Name)
	}
	if len(clusters[0].Ports) != 4 {
		t.Errorf("expected 4 ports in 0-1023, got %d", len(clusters[0].Ports))
	}
}

func TestApply_EmptyPorts(t *testing.T) {
	c := New(DefaultConfig())
	clusters := c.Apply(nil)
	if len(clusters) != 0 {
		t.Errorf("expected 0 clusters for empty input, got %d", len(clusters))
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Strategy != StrategyRange {
		t.Errorf("expected default strategy %q, got %q", StrategyRange, cfg.Strategy)
	}
	if cfg.RangeSize != 1024 {
		t.Errorf("expected default range size 1024, got %d", cfg.RangeSize)
	}
}

func TestApplyFlags_Override(t *testing.T) {
	dst := DefaultConfig()
	src := Config{Strategy: StrategyProto, RangeSize: 512}
	ApplyFlags(&dst, src)
	if dst.Strategy != StrategyProto {
		t.Errorf("expected strategy %q, got %q", StrategyProto, dst.Strategy)
	}
	if dst.RangeSize != 512 {
		t.Errorf("expected range size 512, got %d", dst.RangeSize)
	}
}
