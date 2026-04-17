// Package timeout provides per-entry processing deadline enforcement.
// Entries that exceed the configured duration are dropped and counted.
package timeout

import (
	"context"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
)

// Config holds configuration for the timeout stage.
type Config struct {
	// Duration is the maximum time allowed to process a single entry.
	// Defaults to 200ms if zero.
	Duration time.Duration
}

// Guard enforces a processing deadline on each log entry.
type Guard struct {
	cfg Config
	dropped uint64
}

// New creates a Guard with the provided Config.
// A zero Duration is replaced with the 200 ms default.
func New(cfg Config) *Guard {
	if cfg.Duration <= 0 {
		cfg.Duration = 200 * time.Millisecond
	}
	return &Guard{cfg: cfg}
}

// Run reads entries from in, applies a per-entry deadline, and forwards
// entries that complete in time to the returned channel.
// Entries that exceed the deadline are silently dropped; Dropped reports
// the cumulative count.
func (g *Guard) Run(ctx context.Context, in <-chan reader.Entry) <-chan reader.Entry {
	out := make(chan reader.Entry, 64)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-in:
				if !ok {
					return
				}
				tctx, cancel := context.WithTimeout(ctx, g.cfg.Duration)
				forwarded := g.forward(tctx, out, entry)
				cancel()
				if !forwarded {
					g.dropped++
				}
			}
		}
	}()
	return out
}

func (g *Guard) forward(ctx context.Context, out chan<- reader.Entry, e reader.Entry) bool {
	select {
	case out <- e:
		return true
	case <-ctx.Done():
		return false
	}
}

// Dropped returns the number of entries dropped due to deadline expiry.
func (g *Guard) Dropped() uint64 { return g.dropped }
