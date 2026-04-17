package window_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/window"
)

func TestGuard_WithinLimit_Allows(t *testing.T) {
	g := window.NewGuard(window.GuardConfig{WindowSize: time.Second, Limit: 3})
	for i := 0; i < 3; i++ {
		if !g.Allow() {
			t.Fatalf("expected Allow on call %d", i+1)
		}
	}
}

func TestGuard_ExceedsLimit_Blocks(t *testing.T) {
	g := window.NewGuard(window.GuardConfig{WindowSize: time.Second, Limit: 2})
	g.Allow()
	g.Allow()
	if g.Allow() {
		t.Fatal("expected Allow to return false when limit exceeded")
	}
}

func TestGuard_Reset_ClearsCount(t *testing.T) {
	g := window.NewGuard(window.GuardConfig{WindowSize: time.Second, Limit: 1})
	g.Allow()
	g.Allow() // over limit
	g.Reset()
	if !g.Allow() {
		t.Fatal("expected Allow after reset")
	}
}

func TestGuard_DefaultLimit_Applied(t *testing.T) {
	g := window.NewGuard(window.GuardConfig{WindowSize: time.Second})
	if g == nil {
		t.Fatal("expected non-nil guard")
	}
	if c := g.Count(); c != 0 {
		t.Fatalf("expected 0, got %d", c)
	}
}

func TestGuard_WindowExpiry_ResetsNaturally(t *testing.T) {
	g := window.NewGuard(window.GuardConfig{WindowSize: 50 * time.Millisecond, Limit: 1})
	g.Allow()
	time.Sleep(80 * time.Millisecond)
	if !g.Allow() {
		t.Fatal("expected Allow after window expiry")
	}
}
