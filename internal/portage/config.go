package portage

import (
	"flag"
	"time"
)

// Config holds configuration for the port age tracker.
type Config struct {
	// Path is the file used to persist first-seen timestamps.
	Path string
	// MaxAge is the duration after which a closed port entry is evicted.
	MaxAge time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path:   "portage.json",
		MaxAge: 7 * 24 * time.Hour,
	}
}

// RegisterFlags registers command-line flags for the port age tracker.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "portage.path", cfg.Path, "path to port age persistence file")
	fs.DurationVar(&cfg.MaxAge, "portage.max-age", cfg.MaxAge, "how long to retain closed port age records")
}

// ApplyFlags copies non-zero flag values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Path != "" {
		dst.Path = src.Path
	}
	if src.MaxAge > 0 {
		dst.MaxAge = src.MaxAge
	}
}
