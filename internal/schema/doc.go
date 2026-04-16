// Package schema implements structural validation for JSON log entries.
//
// A Validator is constructed from a Config that specifies:
//
//   - RequiredFields: fields that must appear in every entry.
//   - KnownFields: the exhaustive whitelist of permitted field names
//     (only enforced when AllowUnknown is false).
//   - AllowUnknown: when true, fields outside KnownFields are silently
//     accepted; when false they are treated as violations.
//
// All field name comparisons are case-insensitive.
//
// Example:
//
//	v := schema.New(schema.Config{
//		RequiredFields: []string{"level", "message", "ts"},
//		AllowUnknown:   true,
//	})
//	if err := v.Validate(entry.Fields); err != nil {
//		log.Printf("dropping invalid entry: %v", err)
//	}
package schema
