// Package ratelimit provides a token-bucket style rate limiter for
// controlling the throughput of log entries through the pipeline.
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter controls the rate at which log entries are processed.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to ratePerSec entries per second
// with a burst capacity of burstSize.
func New(ratePerSec float64, burstSize int) *Limiter {
	now := time.Now()
	return &Limiter{
		tokens:   float64(burstSize),
		max:      float64(burstSize),
		rate:     ratePerSec,
		lastTick: now,
		clock:    time.Now,
	}
}

// Wait blocks until a token is available or the context is cancelled.
// Returns ctx.Err() if the context is cancelled before a token is acquired.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if l.tryConsume() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(l.nextAvailable()):
		}
	}
}

// tryConsume attempts to consume one token. Returns true on success.
func (l *Limiter) tryConsume() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}
	l.lastTick = now

	if l.tokens >= 1.0 {
		l.tokens--
		return true
	}
	return false
}

// nextAvailable returns the approximate duration until the next token is ready.
func (l *Limiter) nextAvailable() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.rate <= 0 {
		return time.Second
	}
	need := 1.0 - l.tokens
	if need <= 0 {
		return 0
	}
	return time.Duration(need/l.rate*float64(time.Second)) + time.Millisecond
}
