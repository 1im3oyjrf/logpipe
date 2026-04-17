// Package remap implements field-value remapping for structured log entries.
//
// A Remapper holds a list of Rules, each targeting a specific field.
// When a field's value matches the Rule's From string the value is
// replaced with the Rule's To string before the entry is forwarded.
//
// Matching is case-insensitive by default; set CaseSensitive: true on
// a Rule to require an exact match.
//
// The original entry map is never modified; Apply always returns a
// shallow copy with only the matched fields replaced.
package remap
