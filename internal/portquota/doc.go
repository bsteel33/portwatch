// Package portquota enforces configurable limits on the number of open ports,
// both in total and per protocol (tcp/udp). When a scan result exceeds a
// configured threshold a Violation is returned so callers can alert or block.
package portquota
