// Package copy implements a pipeline stage that duplicates each incoming log
// entry, forwarding both the original and a shallow clone to the downstream
// channel. Optional field overrides can be applied exclusively to the clone,
// making it easy to tag copied entries with a different source, environment,
// or routing label without mutating the original.
//
// Usage:
//
//	copier := copy.New(copy.Config{
//		Overrides: map[string]string{"source": "mirror"},
//	})
//	copier.Run(ctx, in, out)
package copy
