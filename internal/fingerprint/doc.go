// Package fingerprint provides a log-entry fingerprinting stage for the
// logpipe processing pipeline.
//
// A Fingerprinter hashes a configurable subset of entry fields using
// SHA-256 and writes a short hexadecimal digest into a nominated output
// field (default: "_fp"). Downstream stages can use this value to group
// related entries, drive deduplication, or route traffic by identity.
//
// # Usage
//
//	fp := fingerprint.New(fingerprint.Config{
//		Fields:      []string{"level", "message"},
//		OutputField: "_fp",
//	})
//	enriched := fp.Apply(entry)
//
// When Fields is empty every field present in the entry is included in
// the hash, sorted alphabetically for stability.
package fingerprint
