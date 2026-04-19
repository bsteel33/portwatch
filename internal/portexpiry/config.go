package portexpiry

import (
	"flag"
	"time"
)

// Config holds portexpiry settings.
type Config struct {
	Path    string
	MaxAge  time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path:   "portexpiry.json",
		MaxAge: 24 * time.Hour,
	}
}

// RegisterFlags registers CLI flags onto fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "expiry-path", cfg.Path, "path to port expiry state file")
	fs.DurationVar(&cfg.MaxAge, "expiry-max-age", cfg.MaxAge, "maximum duration a port may remain open before alerting")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Path != "" {
		dst.Path = src.Path
	}
	if src.MaxAge != 0 {
		dst.MaxAge = src.MaxAge
	}
}
