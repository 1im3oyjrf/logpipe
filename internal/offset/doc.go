// Package offset provides a log-entry processor that shifts a named numeric
// field by a fixed delta.
//
// # Usage
//
//	p := offset.New(offset.Config{
//		Field: "duration_ms",
//		By:    -100,          // subtract 100 from every duration
//	})
//	processed := p.Apply(entry)
//
// The original entry is never mutated; Apply always returns a shallow copy
// when a modification is made. If the target field differs from the source
// field the source key is removed from the output.
package offset
