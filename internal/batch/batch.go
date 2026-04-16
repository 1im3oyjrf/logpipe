package batch

import (
	"context"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Config holds configuration for the batcher.
type Config struct {
	MaxSize int
	MaxWait time.Duration
}

// Batcher accumulates log entries and flushes them as slices either when
// the batch reaches MaxSize or MaxWait elapses, whichever comes first.
type Batcher struct {
	cfg Config
	in  <-chan reader.Entry
}

// New returns a Batcher that reads from in.
func New(cfg Config, in <-chan reader.Entry) *Batcher {
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = 100
	}
	if cfg.MaxWait <= 0 {
		cfg.MaxWait = 500 * time.Millisecond
	}
	return &Batcher{cfg: cfg, in: in}
}

// Run reads entries from the input channel and emits batches on the returned
// channel. The returned channel is closed when ctx is cancelled.
func (b *Batcher) Run(ctx context.Context) <-chan []reader.Entry {
	out := make(chan []reader.Entry)
	go func() {
		defer close(out)
		buf := make([]reader.Entry, 0, b.cfg.MaxSize)
		ticker := time.NewTicker(b.cfg.MaxWait)
		defer ticker.Stop()
		flush := func() {
			if len(buf) == 0 {
				return
			}
			snap := make([]reader.Entry, len(buf))
			copy(snap, buf)
			select {
			case out <- snap:
			case <-ctx.Done():
			}
			buf = buf[:0]
		}
		for {
			select {
			case <-ctx.Done():
				flush()
				return
			case e, ok := <-b.in:
				if !ok {
					flush()
					return
				}
				buf = append(buf, e)
				if len(buf) >= b.cfg.MaxSize {
					flush()
					ticker.Reset(b.cfg.MaxWait)
				}
			case <-ticker.C:
				flush()
			}
		}
	}()
	return out
}
