package filter

import (
	"flag"
	"strconv"
	"strings"
)

// Config holds include and exclude filter rules.
type Config struct {
	Include []Rule
	Exclude []Rule
}

// DefaultConfig returns an empty filter config (no filtering).
func DefaultConfig() Config {
	return Config{}
}

// RegisterFlags registers filter-related CLI flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.Func("include-port", "include only this port (repeatable, format: port or port/proto)", func(s string) error {
		r, err := parseRule(s)
		if err != nil {
			return err
		}
		cfg.Include = append(cfg.Include, r)
		return nil
	})
	fs.Func("exclude-port", "exclude this port (repeatable, format: port or port/proto)", func(s string) error {
		r, err := parseRule(s)
		if err != nil {
			return err
		}
		cfg.Exclude = append(cfg.Exclude, r)
		return nil
	})
}

// parseRule parses a rule string like "80" or "80/tcp".
func parseRule(s string) (Rule, error) {
	parts := strings.SplitN(s, "/", 2)
	port, err := strconv.Atoi(parts[0])
	if err != nil {
		return Rule{}, err
	}
	proto := ""
	if len(parts) == 2 {
		proto = strings.ToLower(parts[1])
	}
	return Rule{Port: port, Protocol: proto}, nil
}
