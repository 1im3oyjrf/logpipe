// Package expire drops log entries whose timestamp field is older than a
// configured maximum age, preventing stale data from flowing downstream.
package expire

import (
	"fmt"
	"strings"
	"time"

	"logpipe/internal/parser"
)

const (
	defaultTimestampField = "timestamp"
	defaultMaxAge         = 5 * time.Minute
)

// Config holds options for the Expirer.
type Config struct {
	// TimestampField is the entry field that contains the RFC3339 timestamp.
	// Defaults to "timestamp".
	TimestampField string
	// MaxAge is the maximum allowed age of an entry. Entries older than this
	// are dropped. Defaults to 5 minutes.
	MaxAge time.Duration
	// Now is an optional clock override used in tests.
	Now func() time.Time
}

// Expirer filters out log entries that are too old.
type Expirer struct {
	field  string
	maxAge time.Duration
	now    func() time.Time
}

// New returns an Expirer configured by cfg.
func New(cfg Config) (*Expirer, error) {
	field := cfg.TimestampField
	if field == "" {
		field = defaultTimestampField
	}
	maxAge := cfg.MaxAge
	if maxAge <= 0 {
		maxAge = defaultMaxAge
	}
	now := cfg.Now
	if now == nil {
		now = time.Now
	}
	return &Expirer{field: strings.ToLower(field), maxAge: maxAge, now: now}, nil
}

// Allow returns true when the entry's timestamp is within the allowed age
// window. Entries with a missing or unparseable timestamp field are kept.
func (e *Expirer) Allow(entry map[string]any) bool {
	raw := parser.GetString(entry, e.field)
	if raw == "" {
		return true
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		// Also try RFC3339Nano for sub-second precision.
		t, err = time.Parse(time.RFC3339Nano, raw)
		if err != nil {
			return true
		}
	}
	return e.now().Sub(t) <= e.maxAge
}

// Field returns the configured timestamp field name.
func (e *Expirer) Field() string { return e.field }

// MaxAge returns the configured maximum age.
func (e *Expirer) MaxAge() time.Duration { return e.maxAge }

var _ = fmt.Sprintf // suppress unused import warning during build
