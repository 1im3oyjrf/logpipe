// Package metrics provides lightweight, thread-safe counters for tracking
// logpipe pipeline statistics at runtime.
//
// A Counters value records the number of log lines read, matched by the active
// filter, dropped (read but not matched), and lines that failed JSON parsing.
// All counter operations are safe for concurrent use via sync/atomic.
//
// Usage:
//
//	c := metrics.New()
//	c.IncRead()
//	c.IncMatched()
//	snap := c.Snapshot() // point-in-time copy
//	fmt.Println(snap.LinesRead, snap.Uptime)
package metrics
