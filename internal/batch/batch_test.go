package batch_test

import (
	"context"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/batch"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(msg string) reader.Entry {
	return reader.Entry{Message: msg, Level: "info", Fields: map[string]any{}}
}

func feed(entries []reader.Entry) <-chan reader.Entry {
	ch := make(chan reader.Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func TestRun_FlushesOnMaxSize(t *testing.T) {
	entries := []reader.Entry{makeEntry("a"), makeEntry("b"), makeEntry("c")}
	ch := feed(entries)
	b := batch.New(batch.Config{MaxSize: 2, MaxWait: time.Second}, ch)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out := b.Run(ctx)
	var got [][]reader.Entry
	for batch := range out {
		got = append(got, batch)
	}
	if len(got) == 0 {
		t.Fatal("expected at least one batch")
	}
	total := 0
	for _, g := range got {
		total += len(g)
	}
	if total != 3 {
		t.Fatalf("expected 3 entries total, got %d", total)
	}
	if len(got[0]) != 2 {
		t.Fatalf("first batch should have 2 entries, got %d", len(got[0]))
	}
}

func TestRun_FlushesOnTimer(t *testing.T) {
	ch := make(chan reader.Entry, 1)
	ch <- makeEntry("only")
	b := batch.New(batch.Config{MaxSize: 10, MaxWait: 50 * time.Millisecond}, ch)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	out := b.Run(ctx)
	time.AfterFunc(120*time.Millisecond, func() { close(ch) })
	var got [][]reader.Entry
	for batch := range out {
		got = append(got, batch)
	}
	if len(got) == 0 {
		t.Fatal("expected a timer-triggered batch")
	}
}

func TestRun_EmptyInput_NoOutput(t *testing.T) {
	ch := make(chan reader.Entry)
	close(ch)
	b := batch.New(batch.Config{MaxSize: 5, MaxWait: 50 * time.Millisecond}, ch)
	ctx := context.Background()
	out := b.Run(ctx)
	var got [][]reader.Entry
	for batch := range out {
		got = append(got, batch)
	}
	if len(got) != 0 {
		t.Fatalf("expected no batches, got %d", len(got))
	}
}

func TestRun_ContextCancel_FlushesRemainder(t *testing.T) {
	ch := make(chan reader.Entry, 10)
	for i := 0; i < 3; i++ {
		ch <- makeEntry("x")
	}
	b := batch.New(batch.Config{MaxSize: 100, MaxWait: 10 * time.Second}, ch)
	ctx, cancel := context.WithCancel(context.Background())
	out := b.Run(ctx)
	time.AfterFunc(30*time.Millisecond, cancel)
	var total int
	for batch := range out {
		total += len(batch)
	}
	if total != 3 {
		t.Fatalf("expected 3 flushed on cancel, got %d", total)
	}
}
