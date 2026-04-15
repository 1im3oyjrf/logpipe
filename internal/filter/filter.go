// Package filter provides grep-style filtering for structured JSON log entries.
package filter

import (
	"strings"
)

// Entry represents a parsed JSON log entry as a key-value map.
type Entry map[string]interface{}

// Filter holds the configuration for filtering log entries.
type Filter struct {
	// pattern is the substring or keyword to search for.
	pattern string
	// fields restricts matching to specific fields; if empty, all fields are searched.
	fields []string
	// caseSensitive controls whether matching is case-sensitive.
	caseSensitive bool
}

// New creates a new Filter with the given pattern and options.
func New(pattern string, fields []string, caseSensitive bool) *Filter {
	return &Filter{
		pattern:       pattern,
		fields:        fields,
		caseSensitive: caseSensitive,
	}
}

// Match reports whether the given log entry matches the filter pattern.
// If no pattern is set, all entries match.
func (f *Filter) Match(entry Entry) bool {
	if f.pattern == "" {
		return true
	}

	pattern := f.pattern
	if !f.caseSensitive {
		pattern = strings.ToLower(pattern)
	}

	if len(f.fields) > 0 {
		for _, field := range f.fields {
			val, ok := entry[field]
			if !ok {
				continue
			}
			if f.matchValue(val, pattern) {
				return true
			}
		}
		return false
	}

	for _, val := range entry {
		if f.matchValue(val, pattern) {
			return true
		}
	}
	return false
}

// matchValue checks whether a single value contains the pattern.
func (f *Filter) matchValue(val interface{}, pattern string) bool {
	str, ok := val.(string)
	if !ok {
		str = strings.TrimSpace(strings.Replace(strings.Replace(fmt.Sprintf("%v", val), "<nil>", "", -1), "\n", " ", -1))
	}
	if !f.caseSensitive {
		str = strings.ToLower(str)
	}
	return strings.Contains(str, pattern)
}
