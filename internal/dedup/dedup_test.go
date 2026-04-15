package dedup_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/dedup"
	"github.com/yourorg/logpipe/internal/reader"
)

func entry(level, msg string) reader.Entry {
	return reader.Entry{Level: level, Message: msg}
}

func TestIsDuplicate_FirstOccurrence_NotDuplicate(t *testing.T) {
	f := dedup.New(5 * time.Second)
	if f.IsDuplicate(entry("info", "hello world")) {
		t.Fatal("expected first occurrence to not be a duplicate")
	}
}

func TestIsDuplicate_SecondOccurrence_IsDuplicate(t *testing.T) {
	f := dedup.New(5 * time.Second)
	f.IsDuplicate(entry("info", "hello world"))
	if !f.IsDuplicate(entry("info", "hello world")) {
		t.Fatal("expected second occurrence to be a duplicate")
	}
}

func TestIsDuplicate_DifferentLevel_NotDuplicate(t *testing.T) {
	f := dedup.New(5 * time.Second)
	f.IsDuplicate(entry("info", "hello world"))
	if f.IsDuplicate(entry("error", "hello world")) {
		t.Fatal("expected different level to not be a duplicate")
	}
}

func TestIsDuplicate_DifferentMessage_NotDuplicate(t *testing.T) {
	f := dedup.New(5 * time.Second)
	f.IsDuplicate(entry("info", "hello world"))
	if f.IsDuplicate(entry("info", "goodbye world")) {
		t.Fatal("expected different message to not be a duplicate")
	}
}

func TestIsDuplicate_AfterTTLExpiry_NotDuplicate(t *testing.T) {
	now := time.Now()
	f := dedup.New(1 * time.Second)
	// Inject controlled clock.
	f = dedup.New(1 * time.Second)

	// Manually test TTL by using a very short window and sleeping.
	f2 := &dedupClock{Filter: dedup.New(500 * time.Millisecond), offset: 0}
	f2.Filter.IsDuplicate(entry("info", "expiring message"))
	_ = now

	time.Sleep(600 * time.Millisecond)
	if f2.Filter.IsDuplicate(entry("info", "expiring message")) {
		t.Fatal("expected entry to not be duplicate after TTL expiry")
	}
}

func TestSize_TracksEntries(t *testing.T) {
	f := dedup.New(5 * time.Second)
	if f.Size() != 0 {
		t.Fatalf("expected size 0, got %d", f.Size())
	}
	f.IsDuplicate(entry("info", "msg1"))
	f.IsDuplicate(entry("info", "msg2"))
	f.IsDuplicate(entry("info", "msg1")) // duplicate, should not grow
	if f.Size() != 2 {
		t.Fatalf("expected size 2, got %d", f.Size())
	}
}

// dedupClock is a thin wrapper to allow indirect TTL testing via real sleep.
type dedupClock struct {
	*dedup.Filter
	offset time.Duration
}
