package portburst

import (
	"flag"
	"time"
)

// Config holds configuration for burst detection.
type Config struct {
	Threshold int
	Window    time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Threshold: 10,
		Window:    30 * time.Second,
	}
}

// RegisterFlags registers burst-related CLI flags onto the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.IntVar(&cfg.Threshold, "burst-threshold", cfg.Threshold, "number of new ports within window to trigger burst alert")
	fs.DurationVar(&cfg.Window, "burst-window", cfg.Window, "sliding time window for burst detection")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Threshold != 0 {
		dst.Threshold = src.Threshold
	}
	if src.Window != 0 {
		dst.Window = src.Window
	}
}
