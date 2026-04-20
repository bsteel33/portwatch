// Package portcache provides a thread-safe, TTL-based in-memory cache for
// port scan results produced by the scanner package.
//
// Use portcache to avoid issuing back-to-back scans when multiple components
// request the current port list within a short window. Once the TTL elapses
// the next call to Get returns a cache miss and the caller is expected to
// perform a fresh scan and call Set with the new results.
//
// Example:
//
//	cfg := portcache.DefaultConfig()
//	c := portcache.New(cfg)
//
//	if ports, ok := c.Get(); ok {
//		// use cached ports
//	} else {
//		ports, _ = scanner.OpenPorts()
//		c.Set(ports)
//	}
package portcache
