package schedule

import (
	"flag"
	"time"
)

// Config holds schedule settings.
type Config struct {
	Interval   time.Duration
	DelayFirst bool
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:   30 * time.Second,
		DelayFirst: false,
	}
}

// RegisterFlags registers schedule-related flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.Interval, "schedule.interval", cfg.Interval, "how often to run a scan")
	fs.BoolVar(&cfg.DelayFirst, "schedule.delay-first", cfg.DelayFirst, "skip immediate first run and wait for first tick")
}

// ApplyFlags overwrites cfg fields from src when non-zero.
func ApplyFlags(cfg *Config, src Config) {
	if src.Interval != 0 {
		cfg.Interval = src.Interval
	}
	if src.DelayFirst {
		cfg.DelayFirst = true
	}
}
