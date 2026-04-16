// Package parser provides field extraction and value coercion utilities
// for structured log entries. It handles type-safe access to dynamic
// map[string]any fields commonly found in parsed JSON log lines.
package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// Entry represents a parsed log entry as a flat key-value map.
type Entry = map[string]any

// GetString returns the string value for key, searching case-insensitively.
// Returns empty string and false if the key is absent or the value is not a string.
func GetString(e Entry, key string) (string, bool) {
	if v, ok := lookup(e, key); ok {
		switch s := v.(type) {
		case string:
			return s, true
		case fmt.Stringer:
			return s.String(), true
		}
	}
	return "", false
}

// GetFloat returns the float64 value for key.
// Handles float64, int, int64, and string representations.
func GetFloat(e Entry, key string) (float64, bool) {
	v, ok := lookup(e, key)
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case string:
		f, err := strconv.ParseFloat(n, 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}

// GetBool returns the boolean value for key.
// Handles bool and string representations ("true", "1", "false", "0").
func GetBool(e Entry, key string) (bool, bool) {
	v, ok := lookup(e, key)
	if !ok {
		return false, false
	}
	switch b := v.(type) {
	case bool:
		return b, true
	case string:
		parsed, err := strconv.ParseBool(b)
		if err == nil {
			return parsed, true
		}
	}
	return false, false
}

// Keys returns all keys present in the entry in their original casing.
func Keys(e Entry) []string {
	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}
	return keys
}

// HasField reports whether the entry contains the given key (case-insensitive).
func HasField(e Entry, key string) bool {
	_, ok := lookup(e, key)
	return ok
}

// lookup performs a case-insensitive key search on the entry.
func lookup(e Entry, key string) (any, bool) {
	// Exact match first for performance.
	if v, ok := e[key]; ok {
		return v, true
	}
	lower := strings.ToLower(key)
	for k, v := range e {
		if strings.ToLower(k) == lower {
			return v, true
		}
	}
	return nil, false
}
