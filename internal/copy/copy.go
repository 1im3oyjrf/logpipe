// Package copy provides a stage that duplicates each log entry,
// emitting the original plus a deep copy with optional field overrides.
package copy

import (
	"context"

	"logpipe/internal/parser"
)

// Entry mirrors the shared log entry type used across the pipeline.
type Entry = map[string]any

// Config controls how copies are produced.
type Config struct {
	// Overrides are injected into the copy only, not the original.
	Overrides map[string]string
}

// Copier duplicates entries flowing through the pipeline.
type Copier struct {
	cfg Config
}

// New returns a Copier configured with cfg.
func New(cfg Config) *Copier {
	overrides := make(map[string]string, len(cfg.Overrides))
	for k, v := range cfg.Overrides {
		overrides[parser.Lower(k)] = v
	}
	cfg.Overrides = overrides
	return &Copier{cfg: cfg}
}

// Run reads from in, emits each entry followed by its copy to out, and
// closes out when in is drained or ctx is cancelled.
func (c *Copier) Run(ctx context.Context, in <-chan Entry, out chan<- Entry) {
	defer close(out)
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-in:
			if !ok {
				return
			}
			select {
			case out <- e:
			case <-ctx.Done():
				return
			}
			cp := clone(e)
			for k, v := range c.cfg.Overrides {
				cp[k] = v
			}
			select {
			case out <- cp:
			case <-ctx.Done():
				return
			}
		}
	}
}

func clone(e Entry) Entry {
	out := make(Entry, len(e))
	for k, v := range e {
		out[k] = v
	}
	return out
}
