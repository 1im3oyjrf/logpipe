// Package debounce provides a stage that suppresses rapid repeated log entries
// within a configurable quiet period, emitting only the first occurrence and
// a summary when the window closes.
package debounce

import (
	"context"
	"sync"
	"time"

	"github.com/user/logpipe/internal/parser"
)

const defaultQuiet = 200 * time.Millisecond

// Config holds tuning parameters for the debouncer.
type Config struct {
	// QuietPeriod is how long to wait after the last duplicate before flushing.
	QuietPeriod time.Duration
	// KeyField is the entry field used to group duplicates (default: "message").
	KeyField string
}

type state struct {
	first  map[string]any
	count  int
	timer  *time.Timer
}

// Debouncer suppresses bursts of identical log lines.
type Debouncer struct {
	cfg    Config
	mu     sync.Mutex
	groups map[string]*state
}

// New returns a Debouncer configured with cfg.
func New(cfg Config) *Debouncer {
	if cfg.QuietPeriod <= 0 {
		cfg.QuietPeriod = defaultQuiet
	}
	if cfg.KeyField == "" {
		cfg.KeyField = "message"
	}
	return &Debouncer{cfg: cfg, groups: make(map[string]*state)}
}

// Run reads from in, debounces entries and writes survivors to the returned channel.
func (d *Debouncer) Run(ctx context.Context, in <-chan map[string]any) <-chan map[string]any {
	out := make(chan map[string]any, 64)
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
				d.handle(ctx, entry, out)
			}
		}
	}()
	return out
}

func (d *Debouncer) handle(ctx context.Context, entry map[string]any, out chan<- map[string]any) {
	key := parser.GetString(entry, d.cfg.KeyField)
	d.mu.Lock()
	s, exists := d.groups[key]
	if !exists {
		s = &state{first: entry, count: 1}
		d.groups[key] = s
		s.timer = time.AfterFunc(d.cfg.QuietPeriod, func() {
			d.flush(ctx, key, out)
		})
		d.mu.Unlock()
		return
	}
	s.count++
	s.timer.Reset(d.cfg.QuietPeriod)
	d.mu.Unlock()
}

func (d *Debouncer) flush(ctx context.Context, key string, out chan<- map[string]any) {
	d.mu.Lock()
	s, ok := d.groups[key]
	if !ok {
		d.mu.Unlock()
		return
	}
	entry := s.first
	count := s.count
	delete(d.groups, key)
	d.mu.Unlock()

	if count > 1 {
		copy := make(map[string]any, len(entry)+1)
		for k, v := range entry {
			copy[k] = v
		}
		copy["debounce_count"] = count
		entry = copy
	}
	select {
	case <-ctx.Done():
	case out <- entry:
	}
}
