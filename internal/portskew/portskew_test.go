package portskew

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePorts(n int) []scanner.Port {
	ports := make([]scanner.Port, n)
	for i := range ports {
		ports[i] = scanner.Port{Port: i + 1, Proto: "tcp"}
	}
	return ports
}

func TestAnalyze_BelowMinSamples(t *testing.T) {
	cfg := DefaultConfig()
	d := New(cfg)

	r := d.Analyze(makePorts(10), cfg)
	if r.Skewed {
		t.Error("expected no skew before min samples reached")
	}
	if r.Current != 10 {
		t.Errorf("expected current=10, got %d", r.Current)
	}
}

func TestAnalyze_NoSkew(t *testing.T) {
	cfg := DefaultConfig()
	d := New(cfg)

	// Feed stable samples.
	for i := 0; i < 6; i++ {
		r := d.Analyze(makePorts(10), cfg)
		if i >= cfg.MinSamples-1 && r.Skewed {
			t.Errorf("scan %d: unexpected skew for stable port count", i)
		}
	}
}

func TestAnalyze_Skewed(t *testing.T) {
	cfg := DefaultConfig()
	d := New(cfg)

	// Establish a stable baseline of 10 ports.
	for i := 0; i < 6; i++ {
		d.Analyze(makePorts(10), cfg)
	}

	// Sudden spike to 100 ports should trigger skew.
	r := d.Analyze(makePorts(100), cfg)
	if !r.Skewed {
		t.Errorf("expected skew, z=%.2f", r.ZScore)
	}
}

func TestReset_ClearsSamples(t *testing.T) {
	cfg := DefaultConfig()
	d := New(cfg)

	for i := 0; i < 6; i++ {
		d.Analyze(makePorts(10), cfg)
	}
	d.Reset()

	// After reset, below MinSamples again.
	r := d.Analyze(makePorts(100), cfg)
	if r.Skewed {
		t.Error("expected no skew after reset")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Threshold <= 0 {
		t.Errorf("expected positive threshold, got %f", cfg.Threshold)
	}
	if cfg.MinSamples <= 0 {
		t.Errorf("expected positive MinSamples, got %d", cfg.MinSamples)
	}
}

func TestApplyFlags_Override(t *testing.T) {
	dst := DefaultConfig()
	src := Config{Threshold: 3.5, MinSamples: 10}
	ApplyFlags(&dst, src)
	if dst.Threshold != 3.5 {
		t.Errorf("expected threshold 3.5, got %f", dst.Threshold)
	}
	if dst.MinSamples != 10 {
		t.Errorf("expected MinSamples 10, got %d", dst.MinSamples)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	dst := DefaultConfig()
	orig := dst
	ApplyFlags(&dst, Config{})
	if dst.Threshold != orig.Threshold {
		t.Error("threshold should not change when src is zero")
	}
}
