// Package portname resolves well-known port numbers to human-readable service names.
package portname

import (
	"fmt"
	"strings"
)

// Entry holds a resolved port name and optional description.
type Entry struct {
	Name        string
	Description string
}

// Resolver maps port/proto pairs to service names.
type Resolver struct {
	table map[string]Entry
}

// New returns a Resolver seeded with common well-known ports.
func New() *Resolver {
	r := &Resolver{table: make(map[string]Entry)}
	for _, e := range builtinEntries {
		r.table[key(e.port, e.proto)] = Entry{Name: e.name, Description: e.desc}
	}
	return r
}

// Resolve returns the Entry for the given port and protocol.
// proto should be "tcp" or "udp". ok is false when unknown.
func (r *Resolver) Resolve(port int, proto string) (Entry, bool) {
	e, ok := r.table[key(port, strings.ToLower(proto))]
	return e, ok
}

// Name returns just the service name, or a formatted fallback.
func (r *Resolver) Name(port int, proto string) string {
	if e, ok := r.Resolve(port, proto); ok {
		return e.Name
	}
	return fmt.Sprintf("%d/%s", port, strings.ToLower(proto))
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

type builtin struct {
	port  int
	proto string
	name  string
	desc  string
}

var builtinEntries = []builtin{
	{22, "tcp", "ssh", "Secure Shell"},
	{23, "tcp", "telnet", "Telnet"},
	{25, "tcp", "smtp", "Simple Mail Transfer"},
	{53, "tcp", "dns", "Domain Name System"},
	{53, "udp", "dns", "Domain Name System"},
	{80, "tcp", "http", "Hypertext Transfer Protocol"},
	{110, "tcp", "pop3", "Post Office Protocol v3"},
	{143, "tcp", "imap", "Internet Message Access Protocol"},
	{443, "tcp", "https", "HTTP Secure"},
	{3306, "tcp", "mysql", "MySQL Database"},
	{5432, "tcp", "postgres", "PostgreSQL Database"},
	{6379, "tcp", "redis", "Redis"},
	{8080, "tcp", "http-alt", "HTTP Alternate"},
	{27017, "tcp", "mongodb", "MongoDB"},
}
