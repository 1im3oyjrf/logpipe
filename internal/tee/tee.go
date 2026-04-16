// Package tee provides a fan-out writer that duplicates log entries
// to multiple downstream channels simultaneously.
package tee

import (
	"context"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Tee fans out entries from a single source channel to multiple consumers.
type Tee struct {
	mu      sync.RWMutex
	targets []chan reader.Entry
}

// New creates a new Tee with the given per-target buffer size.
func New(bufSize int) *Tee {
	_ = bufSize // stored for future use when adding targets
	return &Tee{}
}

// Add registers a new consumer channel and returns it.
func (t *Tee) Add(bufSize int) <-chan reader.Entry {
	ch := make(chan reader.Entry, bufSize)
	t.mu.Lock()
	t.targets = append(t.targets, ch)
	t.mu.Unlock()
	return ch
}

// Run reads from src and copies each entry to all registered targets.
// It closes all target channels when src is drained or ctx is cancelled.
func (t *Tee) Run(ctx context.Context, src <-chan reader.Entry) {
	defer t.closeAll()
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-src:
			if !ok {
				return
			}
			t.mu.RLock()
			for _, ch := range t.targets {
				select {
				case ch <- entry:
				default:
					// drop if consumer is full
				}
			}
			t.mu.RUnlock()
		}
	}
}

func (t *Tee) closeAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, ch := range t.targets {
		close(ch)
	}
}
