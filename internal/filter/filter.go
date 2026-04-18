package filter

import "github.com/user/portwatch/internal/scanner"

// Rule defines a single port filter rule.
type Rule struct {
	Port     int
	Protocol string // "tcp" or "udp"
	Action   string // "include" or "exclude"
}

// Filter applies include/exclude rules to a list of ports.
type Filter struct {
	cfg Config
}

// New returns a new Filter with the given config.
func New(cfg Config) *Filter {
	return &Filter{cfg: cfg}
}

// Apply filters ports according to the configured rules.
// If include rules exist, only matching ports are kept.
// Exclude rules always remove matching ports.
func (f *Filter) Apply(ports []scanner.Port) []scanner.Port {
	var result []scanner.Port
	for _, p := range ports {
		if f.excluded(p) {
			continue
		}
		if len(f.cfg.Include) > 0 && !f.included(p) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func (f *Filter) included(p scanner.Port) bool {
	for _, r := range f.cfg.Include {
		if matches(r, p) {
			return true
		}
	}
	return false
}

func (f *Filter) excluded(p scanner.Port) bool {
	for _, r := range f.cfg.Exclude {
		if matches(r, p) {
			return true
		}
	}
	return false
}

func matches(r Rule, p scanner.Port) bool {
	portMatch := r.Port == 0 || r.Port == p.Port
	protoMatch := r.Protocol == "" || r.Protocol == p.Protocol
	return portMatch && protoMatch
}
