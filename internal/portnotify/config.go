package portnotify

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Rule defines a single watch condition.
type Rule struct {
	Port  int
	Proto string
	Label string
}

// Config holds Notifier configuration.
type Config struct {
	Rules []Rule
}

// DefaultConfig returns a Config with no rules.
func DefaultConfig() Config {
	return Config{}
}

// RegisterFlags registers CLI flags onto fs and returns a pointer to a raw
// rules string slice that ApplyFlags can consume.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) *[]string {
	var raw []string
	fs.Func("portnotify-rule",
		"watch rule: port[/proto][=label] (repeatable)",
		func(s string) error {
			raw = append(raw, s)
			return nil
		})
	return &raw
}

// ApplyFlags parses raw rule strings into cfg.Rules.
func ApplyFlags(cfg *Config, raw []string) error {
	for _, s := range raw {
		r, err := parseRule(s)
		if err != nil {
			return err
		}
		cfg.Rules = append(cfg.Rules, r)
	}
	return nil
}

// parseRule parses "port[/proto][=label]".
func parseRule(s string) (Rule, error) {
	label := ""
	if idx := strings.Index(s, "="); idx >= 0 {
		label = s[idx+1:]
		s = s[:idx]
	}
	proto := ""
	if idx := strings.Index(s, "/"); idx >= 0 {
		proto = s[idx+1:]
		s = s[:idx]
	}
	port, err := strconv.Atoi(s)
	if err != nil || port < 1 || port > 65535 {
		return Rule{}, fmt.Errorf("portnotify: invalid rule %q", s)
	}
	if label == "" {
		label = fmt.Sprintf("port-%d", port)
	}
	return Rule{Port: port, Proto: proto, Label: label}, nil
}
