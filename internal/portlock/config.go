package portlock

import "flag"

// Config holds portlock configuration.
type Config struct {
	Path string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path: "portlock.json",
	}
}

// RegisterFlags registers portlock flags onto fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "lock-file", cfg.Path, "path to port lock file")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Path != "" {
		dst.Path = src.Path
	}
}
