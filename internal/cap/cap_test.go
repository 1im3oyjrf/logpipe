package cap_test

import (
	"context"
	"testing"
	"time"

	capper "github.com/logpipe/logpipe/internal/cap"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(msg string) reader.Entry {
	return reader.Entry{Fields: map[string]interface{}{"message": msg}}
}

func feed(entries []reader.Entry) <-chan reader.Entry {
	ch := make(chan reader.Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func collect(ch <-chan reader.Entry) []reader.Entry {
	var out []reader.Entry
	for e := range ch {
		out = append(out, e)
	}
	return out
}

func TestRun_ZeroMax_ForwardsAll(t *testing.T) {
	entries := []reader.Entry{makeEntry("a"), makeEntry("b"), makeEntry("c")}
	c := capper.New(capper.Config{Max: 0})
	result := collect(c.Run(context.Background(), feed(entries)))
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestRun_MaxOne_ForwardsOne(t *testing.T) {
	entries := []reader.Entry{makeEntry("a"), makeEntry("b"), makeEntry("c")}
	c := capper.New(capper.Config{Max: 1})
	result := collect(c.Run(context.Background(), feed(entries)))
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
}

func TestRun_MaxExceedsInput_ForwardsAll(t *testing.T) {
	entries := []reader.Entry{makeEntry("a"), makeEntry("b")}
	c := capper.New(capper.Config{Max: 10})
	result := collect(c.Run(context.Background(), feed(entries)))
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestRun_ContextCancellation_Stops(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	blocking := make(chan reader.Entry)
	c := capper.New(capper.Config{Max: 0})
	result := collect(c.Run(ctx, blocking))
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}
