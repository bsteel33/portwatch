package portsampler

import (
	"flag"
	"time"
)

// Config controls sampler behaviour.
type Config struct {
	Interval   time.Duration
	MaxSamples int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:   30 * time.Second,
		MaxSamples: 120,
	}
}

// RegisterFlags registers sampler flags onto fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.Interval, "sampler-interval", cfg.Interval, "how often to sample open ports")
	fs.IntVar(&cfg.MaxSamples, "sampler-max", cfg.MaxSamples, "maximum number of samples to retain (0 = unlimited)")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Interval != 0 {
		dst.Interval = src.Interval
	}
	if src.MaxSamples != 0 {
		dst.MaxSamples = src.MaxSamples
	}
}
