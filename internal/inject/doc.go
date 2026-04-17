// Package inject provides a static field injector for log entries.
//
// An Injector is configured with a map of field names to string values.
// When applied to a log entry it copies the entry and adds the configured
// fields, optionally overwriting any fields that already exist.
//
// Field names are normalised to lower-case so that configuration is
// case-insensitive.
//
// Usage:
//
//	inj := inject.New(inject.Config{
//		Fields:            map[string]string{"env": "production"},
//		OverwriteExisting: false,
//	})
//	enriched := inj.Apply(entry)
package inject
