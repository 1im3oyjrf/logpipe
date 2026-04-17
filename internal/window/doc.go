// Package window implements a sliding time-window counter and rate guard for
// use in the logpipe pipeline.
//
// Window tracks event counts within a rolling duration, evicting entries that
// fall outside the configured window on each access.
//
// Guard wraps a Window with a configurable limit and exposes a simple
// Allow/Count/Reset API suitable for rate-limiting decisions in filters,
// alerts, and throttle stages.
//
// Example:
//
//	g := window.NewGuard(window.GuardConfig{
//		WindowSize: 10 * time.Second,
//		Limit:      50,
//	})
//	if !g.Allow() {
//		// rate limit exceeded
//	}
package window
