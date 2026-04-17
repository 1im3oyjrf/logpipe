package debounce_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/logpipe/internal/debounce"
)

func makeEntry(msg string) map[string]any {
	return map[string]any{"message": msg, "level": "info"}
}

func feed(in chan<- map[string]any, entries ...map[string]any) {
	for _, e := range entries {
		in <- e
	}
}

func collect(ch <-chan map[string]any, timeout time.Duration) []map[string]any {
	var out []map[string]any
	deadline := time.After(timeout)
	for {
		select {
		case e, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, e)
		case <-deadline:
			return out
		}
	}
}

func TestDebounce_UniqueEntries_AllForwarded(t *testing.T) {
	d := debounce.New(debounce.Config{QuietPeriod: 30 * time.Millisecond})
	in := make(chan map[string]any, 4)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := d.Run(ctx, in)

	feed(in, makeEntry("alpha"), makeEntry("beta"))
	close(in)

	results := collect(out, 300*time.Millisecond)
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
}

func TestDebounce_DuplicateBurst_EmitsOnce(t *testing.T) {
	d := debounce.New(debounce.Config{QuietPeriod: 40 * time.Millisecond})
	in := make(chan map[string]any, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := d.Run(ctx, in)

	for i := 0; i < 5; i++ {
		in <- makeEntry("repeated")
	}
	close(in)

	results := collect(out, 400*time.Millisecond)
	if len(results) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(results))
	}
	if v, ok := results[0]["debounce_count"]; !ok || v.(int) != 5 {
		t.Fatalf("expected debounce_count=5, got %v", results[0]["debounce_count"])
	}
}

func TestDebounce_OriginalNotMutated(t *testing.T) {
	d := debounce.New(debounce.Config{QuietPeriod: 30 * time.Millisecond})
	in := make(chan map[string]any, 4)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := d.Run(ctx, in)

	orig := makeEntry("once")
	in <- orig
	in <- makeEntry("once")
	close(in)

	collect(out, 300*time.Millisecond)
	if _, mutated := orig["debounce_count"]; mutated {
		t.Fatal("original entry was mutated")
	}
}

func TestDebounce_ContextCancellation_Stops(t *testing.T) {
	d := debounce.New(debounce.Config{QuietPeriod: 50 * time.Millisecond})
	in := make(chan map[string]any)
	ctx, cancel := context.WithCancel(context.Background())
	out := d.Run(ctx, in)
	cancel()
	select {
	case <-out:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("output channel not closed after context cancellation")
	}
}
