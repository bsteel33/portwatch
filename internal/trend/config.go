package trend

import (
	"flag"
	"time"
)

// Config holds configuration for the Trend tracker.
type Config struct {
	Window time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window: 10 * time.Minute,
	}
}

// RegisterFlags registers trend-related flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.Window, "trend-window", cfg.Window, "sliding window duration for port count trend analysis")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Window != 0 {
		dst.Window = src.Window
	}
}
