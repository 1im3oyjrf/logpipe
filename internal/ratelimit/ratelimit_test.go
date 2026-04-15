package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestNew_InitialBurstIsAvailable(t *testing.T) {
	l := New(10, 5)
	if l.tokens != 5 {
		t.Fatalf("expected 5 initial tokens, got %v", l.tokens)
	}
}

func TestWait_ConsumesBurstImmediately(t *testing.T) {
	l := New(100, 3)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
	}
}

func TestWait_CancelledContext_ReturnsError(t *testing.T) {
	// Zero rate so no tokens are ever replenished.
	l := New(0, 0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := l.Wait(ctx); err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
}

func TestWait_ContextTimeout_ReturnsError(t *testing.T) {
	l := New(0, 0) // no tokens, no refill
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if time.Since(start) < 10*time.Millisecond {
		t.Error("Wait returned too quickly before timeout")
	}
}

func TestWait_RefillsTokensOverTime(t *testing.T) {
	l := New(1000, 1) // 1000 tokens/sec, burst 1

	// Drain the single burst token.
	ctx := context.Background()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("unexpected error draining burst: %v", err)
	}

	// Simulate time passing by advancing lastTick into the past.
	l.mu.Lock()
	l.lastTick = time.Now().Add(-10 * time.Millisecond) // +10 tokens at 1000/s
	l.mu.Unlock()

	if err := l.Wait(ctx); err != nil {
		t.Fatalf("expected token after refill, got error: %v", err)
	}
}

func TestTryConsume_CapsAtMax(t *testing.T) {
	l := New(1000, 4)

	// Wind back lastTick far enough to overflow max.
	l.mu.Lock()
	l.lastTick = time.Now().Add(-10 * time.Second)
	l.mu.Unlock()

	l.tryConsume()

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tokens > l.max {
		t.Errorf("tokens %v exceeded max %v", l.tokens, l.max)
	}
}
