package metrics

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestReporter_WritesSnapshot(t *testing.T) {
	m := New()
	m.IncRead()
	m.IncRead()
	m.IncMatched()
	m.IncDropped()

	var buf bytes.Buffer
	r := NewReporter(m, &buf, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	r.Run(ctx)

	out := buf.String()
	if !strings.Contains(out, "[metrics]") {
		t.Errorf("expected [metrics] prefix, got: %q", out)
	}
	if !strings.Contains(out, "read=2") {
		t.Errorf("expected read=2 in output, got: %q", out)
	}
	if !strings.Contains(out, "matched=1") {
		t.Errorf("expected matched=1 in output, got: %q", out)
	}
}

func TestReporter_WritesOnCancel(t *testing.T) {
	m := New()
	m.IncRead()

	var buf bytes.Buffer
	r := NewReporter(m, &buf, 10*time.Second) // long interval — only cancel write

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	r.Run(ctx)

	if !strings.Contains(buf.String(), "read=1") {
		t.Errorf("expected final write on cancel, got: %q", buf.String())
	}
}

func TestSnapshot_MatchRate(t *testing.T) {
	s := Snapshot{LinesRead: 10, LinesMatched: 4}
	if got := s.MatchRate(); got != 0.4 {
		t.Errorf("expected 0.4, got %f", got)
	}
}

func TestSnapshot_MatchRate_ZeroRead(t *testing.T) {
	s := Snapshot{LinesRead: 0, LinesMatched: 0}
	if got := s.MatchRate(); got != 0 {
		t.Errorf("expected 0, got %f", got)
	}
}

func TestSnapshot_DropRate(t *testing.T) {
	s := Snapshot{LinesRead: 8, LinesDropped: 2}
	if got := s.DropRate(); got != 0.25 {
		t.Errorf("expected 0.25, got %f", got)
	}
}
