// Package normalize provides field-name and log-level normalisation for
// structured log entries.
//
// Use New to create a Normalizer from a Config that specifies:
//
//   - FieldMap  – a mapping from raw field names (matched case-insensitively)
//     to canonical names that the rest of the pipeline expects.
//   - LevelField – the key that holds the severity level (default: "level").
//
// Apply returns a new map with renamed keys and the level value folded to
// lower-case. The original entry is never modified.
package normalize
