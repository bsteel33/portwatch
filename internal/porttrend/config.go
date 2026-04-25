package porttrend

import (
	"flag"
	"time"
)

// Config controls Tracker behaviour.
type Config struct {
	// Window is how far back samples are retained.
	Window time.Duration
	// Threshold is the minimum absolute delta required to declare a
	// directional trend (avoids noise for tiny fluctuations).
	Threshold int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window:    10 * time.Minute,
		Threshold: 2,
	}
}

// RegisterFlags registers porttrend flags on fs and returns a pointer to a
// Config that will be populated after fs.Parse.
func RegisterFlags(fs *flag.FlagSet) *Config {
	cfg := DefaultConfig()
	fs.DurationVar(&cfg.Window, "porttrend.window", cfg.Window,
		"sliding window for port-count trend analysis")
	fs.IntVar(&cfg.Threshold, "porttrend.threshold", cfg.Threshold,
		"minimum delta to declare a directional trend")
	return &cfg
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Window != 0 {
		dst.Window = src.Window
	}
	if src.Threshold != 0 {
		dst.Threshold = src.Threshold
	}
}
