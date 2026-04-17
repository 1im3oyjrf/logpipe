// Package coalesce implements a field-coalescing transformer for structured
// log entries.
//
// Different log producers often emit the same semantic value under different
// field names (e.g. "msg", "message", "text"). The coalesce transformer
// resolves this by scanning a prioritised list of source fields and writing
// the first non-empty value into a single canonical target field.
//
// Example usage:
//
//	tr := coalesce.New(coalesce.Config{
//		Rules: []coalesce.Rule{
//			{Sources: []string{"msg", "message", "text"}, Target: "message"},
//		},
//	})
//	out := tr.Apply(entry)
package coalesce
