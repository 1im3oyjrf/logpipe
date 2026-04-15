// Package source provides abstractions for reading log entries
// from multiple named input sources (files, stdin, etc.).
package source

import (
	"context"
	"io"
	"sync"

	"github.com/user/logpipe/internal/reader"
)

// Entry wraps a log entry with metadata about which source it came from.
type Entry struct {
	Source string
	Fields map[string]interface{}
}

// Source represents a named log input.
type Source struct {
	Name   string
	reader *reader.JSONReader
}

// New creates a new Source with the given name and underlying reader.
func New(name string, r io.Reader) *Source {
	return &Source{
		Name:   name,
		reader: reader.NewJSONReader(r),
	}
}

// Multiplexer fans in entries from multiple sources into a single channel.
type Multiplexer struct {
	sources []*Source
}

// NewMultiplexer creates a Multiplexer from the provided sources.
func NewMultiplexer(sources ...*Source) *Multiplexer {
	return &Multiplexer{sources: sources}
}

// Stream starts reading all sources concurrently and sends tagged entries
// to the returned channel. The channel is closed when all sources are exhausted
// or ctx is cancelled.
func (m *Multiplexer) Stream(ctx context.Context) <-chan Entry {
	out := make(chan Entry, 64)
	var wg sync.WaitGroup

	for _, s := range m.sources {
		wg.Add(1)
		go func(src *Source) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				fields, err := src.reader.Next()
				if err != nil {
					return
				}
				out <- Entry{Source: src.Name, Fields: fields}
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
