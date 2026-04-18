package config

import (
	"flag"
	"time"
)

// Flags holds values parsed from command-line flags.
type Flags struct {
	ConfigPath   string
	SnapshotPath string
	Interval     time.Duration
	Verbose      bool
	Once         bool
}

// ParseFlags parses os.Args and returns a Flags struct.
func ParseFlags() *Flags {
	f := &Flags{}
	flag.StringVar(&f.ConfigPath, "config", "", "path to JSON config file")
	flag.StringVar(&f.SnapshotPath, "snapshot", "", "path to snapshot file (overrides config)")
	flag.DurationVar(&f.Interval, "interval", 0, "scan interval (overrides config)")
	flag.BoolVar(&f.Verbose, "verbose", false, "enable verbose output")
	flag.BoolVar(&f.Once, "once", false, "run a single scan and exit")
	flag.Parse()
	return f
}

// Apply merges CLI flags into a Config, with flags taking precedence.
func Apply(cfg *Config, f *Flags) *Config {
	if f.SnapshotPath != "" {
		cfg.SnapshotPath = f.SnapshotPath
	}
	if f.Interval != 0 {
		cfg.ScanInterval = f.Interval
	}
	if f.Verbose {
		cfg.Verbose = true
	}
	return cfg
}
