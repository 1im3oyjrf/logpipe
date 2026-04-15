package source_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/logpipe/internal/source"
)

func TestMultiplexer_SingleSource(t *testing.T) {
	input := strings.NewReader(`{"level":"info","msg":"hello"}` + "\n")
	s := source.New("app", input)
	mux := source.NewMultiplexer(s)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var entries []source.Entry
	for e := range mux.Stream(ctx) {
		entries = append(entries, e)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Source != "app" {
		t.Errorf("expected source 'app', got %q", entries[0].Source)
	}
	if entries[0].Fields["msg"] != "hello" {
		t.Errorf("unexpected msg field: %v", entries[0].Fields["msg"])
	}
}

func TestMultiplexer_MultipleSources(t *testing.T) {
	a := source.New("svc-a", strings.NewReader("{\"msg\":\"from-a\"}\n"))
	b := source.New("svc-b", strings.NewReader("{\"msg\":\"from-b\"}\n"))
	mux := source.NewMultiplexer(a, b)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	seen := map[string]bool{}
	for e := range mux.Stream(ctx) {
		seen[e.Source] = true
	}

	if !seen["svc-a"] || !seen["svc-b"] {
		t.Errorf("expected entries from both sources, got: %v", seen)
	}
}

func TestMultiplexer_EmptySource(t *testing.T) {
	s := source.New("empty", strings.NewReader(""))
	mux := source.NewMultiplexer(s)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var count int
	for range mux.Stream(ctx) {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 entries from empty source, got %d", count)
	}
}

func TestMultiplexer_ContextCancellation(t *testing.T) {
	// Infinite-like source via a blocking reader is hard to simulate;
	// instead verify cancel stops a multi-line source mid-stream gracefully.
	lines := strings.Repeat("{\"msg\":\"x\"}\n", 1000)
	s := source.New("heavy", strings.NewReader(lines))
	mux := source.NewMultiplexer(s)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	// Should not block indefinitely.
	done := make(chan struct{})
	go func() {
		for range mux.Stream(ctx) {
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for cancelled stream to finish")
	}
}
