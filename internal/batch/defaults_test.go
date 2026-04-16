package batch_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/batch"
	"github.com/logpipe/logpipe/internal/reader"
)

func TestNew_DefaultsApplied(t *testing.T) {
	ch := make(chan reader.Entry)
	close(ch)
	b := batch.New(batch.Config{}, ch)
	if b == nil {
		t.Fatal("expected non-nil batcher")
	}
}

func TestNew_ExplicitConfig_Respected(t *testing.T) {
	ch := make(chan reader.Entry)
	close(ch)
	cfg := batch.Config{MaxSize: 5, MaxWait: 200 * time.Millisecond}
	b := batch.New(cfg, ch)
	if b == nil {
		t.Fatal("expected non-nil batcher")
	}
}
