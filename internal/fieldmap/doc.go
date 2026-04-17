// Package fieldmap rewrites field names in structured log entries according
// to a user-supplied mapping table.
//
// Use cases include normalising vendor-specific field names (e.g. "@message"
// → "msg"), enforcing a canonical schema before forwarding entries to an
// output sink, or stripping internal fields that should not be exposed.
//
// Example:
//
//	m := fieldmap.New(fieldmap.Config{
//		Rules: map[string]string{
//			"@message": "msg",
//			"@timestamp": "time",
//		},
//		DropUnmapped: false,
//	})
//	outFields := m.Apply(entry.Fields)
package fieldmap
