// Package join provides a log-entry transformer that concatenates the values
// of multiple fields into a single new field.
//
// # Usage
//
//	j := join.New(join.Config{
//		Fields:      []string{"service", "host"},
//		Separator:   "@",
//		Target:      "origin",
//		DropSources: true,
//	})
//	out := j.Apply(entry)
//
// Missing or non-string source fields are silently skipped.
// The original entry map is never mutated.
package join
