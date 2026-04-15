package metrics

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Reporter periodically writes metric snapshots to a writer.
type Reporter struct {
	metrics  *Metrics
	out      io.Writer
	interval time.Duration
}

// NewReporter creates a Reporter that writes to out every interval.
func NewReporter(m *Metrics, out io.Writer, interval time.Duration) *Reporter {
	return &Reporter{
		metrics:  m,
		out:      out,
		interval: interval,
	}
}

// Run starts the reporting loop, blocking until ctx is cancelled.
func (r *Reporter) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.write()
			return
		case <-ticker.C:
			r.write()
		}
	}
}

func (r *Reporter) write() {
	s := r.metrics.Snapshot()
	fmt.Fprintf(
		r.out,
		"[metrics] uptime=%.1fs read=%d matched=%d dropped=%d errors=%d match_rate=%.1f%%\n",
		s.Uptime.Seconds(),
		s.LinesRead,
		s.LinesMatched,
		s.LinesDropped,
		s.Errors,
		s.MatchRate()*100,
	)
}
