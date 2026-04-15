package throttle_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/throttle"
)

func TestAllow_FirstOccurrence_IsAllowed(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	if !th.Allow("key1") {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestAllow_WithinCooldown_IsSuppressed(t *testing.T) {
	th := throttle.New(1 * time.Hour)
	th.Allow("key1")
	if th.Allow("key1") {
		t.Fatal("expected second occurrence within cooldown to be suppressed")
	}
}

func TestAllow_AfterCooldown_IsAllowed(t *testing.T) {
	now := time.Now()
	th := throttle.New(50 * time.Millisecond)

	// Inject a controllable clock.
	calls := 0
	th = newWithClock(50*time.Millisecond, func() time.Time {
		calls++
		if calls == 1 {
			return now
		}
		return now.Add(100 * time.Millisecond)
	})

	th.Allow("key1")
	if !th.Allow("key1") {
		t.Fatal("expected entry to be allowed after cooldown expired")
	}
}

func TestAllow_DifferentKeys_IndependentWindows(t *testing.T) {
	th := throttle.New(1 * time.Hour)
	th.Allow("keyA")
	if !th.Allow("keyB") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_ClearsAllKeys(t *testing.T) {
	th := throttle.New(1 * time.Hour)
	th.Allow("k1")
	th.Allow("k2")
	th.Reset()
	if th.Len() != 0 {
		t.Fatalf("expected 0 keys after reset, got %d", th.Len())
	}
	if !th.Allow("k1") {
		t.Fatal("expected k1 to be allowed after reset")
	}
}

func TestEvict_RemovesExpiredKeys(t *testing.T) {
	now := time.Now()
	clock := now
	th := newWithClock(50*time.Millisecond, func() time.Time { return clock })

	th.Allow("old")
	clock = now.Add(200 * time.Millisecond)
	th.Allow("fresh")
	th.Evict()

	if th.Len() != 1 {
		t.Fatalf("expected 1 key after evict, got %d", th.Len())
	}
}

// newWithClock is a test helper that injects a custom clock into Throttle.
func newWithClock(d time.Duration, fn func() time.Time) *throttle.Throttle {
	th := throttle.New(d)
	th.SetClock(fn)
	return th
}
