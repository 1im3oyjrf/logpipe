// Package throttle implements per-key emission throttling for log entries.
//
// A Throttle suppresses repeated entries that share the same key within a
// configurable cooldown window. This is useful for noisy log sources that
// emit the same message at a high rate.
//
// Basic usage:
//
//	th := throttle.New(5 * time.Second)
//
//	for entry := range entries {
//		if th.Allow(entry.Message) {
//			forwardEntry(entry)
//		}
//	}
//
// Call Evict periodically in long-running processes to release memory
// for keys that have naturally expired.
package throttle
