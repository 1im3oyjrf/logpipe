// Package ceiling implements a rate-cap processor for log entry streams.
//
// A Ceiling enforces a hard upper bound on the number of entries forwarded
// within a rolling time window. Entries that arrive after the cap has been
// reached are silently dropped and counted via Dropped().
//
// Typical usage:
//
//	c := ceiling.New(ceiling.Config{Max: 1000, Window: time.Second})
//	for entry := range source {
//		if c.Allow(entry) {
//			// forward entry
//		}
//	}
//
// Alternatively, use Apply to wrap a channel directly:
//
//	filtered := c.Apply(entryChan)
package ceiling
