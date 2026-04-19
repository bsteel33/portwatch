// Package portsampler periodically samples open ports and records counts over time.
package portsampler

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Sample holds a single observation.
type Sample struct {
	At    time.Time
	Count int
	Ports []scanner.Port
}

// Sampler collects port samples at a fixed interval.
type Sampler struct {
	mu      sync.Mutex
	cfg     Config
	samples []Sample
	stop    chan struct{}
	wg      sync.WaitGroup
}

// New returns a new Sampler with the given config.
func New(cfg Config) *Sampler {
	return &Sampler{
		cfg:  cfg,
		stop: make(chan struct{}),
	}
}

// Start begins sampling in a background goroutine.
func (s *Sampler) Start(scan func() ([]scanner.Port, error)) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(s.cfg.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-s.stop:
				return
			case t := <-ticker.C:
				ports, err := scan()
				if err != nil {
					continue
				}
				s.record(t, ports)
			}
		}
	}()
}

// Stop halts background sampling.
func (s *Sampler) Stop() {
	close(s.stop)
	s.wg.Wait()
}

func (s *Sampler) record(at time.Time, ports []scanner.Port) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.samples = append(s.samples, Sample{At: at, Count: len(ports), Ports: ports})
	if s.cfg.MaxSamples > 0 && len(s.samples) > s.cfg.MaxSamples {
		s.samples = s.samples[len(s.samples)-s.cfg.MaxSamples:]
	}
}

// Samples returns a copy of all recorded samples.
func (s *Sampler) Samples() []Sample {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Sample, len(s.samples))
	copy(out, s.samples)
	return out
}

// Last returns the most recent sample, or false if none.
func (s *Sampler) Last() (Sample, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) == 0 {
		return Sample{}, false
	}
	return s.samples[len(s.samples)-1], true
}
