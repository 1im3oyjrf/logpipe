// Package splitter fans out log entries to multiple named output channels
// based on a fixed set of registered targets.
package splitter

import (
	"context"
	"sync"

	"github.com/logpipe/logpipe/internal/parser"
)

// Target is a named output channel for log entries.
type Target struct {
	Name string
	Ch   chan map[string]any
}

// Splitter distributes every received entry to all registered targets.
type Splitter struct {
	mu      sync.RWMutex
	targets []*Target
}

// New returns an empty Splitter.
func New() *Splitter {
	return &Splitter{}
}

// Add registers a new named target. The caller is responsible for creating
// and draining the channel.
func (s *Splitter) Add(name string, ch chan map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.targets = append(s.targets, &Target{Name: name, Ch: ch})
}

// Targets returns a snapshot of currently registered targets.
func (s *Splitter) Targets() []*Target {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Target, len(s.targets))
	copy(out, s.targets)
	return out
}

// Run reads entries from in and fans each one out to every registered target.
// It returns when ctx is cancelled or in is closed.
func (s *Splitter) Run(ctx context.Context, in <-chan map[string]any) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-in:
			if !ok {
				return
			}
			s.dispatch(ctx, entry)
		}
	}
}

func (s *Splitter) dispatch(ctx context.Context, entry map[string]any) {
	s.mu.RLock()
	targets := s.targets
	s.mu.RUnlock()

	for _, t := range targets {
		// shallow copy so downstream mutations don't interfere
		copy := parser.ShallowCopy(entry)
		select {
		case t.Ch <- copy:
		case <-ctx.Done():
			return
		}
	}
}
