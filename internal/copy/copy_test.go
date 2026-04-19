package copy_test

import (
	"context"
	"testing"

	"logpipe/internal/copy"
)

func feed(entries []copy.Entry) <-chan copy.Entry {
	ch := make(chan copy.Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func collect(out <-chan copy.Entry) []copy.Entry {
	var result []copy.Entry
	for e := range out {
		result = append(result, e)
	}
	return result
}

func TestRun_EmitsOriginalAndCopy(t *testing.T) {
	c := copy.New(copy.Config{})
	in := feed([]copy.Entry{{"msg": "hello", "level": "info"}})
	out := make(chan copy.Entry, 10)
	c.Run(context.Background(), in, out)
	entries := collect(out)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0]["msg"] != "hello" || entries[1]["msg"] != "hello" {
		t.Fatal("both entries should carry the original message")
	}
}

func TestRun_CopyHasOverrides(t *testing.T) {
	c := copy.New(copy.Config{Overrides: map[string]string{"source": "copy"}})
	in := feed([]copy.Entry{{"msg": "hi", "source": "original"}})
	out := make(chan copy.Entry, 10)
	c.Run(context.Background(), in, out)
	entries := collect(out)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0]["source"] != "original" {
		t.Errorf("original should be unchanged, got %v", entries[0]["source"])
	}
	if entries[1]["source"] != "copy" {
		t.Errorf("copy should have override, got %v", entries[1]["source"])
	}
}

func TestRun_OriginalNotMutatedByOverride(t *testing.T) {
	c := copy.New(copy.Config{Overrides: map[string]string{"level": "debug"}})
	orig := copy.Entry{"msg": "test", "level": "info"}
	in := feed([]copy.Entry{orig})
	out := make(chan copy.Entry, 10)
	c.Run(context.Background(), in, out)
	entries := collect(out)
	if entries[0]["level"] != "info" {
		t.Errorf("original level should be info, got %v", entries[0]["level"])
	}
	if entries[1]["level"] != "debug" {
		t.Errorf("copy level should be debug, got %v", entries[1]["level"])
	}
}

func TestRun_ContextCancellation_Stops(t *testing.T) {
	c := copy.New(copy.Config{})
	blocking := make(chan copy.Entry)
	out := make(chan copy.Entry, 2)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Run(ctx, blocking, out)
	if _, ok := <-out; ok {
		t.Fatal("expected out to be closed")
	}
}
