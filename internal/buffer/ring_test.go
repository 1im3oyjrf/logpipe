package buffer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/buffer"
	"github.com/yourorg/logpipe/internal/reader"
)

func makeEntry(msg string) reader.Entry {
	return reader.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    map[string]interface{}{},
	}
}

func TestNew_PanicOnZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for zero capacity")
		}
	}()
	buffer.New(0)
}

func TestPush_BelowCapacity(t *testing.T) {
	b := buffer.New(5)
	b.Push(makeEntry("a"))
	b.Push(makeEntry("b"))
	if got := b.Len(); got != 2 {
		t.Fatalf("expected len 2, got %d", got)
	}
}

func TestSnapshot_ChronologicalOrder(t *testing.T) {
	b := buffer.New(4)
	for i := 0; i < 4; i++ {
		b.Push(makeEntry(fmt.Sprintf("msg%d", i)))
	}
	snap := b.Snapshot()
	for i, e := range snap {
		want := fmt.Sprintf("msg%d", i)
		if e.Message != want {
			t.Errorf("index %d: want %q, got %q", i, want, e.Message)
		}
	}
}

func TestPush_OverwritesOldest(t *testing.T) {
	b := buffer.New(3)
	for i := 0; i < 5; i++ {
		b.Push(makeEntry(fmt.Sprintf("msg%d", i)))
	}
	if got := b.Len(); got != 3 {
		t.Fatalf("expected len 3 after overflow, got %d", got)
	}
	snap := b.Snapshot()
	expected := []string{"msg2", "msg3", "msg4"}
	for i, e := range snap {
		if e.Message != expected[i] {
			t.Errorf("index %d: want %q, got %q", i, expected[i], e.Message)
		}
	}
}

func TestSnapshot_EmptyBuffer(t *testing.T) {
	b := buffer.New(10)
	if snap := b.Snapshot(); snap != nil {
		t.Errorf("expected nil snapshot for empty buffer, got %v", snap)
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	b := buffer.New(4)
	b.Push(makeEntry("x"))
	b.Push(makeEntry("y"))
	b.Reset()
	if b.Len() != 0 {
		t.Errorf("expected len 0 after reset, got %d", b.Len())
	}
	if snap := b.Snapshot(); snap != nil {
		t.Errorf("expected nil snapshot after reset")
	}
}

func TestCap_ReturnsCapacity(t *testing.T) {
	b := buffer.New(7)
	if b.Cap() != 7 {
		t.Errorf("expected cap 7, got %d", b.Cap())
	}
}
