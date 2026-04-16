package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/limiter"
)

func TestNew_PanicOnZeroCap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for n=0")
		}
	}()
	limiter.New(0)
}

func TestAcquire_WithinCap_Succeeds(t *testing.T) {
	l := limiter.New(2)
	ctx := context.Background()
	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Available() != 1 {
		t.Fatalf("expected 1 available, got %d", l.Available())
	}
	l.Release()
	if l.Available() != 2 {
		t.Fatalf("expected 2 available after release, got %d", l.Available())
	}
}

func TestAcquire_CancelledContext_ReturnsError(t *testing.T) {
	l := limiter.New(1)
	ctx := context.Background()
	_ = l.Acquire(ctx) // fill the slot

	ctxC, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := l.Acquire(ctxC)
	if err != limiter.ErrLimitExceeded {
		t.Fatalf("expected ErrLimitExceeded, got %v", err)
	}
}

func TestAcquire_ConcurrentGoroutines_RespectsLimit(t *testing.T) {
	const cap = 3
	l := limiter.New(cap)
	ctx := context.Background()

	var mu sync.Mutex
	peak := 0
	current := 0
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Acquire(ctx); err != nil {
				return
			}
			defer l.Release()
			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			current--
			mu.Unlock()
		}()
	}
	wg.Wait()
	if peak > cap {
		t.Fatalf("peak concurrency %d exceeded cap %d", peak, cap)
	}
}

func TestCap_ReturnsConfiguredValue(t *testing.T) {
	l := limiter.New(5)
	if l.Cap() != 5 {
		t.Fatalf("expected cap 5, got %d", l.Cap())
	}
}
