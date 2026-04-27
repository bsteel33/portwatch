package portjournal

import "flag"

// Config holds configuration for the port journal.
type Config struct {
	Path    string
	MaxLast int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path:    "portjournal.json",
		MaxLast: 50,
	}
}

// RegisterFlags registers journal-related CLI flags onto fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "journal-path", cfg.Path, "path to port journal file")
	fs.IntVar(&cfg.MaxLast, "journal-last", cfg.MaxLast, "number of recent entries to show with --journal")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Path != "" {
		dst.Path = src.Path
	}
	if src.MaxLast > 0 {
		dst.MaxLast = src.MaxLast
	}
}
