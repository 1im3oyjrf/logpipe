package splitter_test

import (
	"context"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/splitter"
)

func makeEntry(msg string) map[string]any {
	return map[string]any{"message": msg, "level": "info"}
}

func collect(ch <-chan map[string]any, n int) []map[string]any {
	var out []map[string]any
	for i := 0; i < n; i++ {
		select {
		case e := <-ch:
			out = append(out, e)
		case <-time.After(time.Second):
			return out
		}
	}
	return out
}

func TestSplitter_SingleTarget_ReceivesAllEntries(t *testing.T) {
	s := splitter.New()
	ch := make(chan map[string]any, 4)
	s.Add("a", ch)

	in := make(chan map[string]any, 2)
	in <- makeEntry("hello")
	in <- makeEntry("world")
	close(in)

	s.Run(context.Background(), in)
	got := collect(ch, 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestSplitter_MultipleTargets_EachReceiveAll(t *testing.T) {
	s := splitter.New()
	ch1 := make(chan map[string]any, 4)
	ch2 := make(chan map[string]any, 4)
	s.Add("x", ch1)
	s.Add("y", ch2)

	in := make(chan map[string]any, 3)
	for i := 0; i < 3; i++ {
		in <- makeEntry("msg")
	}
	close(in)

	s.Run(context.Background(), in)

	if len(collect(ch1, 3)) != 3 {
		t.Error("ch1 did not receive all entries")
	}
	if len(collect(ch2, 3)) != 3 {
		t.Error("ch2 did not receive all entries")
	}
}

func TestSplitter_ContextCancellation_Stops(t *testing.T) {
	s := splitter.New()
	ch := make(chan map[string]any) // unbuffered — will block
	s.Add("slow", ch)

	in := make(chan map[string]any, 1)
	in <- makeEntry("blocked")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		s.Run(ctx, in)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
}

func TestSplitter_NoTargets_DropsEntries(t *testing.T) {
	s := splitter.New()
	in := make(chan map[string]any, 2)
	in <- makeEntry("a")
	in <- makeEntry("b")
	close(in)
	s.Run(context.Background(), in) // should complete without blocking
}
