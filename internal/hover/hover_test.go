package hover_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/hover"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(msg, level string) reader.Entry {
	return reader.Entry{
		"message": msg,
		"level":   level,
	}
}

func collect(out <-chan reader.Entry, n int) []reader.Entry {
	var entries []reader.Entry
	for i := 0; i < n; i++ {
		select {
		case e := <-out:
			entries = append(entries, e)
		case <-time.After(time.Second):
			return entries
		}
	}
	return entries
}

func TestHover_NilObserver_PassesThrough(t *testing.T) {
	h := hover.New(nil)
	in := make(chan reader.Entry, 2)
	out := make(chan reader.Entry, 2)

	in <- makeEntry("hello", "info")
	in <- makeEntry("world", "warn")
	close(in)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	h.Run(ctx, in, out)

	results := collect(out, 2)
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
}

func TestHover_ObserverReceivesAllEntries(t *testing.T) {
	var mu sync.Mutex
	var seen []reader.Entry

	obs := func(e reader.Entry) {
		mu.Lock()
		seen = append(seen, e)
		mu.Unlock()
	}

	h := hover.New(obs)
	in := make(chan reader.Entry, 3)
	out := make(chan reader.Entry, 3)

	in <- makeEntry("a", "info")
	in <- makeEntry("b", "error")
	in <- makeEntry("c", "debug")
	close(in)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	h.Run(ctx, in, out)

	mu.Lock()
	defer mu.Unlock()
	if len(seen) != 3 {
		t.Fatalf("observer expected 3 entries, got %d", len(seen))
	}
}

func TestHover_OutputMatchesInput(t *testing.T) {
	var observed reader.Entry
	obs := func(e reader.Entry) { observed = e }

	h := hover.New(obs)
	in := make(chan reader.Entry, 1)
	out := make(chan reader.Entry, 1)

	want := makeEntry("test-msg", "info")
	in <- want
	close(in)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	h.Run(ctx, in, out)

	got := <-out
	if got["message"] != want["message"] {
		t.Errorf("output message = %v, want %v", got["message"], want["message"])
	}
	if observed["message"] != want["message"] {
		t.Errorf("observed message = %v, want %v", observed["message"], want["message"])
	}
}

func TestHover_ContextCancellation_Stops(t *testing.T) {
	h := hover.New(nil)
	in := make(chan reader.Entry) // never sends
	out := make(chan reader.Entry, 1)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		h.Run(ctx, in, out)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
}
