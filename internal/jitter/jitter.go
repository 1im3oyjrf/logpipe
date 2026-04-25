// Package jitter adds a small random delay field to each log entry,
// useful for simulating or annotating processing latency in test pipelines.
package jitter

import (
	"math/rand"
	"time"

	"github.com/your-org/logpipe/internal/parser"
)

// Config controls jitter behaviour.
type Config struct {
	// Field is the output field name. Defaults to "jitter_ms".
	Field string
	// MaxMS is the upper bound (exclusive) for the random value in milliseconds.
	// Defaults to 100.
	MaxMS int
	// Overwrite replaces an existing field when true.
	Overwrite bool
}

// Applier injects a random jitter value into each entry.
type Applier struct {
	field     string
	maxMS     int
	overwrite bool
	rng       *rand.Rand
}

// New returns a configured Applier.
func New(cfg Config) *Applier {
	field := cfg.Field
	if field == "" {
		field = "jitter_ms"
	}
	maxMS := cfg.MaxMS
	if maxMS <= 0 {
		maxMS = 100
	}
	return &Applier{
		field:     field,
		maxMS:     maxMS,
		overwrite: cfg.Overwrite,
		//nolint:gosec // non-cryptographic, intentional
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Apply returns a copy of entry with the jitter field injected.
func (a *Applier) Apply(entry map[string]any) map[string]any {
	key := parser.CanonicalKey(a.field, entry)
	if key == "" {
		key = a.field
	}
	if _, exists := entry[key]; exists && !a.overwrite {
		return entry
	}
	out := shallowCopy(entry)
	out[key] = a.rng.Intn(a.maxMS)
	return out
}

func shallowCopy(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
