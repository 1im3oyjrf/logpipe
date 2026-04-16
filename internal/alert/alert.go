// Package alert provides threshold-based alerting for log entry metrics.
package alert

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Rule defines a condition that triggers an alert.
type Rule struct {
	Name      string
	Field     string
	Value     string
	Level     Level
	Threshold int
	Window    time.Duration
}

// Alert is a fired alert event.
type Alert struct {
	Rule      Rule
	Count     int
	FiredAt   time.Time
}

func (a Alert) String() string {
	return fmt.Sprintf("[ALERT:%s] %s — %d occurrences in %s (at %s)",
		a.Rule.Level, a.Rule.Name, a.Count, a.Rule.Window, a.FiredAt.Format(time.RFC3339))
}

// Evaluator tracks counts per rule and fires alerts to a writer.
type Evaluator struct {
	mu      sync.Mutex
	rules   []Rule
	buckets map[string][]time.Time
	out     io.Writer
	now     func() time.Time
}

// New creates a new Evaluator with the given rules and output writer.
func New(rules []Rule, out io.Writer) *Evaluator {
	return &Evaluator{
		rules:   rules,
		buckets: make(map[string][]time.Time),
		out:     out,
		now:     time.Now,
	}
}

// Observe checks the entry fields against all rules and fires alerts when thresholds are exceeded.
func (e *Evaluator) Observe(fields map[string]string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := e.now()
	for _, rule := range e.rules {
		v, ok := fields[rule.Field]
		if !ok || v != rule.Value {
			continue
		}
		key := rule.Name
		cutoff := now.Add(-rule.Window)
		times := e.buckets[key]
		// evict old entries
		filtered := times[:0]
		for _, t := range times {
			if t.After(cutoff) {
				filtered = append(filtered, t)
			}
		}
		filtered = append(filtered, now)
		e.buckets[key] = filtered

		if len(filtered) == rule.Threshold {
			a := Alert{Rule: rule, Count: len(filtered), FiredAt: now}
			fmt.Fprintln(e.out, a.String())
		}
	}
}
