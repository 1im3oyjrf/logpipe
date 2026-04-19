package sequence

import (
	"sync/atomic"

	"github.com/logpipe/logpipe/internal/parser"
)

// Sequencer attaches a monotonically increasing counter to each log entry.
type Sequencer struct {
	field string
	counter uint64
}

// Config holds options for the Sequencer.
type Config struct {
	// Field is the key injected into each entry. Defaults to "_seq".
	Field string
}

// New returns a Sequencer that stamps entries with a sequence number.
func New(cfg Config) *Sequencer {
	if cfg.Field == "" {
		cfg.Field = "_seq"
	}
	return &Sequencer{field: cfg.Field}
}

// Apply injects the next sequence number into a copy of the entry.
func (s *Sequencer) Apply(entry map[string]any) map[string]any {
	n := atomic.AddUint64(&s.counter, 1)
	out := shallowCopy(entry)
	out[s.field] = int64(n)
	return out
}

// Reset resets the internal counter to zero.
func (s *Sequencer) Reset() {
	atomic.StoreUint64(&s.counter, 0)
}

// Field returns the configured field name.
func (s *Sequencer) Field() string { return s.field }

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src)+1)
	for k, v := range src {
		out[k] = v
	}
	_ = parser.Keys // ensure import is used
	return out
}
