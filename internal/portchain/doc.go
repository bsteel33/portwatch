// Package portchain implements a composable processing pipeline for port lists.
//
// A Chain holds an ordered sequence of named Stage functions. Each stage
// receives the output of the previous stage, allowing callers to compose
// filtering, labelling, ranking, and any other transformation in a single
// readable pipeline:
//
//	c := portchain.New(verbose)
//	c.Add("filter",  filter.Apply)
//	c.Add("rank",    rank.Apply)
//	c.Add("label",   label.Apply)
//	result := c.Run(scanned)
//
// Stages are executed in registration order. When verbose mode is enabled
// the chain prints per-stage port counts to stderr for debugging.
package portchain
