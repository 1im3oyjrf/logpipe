// Package hover provides a pass-through stage that emits a copy of every
// log entry to a side-channel observer without affecting the main pipeline.
// It is useful for attaching metrics collectors, alerting hooks, or debug
// listeners to an existing pipeline without modifying its topology.
package hover

import (
	"context"

	"github.com/logpipe/logpipe/internal/reader"
)

// Observer is a function that receives a read-only view of each entry.
// Implementations must not block; slow observers should buffer internally.
type Observer func(entry reader.Entry)

// Hover forwards every entry from in to out unchanged, calling obs for each
// one. If obs is nil the stage acts as a transparent pass-through.
type Hover struct {
	obs Observer
}

// New creates a Hover stage with the given observer.
func New(obs Observer) *Hover {
	return &Hover{obs: obs}
}

// Run reads from in, notifies the observer, and writes to out.
// It returns when ctx is cancelled or in is closed.
func (h *Hover) Run(ctx context.Context, in <-chan reader.Entry, out chan<- reader.Entry) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-in:
			if !ok {
				return
			}
			if h.obs != nil {
				h.obs(entry)
			}
			select {
			case out <- entry:
			case <-ctx.Done():
				return
			}
		}
	}
}
