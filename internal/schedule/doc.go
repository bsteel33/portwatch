// Package schedule provides a simple periodic scheduler used by portwatch
// to trigger port scans at a configured interval.
//
// Usage:
//
//	cfg := schedule.DefaultConfig()
//	sched := schedule.New(cfg)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	sched.Run(ctx, func() {
//		// scan and alert
//	})
package schedule
