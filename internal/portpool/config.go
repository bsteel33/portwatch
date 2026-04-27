package portpool

import "flag"

// Config holds configuration for a Pool.
type Config struct {
	Name     string
	Capacity int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Name:     "default",
		Capacity: 0,
	}
}

// RegisterFlags registers pool-related CLI flags into the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Name, "pool-name", cfg.Name, "name of the port pool")
	fs.IntVar(&cfg.Capacity, "pool-capacity", cfg.Capacity, "maximum number of ports in the pool (0 = unlimited)")
}

// ApplyFlags copies non-zero flag values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Name != "" && src.Name != "default" {
		dst.Name = src.Name
	}
	if src.Capacity > 0 {
		dst.Capacity = src.Capacity
	}
}
