package window_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/window"
)

func TestNew_DefaultSize(t *testing.T) {
	w := window.New(window.Config{})
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestAdd_IncreasesCount(t *testing.T) {
	w := window.New(window.Config{Size: time.Second})
	w.Add(3)
	w.Add(2)
	if got := w.Count(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestCount_EvictsExpiredBuckets(t *testing.T) {
	w := window.New(window.Config{Size: 50 * time.Millisecond})
	w.Add(10)
	time.Sleep(80 * time.Millisecond)
	w.Add(4)
	if got := w.Count(); got != 4 {
		t.Fatalf("expected 4 after eviction, got %d", got)
	}
}

func TestReset_ClearsAllBuckets(t *testing.T) {
	w := window.New(window.Config{Size: time.Second})
	w.Add(7)
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCount_EmptyWindow_ReturnsZero(t *testing.T) {
	w := window.New(window.Config{Size: time.Second})
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_NegativeSizeUsesDefault(t *testing.T) {
	w := window.New(window.Config{Size: -1})
	w.Add(1)
	if got := w.Count(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}
