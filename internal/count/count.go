// Package count provides a processor that injects a running total of
// entries seen into each log entry as a numeric field.
package count

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/your-org/logpipe/internal/parser"
)

const (
	defaultField = "_count"
)

// Config controls the behaviour of the counter processor.
type Config struct {
	// Field is the name of the field injected into each entry.
	// Defaults to "_count".
	Field string

	// Overwrite replaces an existing field with the same name when true.
	// When false (the default) the entry is passed through unchanged if
	// the field already exists.
	Overwrite bool
}

// Processor injects a monotonically increasing counter into every entry.
type Processor struct {
	field     string
	overwrite bool
	counter   atomic.Uint64
}

// New returns a Processor configured by cfg.
func New(cfg Config) *Processor {
	f := strings.TrimSpace(cfg.Field)
	if f == "" {
		f = defaultField
	}
	return &Processor{
		field:     f,
		overwrite: cfg.Overwrite,
	}
}

// Apply injects the current counter value into entry and returns the
// modified copy. The original entry is never mutated.
func (p *Processor) Apply(entry map[string]any) map[string]any {
	key := p.field

	if !p.overwrite {
		if _, exists := parser.HasField(entry, key); exists {
			return entry
		}
	}

	n := p.counter.Add(1)

	out := make(map[string]any, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}
	out[key] = fmt.Sprintf("%d", n)
	return out
}

// Reset sets the internal counter back to zero.
func (p *Processor) Reset() {
	p.counter.Store(0)
}

// Value returns the current counter value without incrementing it.
func (p *Processor) Value() uint64 {
	return p.counter.Load()
}
