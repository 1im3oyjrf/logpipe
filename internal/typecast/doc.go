// Package typecast provides field-level type coercion for structured log entries.
//
// Rules map field names (case-insensitive) to a target type: "string", "int",
// "float", or "bool". Fields whose current value cannot be parsed into the
// target type are left unchanged. The original entry is never mutated.
//
// Example usage:
//
//	c := typecast.New(typecast.Config{
//		Rules: []typecast.Rule{
//			{Field: "status_code", Target: "int"},
//			{Field: "latency_ms",  Target: "float"},
//		},
//	})
//	out := c.Apply(entry)
package typecast
