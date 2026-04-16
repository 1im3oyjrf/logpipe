// Package mask replaces the values of nominated log entry fields with a
// fixed placeholder string, preventing sensitive data such as passwords,
// tokens, and API keys from appearing in rendered output.
//
// Usage:
//
//	m := mask.New(mask.Config{
//		Fields:      []string{"password", "token"},
//		Placeholder: "[MASKED]",
//	})
//	clean := m.Apply(entry)
//
// Field matching is case-insensitive. The original entry map is never
// modified; Apply always returns a shallow copy.
package mask
