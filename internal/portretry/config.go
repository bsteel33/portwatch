package portretry

import (
	"flag"
	"time"
)

// Config holds configuration for the Retryer.
type Config struct {
	MaxAttempts int
	Delay       time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
	}
}

// RegisterFlags registers portretry flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.IntVar(&cfg.MaxAttempts, "retry-attempts", cfg.MaxAttempts,
		"number of scan retry attempts on failure")
	fs.DurationVar(&cfg.Delay, "retry-delay", cfg.Delay,
		"delay between retry attempts")
}

// ApplyFlags overwrites cfg fields from parsed flag values when non-zero.
func ApplyFlags(cfg *Config, attempts int, delay time.Duration) {
	if attempts > 0 {
		cfg.MaxAttempts = attempts
	}
	if delay > 0 {
		cfg.Delay = delay
	}
}
