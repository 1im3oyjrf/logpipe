// Package limiter provides a concurrency limiter that caps the number of
// goroutines processing log entries simultaneously.
package limiter

import (
	"context"
	"errors"
)

// ErrLimitExceeded is returned when the semaphore is full and the context is
// already done.
var ErrLimitExceeded = errors.New("limiter: concurrency limit exceeded")

// Limiter is a semaphore-based concurrency limiter.
type Limiter struct {
	sem chan struct{}
}

// New creates a Limiter that allows at most n concurrent acquisitions.
// Panics if n < 1.
func New(n int) *Limiter {
	if n < 1 {
		panic("limiter: n must be >= 1")
	}
	return &Limiter{sem: make(chan struct{}, n)}
}

// Acquire blocks until a slot is available or ctx is done.
// Returns ErrLimitExceeded if ctx is cancelled while waiting.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ErrLimitExceeded
	}
}

// Release frees one slot. Must be called after each successful Acquire.
func (l *Limiter) Release() {
	<-l.sem
}

// Available returns the number of free slots at this instant.
func (l *Limiter) Available() int {
	return cap(l.sem) - len(l.sem)
}

// Cap returns the maximum concurrency configured for this Limiter.
func (l *Limiter) Cap() int {
	return cap(l.sem)
}
