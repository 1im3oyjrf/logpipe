// Package debounce implements a log pipeline stage that collapses rapid bursts
// of identical log entries into a single representative entry.
//
// When multiple entries sharing the same key field arrive within the configured
// quiet period, only the first is forwarded. A "debounce_count" field is
// injected when more than one occurrence was suppressed, giving downstream
// consumers visibility into how many lines were collapsed.
//
// Typical usage:
//
//	d := debounce.New(debounce.Config{
//		QuietPeriod: 100 * time.Millisecond,
//		KeyField:    "message",
//	})
//	out := d.Run(ctx, in)
package debounce
