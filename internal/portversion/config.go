package portversion

import "flag"

// Config holds configuration for the version tracker.
type Config struct {
	Path    string
	Enabled bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path:    "portversion.json",
		Enabled: true,
	}
}

// RegisterFlags registers CLI flags into the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "version-path", cfg.Path, "path to port version state file")
	fs.BoolVar(&cfg.Enabled, "version-track", cfg.Enabled, "enable port version/banner tracking")
}

// ApplyFlags copies non-zero flag values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Path != "" {
		dst.Path = src.Path
	}
	dst.Enabled = src.Enabled
}
