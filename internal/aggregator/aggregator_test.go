package aggregator

import (
	"testing"
	"time"

	"github.com/user/logpipe/internal/reader"
)

func makeEntry(level, fieldVal, fieldKey string) reader.Entry {
	fields := map[string]string{}
	if fieldKey != "" {
		fields[fieldKey] = fieldVal
	}
	return reader.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   "test message",
		Fields:    fields,
	}
}

func TestAdd_SingleEntry_CreatesBucket(t *testing.T) {
	a := New("level", 10*time.Second)
	a.Add(makeEntry("info", "info", "level"))

	snap := a.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap))
	}
	if snap[0].Count != 1 {
		t.Errorf("expected count 1, got %d", snap[0].Count)
	}
	if snap[0].Key != "info" {
		t.Errorf("expected key 'info', got %q", snap[0].Key)
	}
}

func TestAdd_MultipleEntries_SameKey_IncrementsCount(t *testing.T) {
	a := New("level", 10*time.Second)
	for i := 0; i < 5; i++ {
		a.Add(makeEntry("error", "error", "level"))
	}

	snap := a.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap))
	}
	if snap[0].Count != 5 {
		t.Errorf("expected count 5, got %d", snap[0].Count)
	}
}

func TestAdd_DifferentKeys_SeparateBuckets(t *testing.T) {
	a := New("level", 10*time.Second)
	a.Add(makeEntry("info", "info", "level"))
	a.Add(makeEntry("error", "error", "level"))
	a.Add(makeEntry("warn", "warn", "level"))

	snap := a.Snapshot()
	if len(snap) != 3 {
		t.Errorf("expected 3 buckets, got %d", len(snap))
	}
}

func TestAdd_MissingField_UsesUnknownKey(t *testing.T) {
	a := New("service", 10*time.Second)
	a.Add(makeEntry("info", "", ""))

	snap := a.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap))
	}
	if snap[0].Key != "(unknown)" {
		t.Errorf("expected key '(unknown)', got %q", snap[0].Key)
	}
}

func TestAdd_ExpiredWindow_ResetsBucket(t *testing.T) {
	a := New("level", 50*time.Millisecond)
	fixedNow := time.Now()
	a.now = func() time.Time { return fixedNow }
	a.Add(makeEntry("info", "info", "level"))

	// Advance time beyond window
	a.now = func() time.Time { return fixedNow.Add(100 * time.Millisecond) }
	a.Add(makeEntry("info", "info", "level"))

	snap := a.Snapshot()
	if snap[0].Count != 1 {
		t.Errorf("expected bucket reset to count 1, got %d", snap[0].Count)
	}
}

func TestReset_ClearsAllBuckets(t *testing.T) {
	a := New("level", 10*time.Second)
	a.Add(makeEntry("info", "info", "level"))
	a.Add(makeEntry("error", "error", "level"))
	a.Reset()

	snap := a.Snapshot()
	if len(snap) != 0 {
		t.Errorf("expected 0 buckets after reset, got %d", len(snap))
	}
}

func TestNew_DefaultField_IsLevel(t *testing.T) {
	a := New("", 10*time.Second)
	if a.field != "level" {
		t.Errorf("expected default field 'level', got %q", a.field)
	}
}

func TestNew_ZeroWindow_UsesDefault(t *testing.T) {
	a := New("level", 0)
	if a.window != 10*time.Second {
		t.Errorf("expected default window 10s, got %v", a.window)
	}
}
