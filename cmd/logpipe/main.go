package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/logpipe/internal/config"
	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/pipeline"
	"github.com/yourorg/logpipe/internal/source"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "logpipe: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := config.Parse(args)
	if err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var sources []*source.Source
	for _, s := range cfg.Sources {
		src, err := source.New(s)
		if err != nil {
			return fmt.Errorf("opening source %q: %w", s, err)
		}
		sources = append(sources, src)
	}

	mux := source.NewMultiplexer(sources...)

	f, err := filter.New(cfg.Pattern, cfg.CaseSensitive)
	if err != nil {
		return fmt.Errorf("building filter: %w", err)
	}

	fmt := output.New(os.Stdout, cfg.NoColor, cfg.Fields)

	return pipeline.Run(ctx, mux, f, fmt)
}
