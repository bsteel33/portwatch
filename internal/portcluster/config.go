package portcluster

import "flag"

// Strategy controls how ports are grouped into clusters.
type Strategy string

const (
	// StrategyProto groups ports by network protocol (tcp/udp).
	StrategyProto Strategy = "proto"
	// StrategyRange groups ports by numeric range bands.
	StrategyRange Strategy = "range"
)

// Config holds configuration for the Clusterer.
type Config struct {
	// Strategy determines the clustering algorithm.
	Strategy Strategy
	// RangeSize is the band width used when Strategy is StrategyRange.
	RangeSize int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Strategy:  StrategyRange,
		RangeSize: 1024,
	}
}

// RegisterFlags registers CLI flags into the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar((*string)(&cfg.Strategy), "cluster-strategy", string(cfg.Strategy),
		"port clustering strategy: proto or range")
	fs.IntVar(&cfg.RangeSize, "cluster-range-size", cfg.RangeSize,
		"band width for range-based clustering")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Strategy != "" {
		dst.Strategy = src.Strategy
	}
	if src.RangeSize > 0 {
		dst.RangeSize = src.RangeSize
	}
}
