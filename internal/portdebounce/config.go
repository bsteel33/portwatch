package portdebounce

import (
	"flag"
	"time"
)

// Config controls debounce behaviour.
type Config struct {
	// Window is the minimum duration a port must be continuously observed
	// before it is considered stable and acted upon.
	Window time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Window: 10 * time.Second,
	}
}

// RegisterFlags registers debounce flags on fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.Window, "debounce-window", cfg.Window,
		"duration a port must be stable before triggering an alert")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Window > 0 {
		dst.Window = src.Window
	}
}
