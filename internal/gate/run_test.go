package gate_test

import (
	"context"
	"testing"

	"github.com/logpipe/internal/gate"
)

func feed(entries []map[string]any) <-chan map[string]any {
	ch := make(chan map[string]any, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func collect(ch <-chan map[string]any) []map[string]any {
	var out []map[string]any
	for e := range ch {
		out = append(out, e)
	}
	return out
}

func TestRun_PassingEntries_Forwarded(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "level", Op: gate.OpEq, Value: "error"})
	in := feed([]map[string]any{
		{"level": "error", "msg": "boom"},
		{"level": "info", "msg": "ok"},
		{"level": "error", "msg": "again"},
	})
	out := make(chan map[string]any, 10)
	g.Run(context.Background(), in, out)
	results := collect(out)
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
}

func TestRun_NoMatch_EmitsNothing(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "level", Op: gate.OpEq, Value: "debug"})
	in := feed([]map[string]any{
		{"level": "info"},
		{"level": "error"},
	})
	out := make(chan map[string]any, 10)
	g.Run(context.Background(), in, out)
	results := collect(out)
	if len(results) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(results))
	}
}

func TestRun_ContextCancellation_Stops(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "level", Op: gate.OpEq, Value: "error"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	in := make(chan map[string]any)
	out := make(chan map[string]any, 1)
	g.Run(ctx, in, out)
	if _, open := <-out; open {
		t.Fatal("expected out to be closed")
	}
}
