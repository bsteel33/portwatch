package portpoll

import (
	"flag"
	"time"
)

// Config controls Poller behaviour.
type Config struct {
	Timeout  time.Duration
	Host     string
	Proto    string
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout: 2 * time.Second,
		Host:    "127.0.0.1",
		Proto:   "tcp",
	}
}

// RegisterFlags registers portpoll flags on the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.Timeout, "poll-timeout", cfg.Timeout, "dial timeout per port poll")
	fs.StringVar(&cfg.Host, "poll-host", cfg.Host, "host to poll ports on")
	fs.StringVar(&cfg.Proto, "poll-proto", cfg.Proto, "default protocol for polling (tcp/udp)")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Timeout > 0 {
		dst.Timeout = src.Timeout
	}
	if src.Host != "" {
		dst.Host = src.Host
	}
	if src.Proto != "" {
		dst.Proto = src.Proto
	}
}
