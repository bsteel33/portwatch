package portscore

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Proto: "tcp", State: "open"},
		{Port: 23, Proto: "tcp", State: "open"},
		{Port: 8080, Proto: "tcp", State: "open"},
		{Port: 53, Proto: "udp", State: "open"},
	}
}

func TestScore_WellKnownAndTCP(t *testing.T) {
	cfg := DefaultConfig()
	s := New(cfg)
	scores := s.Score([]scanner.Port{{Port: 80, Proto: "tcp", State: "open"}})
	if len(scores) != 1 {
		t.Fatalf("expected 1 score, got %d", len(scores))
	}
	got := scores[0].Total
	want := cfg.WellKnownWeight + cfg.TCPWeight
	if got != want {
		t.Errorf("port 80/tcp: got total %d, want %d", got, want)
	}
}

func TestScore_HighRiskPort(t *testing.T) {
	cfg := DefaultConfig()
	s := New(cfg)
	scores := s.Score([]scanner.Port{{Port: 23, Proto: "tcp", State: "open"}})
	sc := scores[0]
	if sc.Factors["high_risk"] != 50 {
		t.Errorf("expected high_risk factor 50, got %d", sc.Factors["high_risk"])
	}
	if sc.Total > 100 {
		t.Errorf("total %d exceeds 100", sc.Total)
	}
}

func TestScore_UDPNonWellKnown(t *testing.T) {
	cfg := DefaultConfig()
	s := New(cfg)
	scores := s.Score([]scanner.Port{{Port: 9999, Proto: "udp", State: "open"}})
	sc := scores[0]
	if sc.Total != 0 {
		t.Errorf("expected 0 for high non-well-known udp port, got %d", sc.Total)
	}
}

func TestScore_CappedAt100(t *testing.T) {
	cfg := DefaultConfig()
	cfg.WellKnownWeight = 60
	cfg.TCPWeight = 60
	s := New(cfg)
	scores := s.Score([]scanner.Port{{Port: 23, Proto: "tcp", State: "open"}})
	if scores[0].Total != 100 {
		t.Errorf("expected total capped at 100, got %d", scores[0].Total)
	}
}

func TestScore_AllPorts(t *testing.T) {
	cfg := DefaultConfig()
	s := New(cfg)
	scores := s.Score(samplePorts())
	if len(scores) != 4 {
		t.Fatalf("expected 4 scores, got %d", len(scores))
	}
	for _, sc := range scores {
		if sc.Total < 0 || sc.Total > 100 {
			t.Errorf("score out of range [0,100]: %d for %v", sc.Total, sc.Port)
		}
	}
}
