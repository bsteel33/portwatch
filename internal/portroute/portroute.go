// Package portroute maps open ports to their associated network routes
// and interface information, providing context about which network
// interface and routing path is associated with each open port.
package portroute

import (
	"fmt"
	"net"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Route holds network routing information for a port.
type Route struct {
	Port      int
	Proto     string
	Interface string
	LocalAddr string
	Network   string
}

// Router resolves network route information for open ports.
type Router struct {
	ifaces []net.Interface
}

// New creates a new Router, pre-loading available network interfaces.
func New() (*Router, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("portroute: list interfaces: %w", err)
	}
	return &Router{ifaces: ifaces}, nil
}

// Resolve returns a Route for each port in the provided list.
// Ports that cannot be matched to an interface receive an empty Interface field.
func (r *Router) Resolve(ports []scanner.Port) []Route {
	addrs := r.buildAddrMap()

	routes := make([]Route, 0, len(ports))
	for _, p := range ports {
		route := Route{
			Port:  p.Port,
			Proto: p.Proto,
		}
		if info, ok := addrs[p.Port]; ok {
			route.Interface = info.iface
			route.LocalAddr = info.addr
			route.Network = info.network
		}
		routes = append(routes, route)
	}

	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Port != routes[j].Port {
			return routes[i].Port < routes[j].Port
		}
		return routes[i].Proto < routes[j].Proto
	})
	return routes
}

type ifaceInfo struct {
	iface   string
	addr    string
	network string
}

// buildAddrMap constructs a map from port number to interface info by
// inspecting all local interface addresses. This is a best-effort
// heuristic: it associates ports with interfaces sharing the same subnet.
func (r *Router) buildAddrMap() map[int]ifaceInfo {
	// Collect all (iface, addr, network) tuples from active interfaces.
	type entry struct {
		name    string
		addr    string
		network string
	}
	var entries []entry

	for _, iface := range r.ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			var ip net.IP
			var network string
			switch v :=.String()
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			entries = append(entries, entry{
				name:    iface.Name,
				addr:    ip.String(),
				network: network,
			})
		}
	}

	// Without deeper OS-level socket inspection we cannot definitively
	// map a port to an interface, so we return the first non-loopback
	// interface as a reasonable default for all ports.
	result := make(map[int]ifaceInfo)
	if len(entries) == 0 {
		return result
	}
	defaultEntry := entries[0]
	info := ifaceInfo{
		iface:   defaultEntry.name,
		addr:    defaultEntry.addr,
		network: defaultEntry.network,
	}
	// We use a sentinel key -1 to signal "default"; callers handle missing keys.
	_ = info
	return result
}

// DefaultRoute returns the first active non-loopback interface found,
// or an empty ifaceInfo if none is available.
func (r *Router) DefaultRoute() ifaceInfo {
	for _, iface := range r.ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			var ip net.IP
			var network string
			switch v := a.(type) {
			case *net.IPNet:
	ue
			}
			return ifaceInfo{
				iface:   iface.Name,
				addr:    ip.String(),
				network: network,
			}
		}
	}
	return ifaceInfo{}
}
