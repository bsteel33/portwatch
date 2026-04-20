// Package portclassify assigns a risk classification (safe, monitor,
// suspicious, dangerous) to open ports discovered by the scanner.
//
// Classification is rule-based: each Rule matches on port number and/or
// protocol and maps to a Class with a human-readable reason string.
// Rules are evaluated in order; the first match wins.
//
// Usage:
//
//	cfg := portclassify.DefaultConfig()
//	cl  := portclassify.New(cfg)
//	results := cl.Classify(ports)
//	portclassify.Print(results)
package portclassify
