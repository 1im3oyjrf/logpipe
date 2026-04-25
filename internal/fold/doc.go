// Package fold provides a log entry transformer that collapses multiple
// string fields into a single target field.
//
// # Overview
//
// A Folder reads a list of source fields from each log entry, joins their
// string values with a configurable separator, and writes the result to a
// target field. Source fields that are absent or empty are silently skipped.
//
// # Configuration
//
//   - Fields       – two or more source field names (required).
//   - Target       – destination field name (required).
//   - Separator    – string placed between values; defaults to a single space.
//   - DropSources  – when true, source fields are removed from the output.
//   - CaseInsensitive – match source field names case-insensitively.
//
// The original entry map is never mutated; Apply always returns a new copy.
package fold
