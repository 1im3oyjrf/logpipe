// Package condense provides a Condenser that collapses consecutive duplicate
// log entries into a single annotated entry.
//
// Two entries are considered duplicates when they share the same normalised
// "message" and "level" values. When a new distinct entry arrives the
// previously buffered entry is emitted, optionally carrying a "_repeat" field
// (or a custom field name) that records how many times it was seen.
//
// Usage:
//
//	c := condense.New(condense.Config{})
//	for _, e := range entries {
//		if out, ok := c.Apply(e); ok {
//			process(out)
//		}
//	}
//	if out, ok := c.Flush(); ok {
//		process(out)
//	}
package condense
