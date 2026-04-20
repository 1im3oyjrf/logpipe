// Package format provides the Formatter processor, which renders a
// Go fmt-style template against a log entry's fields and writes the result
// into a configurable target field.
//
// Placeholders take the form {fieldName} and are replaced with the string
// representation of the corresponding field value. Unknown placeholders are
// left as-is in the output string.
//
// Example usage:
//
//	f, err := format.New(format.Config{
//		Target:   "summary",
//		Template: "[{level}] {message} on {host}",
//	})
//	if err != nil { ... }
//	processed := f.Apply(entry)
package format
