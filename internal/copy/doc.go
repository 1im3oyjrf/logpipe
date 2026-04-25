// Package copy implements a pipeline stage that duplicates each incoming log
// entry, forwarding both the original and a shallow clone to the downstream
// channel. Optional field overrides can be applied exclusively to the clone,
// making it easy to tag copied entries with a different source, environment,
// or routing label without mutating the original.
//
// The stage guarantees that the original entry is always forwarded before its
// clone, preserving relative ordering within the output channel. If the
// context is cancelled, both sends are abandoned cleanly without blocking.
//
// Usage:
//
//	copier := copy.New(copy.Config{
//		Overrides: map[string]string{"source": "mirror"},
//	})
//	copier.Run(ctx, in, out)
package copy
