package portacl

import (
	"flag"
	"strings"
)

// Config holds ACL configuration.
type Config struct {
	Rules         []string
	DefaultAction Action
}

// DefaultConfig returns a permissive default configuration.
func DefaultConfig() Config {
	return Config{
		DefaultAction: Allow,
	}
}

// RegisterFlags registers ACL flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.Func("acl-rule", "ACL rule (e.g. deny:23/tcp). May be repeated.", func(s string) error {
		cfg.Rules = append(cfg.Rules, s)
		return nil
	})
	fs.Func("acl-default", "Default ACL action when no rule matches (allow|deny)", func(s string) error {
		a := Action(strings.ToLower(s))
		if a != Allow && a != Deny {
			return nil
		}
		cfg.DefaultAction = a
		return nil
	})
}

// Build constructs an ACL from the config, returning an error if any rule
// fails to parse.
func Build(cfg Config) (*ACL, error) {
	rules := make([]Rule, 0, len(cfg.Rules))
	for _, raw := range cfg.Rules {
		r, err := ParseRule(raw)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return New(rules, cfg.DefaultAction), nil
}
