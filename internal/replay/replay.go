// Package replay provides support for replaying historical log entries
// from a ring buffer, optionally filtered by a pattern.
package replay

import (
	"context"

	"github.com/logpipe/internal/buffer"
	"github.com/logpipe/internal/filter"
	"github.com/logpipe/internal/reader"
)

// Config holds options for a replay operation.
type Config struct {
	// Pattern is an optional grep pattern applied during replay.
	Pattern string
	// CaseSensitive controls whether pattern matching is case-sensitive.
	CaseSensitive bool
}

// Replayer replays buffered log entries to a channel.
type Replayer struct {
	buf *buffer.Ring
	cfg Config
}

// New creates a new Replayer backed by the provided ring buffer.
func New(buf *buffer.Ring, cfg Config) *Replayer {
	return &Replayer{buf: buf, cfg: cfg}
}

// Replay sends all matching buffered entries to the returned channel.
// The channel is closed once all entries have been sent or ctx is cancelled.
func (r *Replayer) Replay(ctx context.Context) <-chan reader.Entry {
	out := make(chan reader.Entry, 64)
	go func() {
		defer close(out)
		f := filter.New(filter.Config{
			Pattern:       r.cfg.Pattern,
			CaseSensitive: r.cfg.CaseSensitive,
		})
		for _, entry := range r.buf.Snapshot() {
			if !f.Match(entry) {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- entry:
			}
		}
	}()
	return out
}
