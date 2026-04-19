package portsampler

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func fakePorts(n int) []scanner.Port {
	ports := make([]scanner.Port, n)
	for i := range ports {
		ports[i] = scanner.Port{Port: 8000 + i, Proto: "tcp"}
	}
	return ports
}

func TestSampler_RecordsSamples(t *testing.T) {
	cfg := Config{Interval: 20 * time.Millisecond, MaxSamples: 10}
	s := New(cfg)
	calls := 0
	s.Start(func() ([]scanner.Port, error) {
		calls++
		return fakePorts(3), nil
	})
	time.Sleep(70 * time.Millisecond)
	s.Stop()
	samples := s.Samples()
	if len(samples) == 0 {
		t.Fatal("expected at least one sample")
	}
	if samples[0].Count != 3 {
		t.Errorf("expected count 3, got %d", samples[0].Count)
	}
}

func TestSampler_MaxSamples(t *testing.T) {
	cfg := Config{Interval: 10 * time.Millisecond, MaxSamples: 3}
	s := New(cfg)
	s.Start(func() ([]scanner.Port, error) { return fakePorts(1), nil })
	time.Sleep(60 * time.Millisecond)
	s.Stop()
	if got := len(s.Samples()); got > 3 {
		t.Errorf("expected at most 3 samples, got %d", got)
	}
}

func TestSampler_Last_Empty(t *testing.T) {
	s := New(DefaultConfig())
	_, ok := s.Last()
	if ok {
		t.Error("expected no sample on empty sampler")
	}
}

func TestSampler_Last_ReturnsMostRecent(t *testing.T) {
	s := New(Config{Interval: 10 * time.Millisecond, MaxSamples: 10})
	s.Start(func() ([]scanner.Port, error) { return fakePorts(5), nil })
	time.Sleep(40 * time.Millisecond)
	s.Stop()
	last, ok := s.Last()
	if !ok {
		t.Fatal("expected a sample")
	}
	if last.Count != 5 {
		t.Errorf("expected count 5, got %d", last.Count)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval <= 0 {
		t.Error("expected positive interval")
	}
	if cfg.MaxSamples <= 0 {
		t.Error("expected positive MaxSamples")
	}
}
