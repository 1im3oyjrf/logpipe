package replay_test

import (
	"context"
	"testing"
	"time"

	"github.com/logpipe/internal/buffer"
	"github.com/logpipe/internal/reader"
	"github.com/logpipe/internal/replay"
)

func makeEntry(msg, level string) reader.Entry {
	return reader.Entry{
		Message: msg,
		Level:   level,
		Fields:  map[string]any{},
	}
}

func collect(ch <-chan reader.Entry) []reader.Entry {
	var out []reader.Entry
	for e := range ch {
		out = append(out, e)
	}
	return out
}

func TestReplay_AllEntries_NoFilter(t *testing.T) {
	buf := buffer.New(10)
	buf.Push(makeEntry("hello", "info"))
	buf.Push(makeEntry("world", "warn"))
	r := replay.New(buf, replay.Config{})
	entries := collect(r.Replay(context.Background()))
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestReplay_WithPattern_FiltersEntries(t *testing.T) {
	buf := buffer.New(10)
	buf.Push(makeEntry("hello logpipe", "info"))
	buf.Push(makeEntry("unrelated", "info"))
	r := replay.New(buf, replay.Config{Pattern: "logpipe"})
	entries := collect(r.Replay(context.Background()))
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Message != "hello logpipe" {
		t.Errorf("unexpected message: %s", entries[0].Message)
	}
}

func TestReplay_EmptyBuffer_ClosesImmediately(t *testing.T) {
	buf := buffer.New(10)
	r := replay.New(buf, replay.Config{})
	entries := collect(r.Replay(context.Background()))
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestReplay_ContextCancellation_StopsEarly(t *testing.T) {
	buf := buffer.New(100)
	for i := 0; i < 100; i++ {
		buf.Push(makeEntry("msg", "info"))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	r := replay.New(buf, replay.Config{})
	// Should not block indefinitely.
	_ = collect(r.Replay(ctx))
}
