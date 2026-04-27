package portbatch

import (
	"flag"
	"time"
)

// Config controls Batcher behaviour.
type Config struct {
	// BatchSize is the maximum number of ports per batch. 0 means no size limit.
	BatchSize int
	// FlushInterval is the maximum time between automatic flushes. 0 disables
	// timer-based flushing.
	FlushInterval time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BatchSize:     50,
		FlushInterval: 5 * time.Second,
	}
}

// RegisterFlags registers portbatch flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.IntVar(&cfg.BatchSize, "batch-size", cfg.BatchSize,
		"maximum number of ports per processing batch (0 = unlimited)")
	fs.DurationVar(&cfg.FlushInterval, "batch-flush-interval", cfg.FlushInterval,
		"maximum interval between automatic batch flushes")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.BatchSize != 0 {
		dst.BatchSize = src.BatchSize
	}
	if src.FlushInterval != 0 {
		dst.FlushInterval = src.FlushInterval
	}
}
