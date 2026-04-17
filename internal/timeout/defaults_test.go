package timeout_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/timeout"
)

func TestNew_ZeroDuration_UsesDefault(t *testing.T) {
	g := timeout.New(timeout.Config{Duration: 0})
	d := timeout.DurationOf(g)
	if d.String() != (200 * time.Millisecond).String() {
		t.Fatalf("expected default 200ms, got %s", d.String())
	}
}

func TestNew_ExplicitDuration_Respected(t *testing.T) {
	g := timeout.New(timeout.Config{Duration: 500 * time.Millisecond})
	d := timeout.DurationOf(g)
	if d.String() != (500 * time.Millisecond).String() {
		t.Fatalf("expected 500ms, got %s", d.String())
	}
}

func TestNew_NegativeDuration_UsesDefault(t *testing.T) {
	g := timeout.New(timeout.Config{Duration: -1 * time.Second})
	d := timeout.DurationOf(g)
	if d.String() != (200 * time.Millisecond).String() {
		t.Fatalf("expected default 200ms for negative input, got %s", d.String())
	}
}
