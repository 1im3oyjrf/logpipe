// Package timeout provides a pipeline stage that enforces a per-entry
// processing deadline. Entries that cannot be forwarded to the downstream
// consumer within the configured duration are dropped, allowing the pipeline
// to remain live under backpressure conditions.
//
// # Usage
//
//	g := timeout.New(timeout.Config{Duration: 100 * time.Millisecond})
//	out := g.Run(ctx, in)
//	// consume out …
//	fmt.Println("dropped:", g.Dropped())
package timeout
