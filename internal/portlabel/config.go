package portlabel

import (
	"flag"
	"strings"
)

// Config holds portlabel configuration.
type Config struct {
	// Rules is a list of "port/proto=label" strings.
	Rules []string
}

// DefaultConfig returns a Config with no rules.
func DefaultConfig() Config {
	return Config{}
}

// RegisterFlags registers portlabel flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.Func("label", "port label rule: port/proto=name (repeatable)", func(s string) error {
		cfg.Rules = append(cfg.Rules, s)
		return nil
	})
}

// ParseRules parses config rules into Label slices.
func ParseRules(cfg Config) ([]Label, error) {
	out := make([]Label, 0, len(cfg.Rules))
	for _, r := range cfg.Rules {
		lb, err := parseRule(r)
		if err != nil {
			return nil, err
		}
		out = append(out, lb)
	}
	return out, nil
}

// parseRule parses a single "port/proto=name" rule.
func parseRule(s string) (Label, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 || parts[1] == "" {
		return Label{}, fmt.Errorf("portlabel: invalid rule %q, want port/proto=name", s)
	}
	var port int
	var proto string
	if _, err := fmt.Sscanf(parts[0], "%d/%s", &port, &proto); err != nil {
		return Label{}, fmt.Errorf("portlabel: invalid port/proto %q: %w", parts[0], err)
	}
	return Label{Port: port, Proto: proto, Name: parts[1]}, nil
}
