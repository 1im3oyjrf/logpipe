package pipeline

import (
	"context"
	"io"

	"github.com/your-org/logpipe/internal/dedup"
	"github.com/your-org/logpipe/internal/filter"
	"github.com/your-org/logpipe/internal/metrics"
	"github.com/your-org/logpipe/internal/output"
	"github.com/your-org/logpipe/internal/sampler"
	"github.com/your-org/logpipe/internal/source"
)

// Config holds all tunables for a pipeline run.
type Config struct {
	Source      *source.Multiplexer
	Filter      *filter.Filter
	Formatter   *output.Formatter
	Writer      io.Writer
	Metrics     *metrics.Metrics
	Dedup       *dedup.Deduplicator
	Sampler     *sampler.Sampler
	SampleRate  uint64 // 0/1 = disabled
}

// Run reads entries from cfg.Source, applies dedup/sampling/filtering,
// formats matching entries and writes them to cfg.Writer.
// It blocks until ctx is cancelled or the source is exhausted.
func Run(ctx context.Context, cfg Config) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		se, ok := cfg.Source.Next(ctx)
		if !ok {
			return nil
		}

		if cfg.Metrics != nil {
			cfg.Metrics.IncRead()
		}

		if cfg.Dedup != nil && cfg.Dedup.IsDuplicate(se.Entry) {
			if cfg.Metrics != nil {
				cfg.Metrics.IncDropped()
			}
			continue
		}

		if cfg.Sampler != nil && !cfg.Sampler.Keep(se.Entry) {
			if cfg.Metrics != nil {
				cfg.Metrics.IncDropped()
			}
			continue
		}

		if cfg.Filter != nil && !cfg.Filter.Match(se.Entry) {
			continue
		}

		if cfg.Metrics != nil {
			cfg.Metrics.IncMatched()
		}

		line := cfg.Formatter.Format(se.Source, se.Entry)
		_, _ = io.WriteString(cfg.Writer, line+"\n")
	}
}
