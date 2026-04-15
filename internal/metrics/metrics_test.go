package metrics_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/metrics"
)

func TestNew_InitialCountersAreZero(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()

	if s.LinesRead != 0 {
		t.Errorf("expected LinesRead=0, got %d", s.LinesRead)
	}
	if s.LinesMatched != 0 {
		t.Errorf("expected LinesMatched=0, got %d", s.LinesMatched)
	}
	if s.LinesDropped != 0 {
		t.Errorf("expected LinesDropped=0, got %d", s.LinesDropped)
	}
	if s.ParseErrors != 0 {
		t.Errorf("expected ParseErrors=0, got %d", s.ParseErrors)
	}
}

func TestCounters_Increments(t *testing.T) {
	c := metrics.New()

	for i := 0; i < 5; i++ {
		c.IncRead()
	}
	for i := 0; i < 3; i++ {
		c.IncMatched()
	}
	for i := 0; i < 2; i++ {
		c.IncDropped()
	}
	c.IncParseError()

	s := c.Snapshot()

	if s.LinesRead != 5 {
		t.Errorf("expected LinesRead=5, got %d", s.LinesRead)
	}
	if s.LinesMatched != 3 {
		t.Errorf("expected LinesMatched=3, got %d", s.LinesMatched)
	}
	if s.LinesDropped != 2 {
		t.Errorf("expected LinesDropped=2, got %d", s.LinesDropped)
	}
	if s.ParseErrors != 1 {
		t.Errorf("expected ParseErrors=1, got %d", s.ParseErrors)
	}
}

func TestSnapshot_UptimeIsPositive(t *testing.T) {
	c := metrics.New()
	time.Sleep(5 * time.Millisecond)
	s := c.Snapshot()

	if s.Uptime <= 0 {
		t.Errorf("expected positive uptime, got %v", s.Uptime)
	}
}

func TestSnapshot_IsImmutable(t *testing.T) {
	c := metrics.New()
	c.IncRead()
	s1 := c.Snapshot()

	c.IncRead()
	s2 := c.Snapshot()

	if s1.LinesRead == s2.LinesRead {
		t.Error("expected snapshots to differ after additional increment")
	}
}
