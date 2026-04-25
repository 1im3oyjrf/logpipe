// Package sparse provides a processor that emits only every Nth entry
// matching a given field value, discarding the rest. Unlike the sampler
// (which operates on a global rate), sparse applies an independent counter
// per distinct field value so that low-volume keys are not starved.
package sparse

import (
	"errors"
	"strings"
	"sync"

	"github.com/logpipe/logpipe/internal/parser"
)

// Entry is the minimal log-entry type used across the pipeline.
type Entry = map[string]any

// Config controls the behaviour of the Sparse processor.
type Config struct {
	// Field is the entry key whose value is used to bucket counters.
	// Defaults to "level" when empty.
	Field string

	// Every N controls how many entries are kept per bucket.
	// A value of 1 (or 0, treated as 1) keeps every entry.
	Every int

	// CaseInsensitive folds the field value to lower-case before bucketing.
	CaseInsensitive bool
}

// Sparse drops all but every Nth entry within each field-value bucket.
type Sparse struct {
	field           string
	every           int
	caseInsensitive bool

	mu       sync.Mutex
	counters map[string]int
}

const defaultField = "level"
const defaultEvery = 1

// New creates a Sparse processor from cfg.
// Returns an error if Every is negative.
func New(cfg Config) (*Sparse, error) {
	if cfg.Every < 0 {
		return nil, errors.New("sparse: Every must be >= 0")
	}

	field := cfg.Field
	if field == "" {
		field = defaultField
	}

	every := cfg.Every
	if every == 0 {
		every = defaultEvery
	}

	return &Sparse{
		field:           field,
		every:           every,
		caseInsensitive: cfg.CaseInsensitive,
		counters:        make(map[string]int),
	}, nil
}

// Allow returns true if the entry should be forwarded.
func (s *Sparse) Allow(e Entry) bool {
	raw := parser.GetString(e, s.field)
	key := raw
	if s.caseInsensitive {
		key = strings.ToLower(raw)
	}

	s.mu.Lock()
	s.counters[key]++
	n := s.counters[key]
	s.mu.Unlock()

	return n%s.every == 1
}

// Reset clears all per-bucket counters.
func (s *Sparse) Reset() {
	s.mu.Lock()
	s.counters = make(map[string]int)
	s.mu.Unlock()
}
