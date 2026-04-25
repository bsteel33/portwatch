package portevents

import "flag"

// Config controls optional behaviour of the Bus.
type Config struct {
	// BufferSize is reserved for future async dispatch; currently unused.
	BufferSize int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BufferSize: 0,
	}
}

// RegisterFlags binds Config fields to the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.IntVar(&cfg.BufferSize, "events-buffer", cfg.BufferSize,
		"event bus buffer size (0 = synchronous)")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.BufferSize > 0 {
		dst.BufferSize = src.BufferSize
	}
}
