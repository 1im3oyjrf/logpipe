// Package truncate caps string field values in structured log entries to a
// configured maximum byte length. This prevents runaway log payloads from
// overwhelming downstream consumers such as the formatter or multiwriter.
//
// Usage:
//
//	tr := truncate.New(truncate.Config{
//		MaxLen: 512,
//		Fields: []string{"message", "error"},
//	})
//	processed := tr.Apply(entry)
package truncate
