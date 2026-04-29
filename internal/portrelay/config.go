package portrelay

import "flag"

// Config holds portrelay configuration.
type Config struct {
	// Enabled controls whether forwarding is active.
	Enabled bool
	// BufferSize is the channel buffer used by async destinations.
	BufferSize int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:    true,
		BufferSize: 32,
	}
}

// RegisterFlags registers portrelay flags on fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.BoolVar(&cfg.Enabled, "relay.enabled", cfg.Enabled, "enable port relay forwarding")
	fs.IntVar(&cfg.BufferSize, "relay.buffer", cfg.BufferSize, "async destination channel buffer size")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.BufferSize > 0 {
		dst.BufferSize = src.BufferSize
	}
	dst.Enabled = src.Enabled
}
