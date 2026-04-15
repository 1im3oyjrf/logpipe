// Package transform provides field-level transformations for structured log entries.
//
// A Transformer can:
//   - Redact sensitive field values (e.g. passwords, tokens) by replacing them
//     with the literal string "[REDACTED]".
//   - Rename fields to normalise key names across different log sources.
//   - Inject static key/value pairs into every entry (e.g. environment tags).
//
// Transformations are applied after filtering and before output formatting,
// ensuring that sensitive data is never written to the output sink.
//
// Example:
//
//	tr := transform.New(transform.Config{
//		RedactFields: []string{"password", "token"},
//		RenameFields: map[string]string{"msg": "message"},
//		AddFields:    map[string]string{"env": "production"},
//	})
//	clean := tr.Apply(entry)
package transform
