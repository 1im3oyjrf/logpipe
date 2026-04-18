package ceiling_test

import (
	"testing"
	"time"

	"github.com/your-org/logpipe/internal/ceiling"
)

func TestNew_DefaultWindow_Applied(t *testing.T) {
	c := ceiling.New(ceiling.Config{Max: 10})
	got := ceiling.WindowOf(c)
	if got != time.Second {
		t.Fatalf("expected default window=1s, got %v", got)
	}
}

func TestNew_ExplicitWindow_Respected(t *testing.T) {
	w := 5 * time.Second
	c := ceiling.New(ceiling.Config{Max: 10, Window: w})
	if ceiling.WindowOf(c) != w {
		t.Fatalf("expected window=%v, got %v", w, ceiling.WindowOf(c))
	}
}

func TestNew_MaxStored(t *testing.T) {
	c := ceiling.New(ceiling.Config{Max: 42})
	if ceiling.MaxOf(c) != 42 {
		t.Fatalf("expected max=42, got %d", ceiling.MaxOf(c))
	}
}
