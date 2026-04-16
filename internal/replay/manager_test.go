package replay_test

import (
	"context"
	"testing"

	"github.com/logpipe/internal/replay"
)

func TestManager_RecordAndSize(t *testing.T) {
	m := replay.NewManager(10)
	if m.Size() != 0 {
		t.Fatalf("expected empty manager")
	}
	m.Record(makeEntry("a", "info"))
	m.Record(makeEntry("b", "warn"))
	if m.Size() != 2 {
		t.Fatalf("expected size 2, got %d", m.Size())
	}
}

func TestManager_Replay_ReturnsRecordedEntries(t *testing.T) {
	m := replay.NewManager(10)
	m.Record(makeEntry("first", "info"))
	m.Record(makeEntry("second", "error"))
	entries := collect(m.Replay(context.Background(), replay.Config{}))
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestManager_Replay_WithPattern(t *testing.T) {
	m := replay.NewManager(10)
	m.Record(makeEntry("matched line", "info"))
	m.Record(makeEntry("other", "info"))
	entries := collect(m.Replay(context.Background(), replay.Config{Pattern: "matched"}))
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestManager_Replay_CapacityEvictsOldest(t *testing.T) {
	m := replay.NewManager(2)
	m.Record(makeEntry("old", "info"))
	m.Record(makeEntry("newer", "info"))
	m.Record(makeEntry("newest", "info"))
	// capacity is 2, so only 2 entries retained
	if m.Size() != 2 {
		t.Fatalf("expected size 2, got %d", m.Size())
	}
	entries := collect(m.Replay(context.Background(), replay.Config{}))
	if entries[0].Message == "old" {
		t.Error("oldest entry should have been evicted")
	}
}
