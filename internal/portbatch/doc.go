// Package portbatch provides a size- and time-bounded batching mechanism for
// port sets produced by the scanner.
//
// Ports are accumulated in an internal buffer. The buffer is flushed when:
//   - the number of buffered ports reaches BatchSize, or
//   - Flush is called explicitly.
//
// A Processor callback is invoked on every flush with the collected Batch.
//
// Example:
//
//	cfg := portbatch.DefaultConfig()
//	b := portbatch.New(cfg, func(batch portbatch.Batch) {
//		fmt.Printf("flushed %d ports\n", len(batch.Ports))
//	})
//	b.AddAll(ports)
//	b.Flush()
package portbatch
