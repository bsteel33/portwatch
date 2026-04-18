package watch

import (
	"flag"
	"time"
)

// Config holds watcher-specific configuration.
type Config struct {
	Interval     time.Duration
	SnapshotPath string
}

// DefaultConfig returns sensible defaults for the watcher.
func DefaultConfig() Config {
	return Config{
		Interval:     30 * time.Second,
		SnapshotPath: "/var/lib/portwatch/snapshot.json",
	}
}

// ApplyFlags overrides Config fields from parsed CLI flags.
func ApplyFlags(cfg *Config, fs *flag.FlagSet) {
	if f := fs.Lookup("interval"); f != nil && f.Value.String() != f.DefValue {
		if d, err := time.ParseDuration(f.Value.String()); err == nil {
			cfg.Interval = d
		}
	}
	if f := fs.Lookup("snapshot"); f != nil && f.Value.String() != f.DefValue {
		cfg.SnapshotPath = f.Value.String()
	}
}
