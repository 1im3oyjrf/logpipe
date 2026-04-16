package alert_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/alert"
)

func newEval(rules []alert.Rule, buf *bytes.Buffer) *alert.Evaluator {
	return alert.New(rules, buf)
}

func TestObserve_BelowThreshold_NoAlert(t *testing.T) {
	var buf bytes.Buffer
	rule := alert.Rule{Name: "too-many-errors", Field: "level", Value: "error",
		Level: alert.LevelError, Threshold: 3, Window: time.Minute}
	ev := newEval([]alert.Rule{rule}, &buf)

	ev.Observe(map[string]string{"level": "error"})
	ev.Observe(map[string]string{"level": "error"})

	if buf.Len() != 0 {
		t.Fatalf("expected no alert, got: %s", buf.String())
	}
}

func TestObserve_AtThreshold_FiresAlert(t *testing.T) {
	var buf bytes.Buffer
	rule := alert.Rule{Name: "too-many-errors", Field: "level", Value: "error",
		Level: alert.LevelError, Threshold: 3, Window: time.Minute}
	ev := newEval([]alert.Rule{rule}, &buf)

	for i := 0; i < 3; i++ {
		ev.Observe(map[string]string{"level": "error"})
	}

	if !strings.Contains(buf.String(), "too-many-errors") {
		t.Fatalf("expected alert output, got: %s", buf.String())
	}
}

func TestObserve_NonMatchingField_NoAlert(t *testing.T) {
	var buf bytes.Buffer
	rule := alert.Rule{Name: "warn-alert", Field: "level", Value: "warn",
		Level: alert.LevelWarn, Threshold: 2, Window: time.Minute}
	ev := newEval([]alert.Rule{rule}, &buf)

	ev.Observe(map[string]string{"level": "info"})
	ev.Observe(map[string]string{"level": "info"})

	if buf.Len() != 0 {
		t.Fatalf("unexpected alert: %s", buf.String())
	}
}

func TestObserve_WindowExpiry_ResetsCount(t *testing.T) {
	var buf bytes.Buffer
	rule := alert.Rule{Name: "burst", Field: "level", Value: "error",
		Level: alert.LevelError, Threshold: 2, Window: 50 * time.Millisecond}
	ev := newEval([]alert.Rule{rule}, &buf)

	ev.Observe(map[string]string{"level": "error"})
	time.Sleep(60 * time.Millisecond)
	ev.Observe(map[string]string{"level": "error"})

	if buf.Len() != 0 {
		t.Fatalf("expected no alert after window expiry, got: %s", buf.String())
	}
}

func TestObserve_ConcurrentSafe(t *testing.T) {
	var buf syncBuffer
	rule := alert.Rule{Name: "concurrent", Field: "level", Value: "error",
		Level: alert.LevelError, Threshold: 100, Window: time.Minute}
	ev := alert.New([]alert.Rule{rule}, &buf)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ev.Observe(map[string]string{"level": "error"})
		}()
	}
	wg.Wait()
}

type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}
