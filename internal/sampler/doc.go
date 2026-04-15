// Package sampler provides deterministic log-entry sampling for logpipe.
//
// When processing high-volume log streams it is often desirable to forward
// only a representative subset of entries rather than every line. Sampler
// keeps exactly 1 out of every N entries in arrival order, making the
// behaviour predictable and reproducible across runs.
//
// Usage:
//
//	s := sampler.New(10) // keep 1 in 10
//	if s.Keep(entry) {
//	    output.Write(entry)
//	}
//
// The zero value is not usable; always construct with New.
package sampler
