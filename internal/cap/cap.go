// Package cap limits the total number of log entries forwarded
// through the pipeline, closing the output channel once the cap is reached.
package cap

import (
	"context"

	"github.com/logpipe/logpipe/internal/reader"
)

// Config holds configuration for the entry cap.
type Config struct {
	// Max is the maximum number of entries to forward. Zero means no limit.
	Max int
}

// Capper forwards at most Max entries from in to the returned channel.
type Capper struct {
	cfg Config
}

// New returns a Capper with the given configuration.
func New(cfg Config) *Capper {
	return &Capper{cfg: cfg}
}

// Run reads from in and writes to the returned channel, stopping after Max
// entries have been forwarded. If Max is zero all entries are forwarded until
// in is closed or ctx is cancelled.
func (c *Capper) Run(ctx context.Context, in <-chan reader.Entry) <-chan reader.Entry {
	out := make(chan reader.Entry)
	go func() {
		defer close(out)
		count := 0
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
				count++
				if c.cfg.Max > 0 && count >= c.cfg.Max {
					return
				}
			}
		}
	}()
	return out
}
