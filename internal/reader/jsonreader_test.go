package reader

import (
	"strings"
	"testing"
	"time"
)

func TestJSONReader_ValidEntries(t *testing.T) {
	input := `{"level":"info","msg":"started","time":"2024-01-15T10:00:00Z"}
{"level":"error","msg":"failed","code":500}
`
	r := NewJSONReader("test-source", strings.NewReader(input))
	out := make(chan LogEntry, 10)

	if err := r.ReadAll(out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	close(out)

	entries := collectEntries(out)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Source != "test-source" {
		t.Errorf("expected source %q, got %q", "test-source", entries[0].Source)
	}
	if entries[0].Fields["level"] != "info" {
		t.Errorf("expected level=info, got %v", entries[0].Fields["level"])
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2024-01-15T10:00:00Z")
	if !entries[0].Timestamp.Equal(expectedTime) {
		t.Errorf("expected timestamp %v, got %v", expectedTime, entries[0].Timestamp)
	}

	// Second entry has no timestamp — should be zero value.
	if !entries[1].Timestamp.IsZero() {
		t.Errorf("expected zero timestamp for entry without time field")
	}
}

func TestJSONReader_SkipsInvalidLines(t *testing.T) {
	input := `not json at all
{"level":"warn","msg":"ok"}
`
	r := NewJSONReader("src", strings.NewReader(input))
	out := make(chan LogEntry, 10)

	if err := r.ReadAll(out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	close(out)

	entries := collectEntries(out)
	if len(entries) != 1 {
		t.Fatalf("expected 1 valid entry, got %d", len(entries))
	}
	if entries[0].Fields["msg"] != "ok" {
		t.Errorf("unexpected entry: %v", entries[0].Fields)
	}
}

func TestJSONReader_EmptyInput(t *testing.T) {
	r := NewJSONReader("empty", strings.NewReader(""))
	out := make(chan LogEntry, 10)

	if err := r.ReadAll(out); err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}
	close(out)

	if len(collectEntries(out)) != 0 {
		t.Error("expected no entries for empty input")
	}
}

func collectEntries(ch <-chan LogEntry) []LogEntry {
	var entries []LogEntry
	for e := range ch {
		entries = append(entries, e)
	}
	return entries
}
