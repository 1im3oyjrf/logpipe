// Package pipeline wires together source multiplexing, filtering,
// and formatted output into a single cohesive processing loop.
package pipeline

import (
	"context"
	"io"

	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/output"
	"github.com/user/logpipe/internal/source"
)

// Config holds the runtime options for a pipeline run.
type Config struct {
	// Sources is the list of named log inputs to read from.
	Sources []*source.Source
	// Pattern is an optional grep pattern forwarded to the filter.
	Pattern string
	// CaseSensitive controls whether pattern matching is case-sensitive.
	CaseSensitive bool
	// FilterFields limits pattern matching to specific JSON field names.
	FilterFields []string
	// NoColor disables ANSI colour codes in output.
	NoColor bool
	// Out is the writer used for formatted output (defaults to os.Stdout).
	Out io.Writer
}

// Run starts the pipeline and blocks until all sources are drained or
// ctx is cancelled. It returns the number of entries written.
func Run(ctx context.Context, cfg Config) (int, error) {
	f := filter.New(filter.Options{
		Pattern:       cfg.Pattern,
		CaseSensitive: cfg.CaseSensitive,
		Fields:        cfg.FilterFields,
	})

	fmt := output.New(output.Options{
		NoColor: cfg.NoColor,
		Out:     cfg.Out,
	})

	mux := source.NewMultiplexer(cfg.Sources...)
	stream := mux.Stream(ctx)

	var written int
	for entry := range stream {
		if !f.Match(entry.Fields) {
			continue
		}
		if err := fmt.Write(entry.Source, entry.Fields); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}
