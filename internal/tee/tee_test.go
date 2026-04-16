package tee_test

import (
	"context"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/tee"
)

func makeEntry(msg string) reader.Entry {
	return reader.Entry{Message: msg, Level: "info", Fields: map[string]any{}}
}

func collect(ch <-chan reader.Entry) []reader.Entry {
	var out []reader.Entry
	for e := range ch {
		out = append(out, e)
	}
	return out
}

func TestTee_SingleConsumer_ReceivesAllEntries(t *testing.T) {
	tr := tee.New(8)
	ch := tr.Add(8)

	src := make(chan reader.Entry, 3)
	src <- makeEntry("a")
	src <- makeEntry("b")
	src <- makeEntry("c")
	close(src)

	tr.Run(context.Background(), src)

	got := collect(ch)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestTee_MultipleConsumers_EachReceiveAllEntries(t *testing.T) {
	tr := tee.New(8)
	ch1 := tr.Add(8)
	ch2 := tr.Add(8)

	src := make(chan reader.Entry, 2)
	src <- makeEntry("x")
	src <- makeEntry("y")
	close(src)

	tr.Run(context.Background(), src)

	for i, ch := range []<-chan reader.Entry{ch1, ch2} {
		got := collect(ch)
		if len(got) != 2 {
			t.Errorf("consumer %d: expected 2 entries, got %d", i, len(got))
		}
	}
}

func TestTee_ContextCancellation_StopsGracefully(t *testing.T) {
	tr := tee.New(8)
	_ = tr.Add(8)

	src := make(chan reader.Entry) // never sends
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		tr.Run(ctx, src)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not stop after context cancellation")
	}
}
