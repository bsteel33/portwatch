package portquota

import "flag"

// Config holds quota limits.
type Config struct {
	TotalLimit  int
	ProtoLimits map[string]int
}

// DefaultConfig returns a Config with no limits enforced.
func DefaultConfig() Config {
	return Config{
		ProtoLimits: make(map[string]int),
	}
}

// RegisterFlags registers CLI flags for quota configuration.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.IntVar(&cfg.TotalLimit, "quota-total", cfg.TotalLimit, "max total open ports (0 = unlimited)")
}

// ApplyFlags copies non-zero override values into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.TotalLimit > 0 {
		dst.TotalLimit = src.TotalLimit
	}
	for proto, limit := range src.ProtoLimits {
		if dst.ProtoLimits == nil {
			dst.ProtoLimits = make(map[string]int)
		}
		dst.ProtoLimits[proto] = limit
	}
}
