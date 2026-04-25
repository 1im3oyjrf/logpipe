// Package jitter provides an Applier that injects a random integer
// millisecond value into each log entry under a configurable field name.
//
// This is useful when testing downstream components that consume latency
// annotations, or when you want to simulate variable processing delays in
// a pipeline without modifying real timestamps.
//
// Usage:
//
//	a := jitter.New(jitter.Config{
//		Field: "proc_jitter_ms",
//		MaxMS: 250,
//	})
//	out := a.Apply(entry)
package jitter
