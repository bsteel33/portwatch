// Package ratelimit provides a sliding-window rate limiter for portwatch
// event pipelines. It is used to suppress bursts of alerts when many port
// changes are detected in a short period, ensuring downstream consumers
// (notify, report) are not overwhelmed.
//
// Usage:
//
//	cfg := ratelimit.DefaultConfig()
//	limiter := ratelimit.New(cfg)
//
//	if limiter.Allow() {
//		// forward event
//	}
package ratelimit
