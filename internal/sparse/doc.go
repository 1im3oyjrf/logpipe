// Package sparse implements a per-bucket rate-limiting processor for
// structured log entries.
//
// Unlike a global sampler, sparse maintains an independent counter for each
// distinct value of a configurable field (defaulting to "level"). Only the
// first entry in every N consecutive entries within a bucket is forwarded;
// the remainder are silently dropped.
//
// This is useful for suppressing repetitive log lines (e.g. frequent DEBUG
// messages) while still preserving a representative sample, without
// penalising low-volume buckets that would otherwise be starved by a global
// sample rate.
//
// Example usage:
//
//	s, err := sparse.New(sparse.Config{
//		Field: "level",
//		Every: 10,
//		CaseInsensitive: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if s.Allow(entry) {
//		// forward entry downstream
//	}
package sparse
