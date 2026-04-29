package portcooldown

import (
	"flag"
	"time"
)

// Config holds configuration for the cooldown tracker.
type Config struct {
	Window time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window: 5 * time.Minute,
	}
}

// RegisterFlags registers cooldown-related CLI flags into the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.Window, "cooldown-window", cfg.Window,
		"duration to suppress repeated alerts for the same port after a state change")
}

// ApplyFlags copies non-zero flag values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Window > 0 {
		dst.Window = src.Window
	}
}
