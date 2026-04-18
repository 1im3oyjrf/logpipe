// Package distinct provides a processor that forwards only entries whose
// configured field value has not been seen before within the current session.
package distinct

import (
	"strings"
	"sync"

	"github.com/your-org/logpipe/internal/parser"
)

// Entry mirrors the minimal log-entry shape used across the project.
type Entry = map[string]any

// Config controls which field is used as the uniqueness key.
type Config struct {
	// Field is the entry field whose value must be unique (case-insensitive).
	// Defaults to "message".
	Field string
}

// Processor drops entries whose key-field value has already been seen.
type Processor struct {
	field string
	mu    sync.Mutex
	seen  map[string]struct{}
}

// New returns a Processor configured by cfg.
func New(cfg Config) *Processor {
	f := strings.ToLower(strings.TrimSpace(cfg.Field))
	if f == "" {
		f = "message"
	}
	return &Processor{
		field: f,
		seen:  make(map[string]struct{}),
	}
}

// Apply returns the entry and true when the field value is seen for the first
// time, or the entry and false when it is a duplicate.
func (p *Processor) Apply(e Entry) (Entry, bool) {
	v := parser.GetString(e, p.field)
	if v == "" {
		// No value — always forward to avoid silently dropping entries.
		return e, true
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.seen[v]; exists {
		return e, false
	}
	p.seen[v] = struct{}{}
	return e, true
}

// Reset clears the set of seen values, allowing duplicates to pass again.
func (p *Processor) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.seen = make(map[string]struct{})
}

// Len returns the number of distinct values observed so far.
func (p *Processor) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.seen)
}
