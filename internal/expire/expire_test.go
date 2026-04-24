package expire_test

import (
	"testing"
	"time"

	"logpipe/internal/expire"
)

func fixed(t time.Time) func() time.Time { return func() time.Time { return t } }

var epoch = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestAllow_RecentEntry_IsKept(t *testing.T) {
	e, _ := expire.New(expire.Config{MaxAge: time.Minute, Now: fixed(epoch)})
	entry := map[string]any{
		"timestamp": epoch.Add(-30 * time.Second).Format(time.RFC3339),
	}
	if !e.Allow(entry) {
		t.Fatal("expected recent entry to be kept")
	}
}

func TestAllow_OldEntry_IsDropped(t *testing.T) {
	e, _ := expire.New(expire.Config{MaxAge: time.Minute, Now: fixed(epoch)})
	entry := map[string]any{
		"timestamp": epoch.Add(-2 * time.Minute).Format(time.RFC3339),
	}
	if e.Allow(entry) {
		t.Fatal("expected old entry to be dropped")
	}
}

func TestAllow_MissingTimestamp_IsKept(t *testing.T) {
	e, _ := expire.New(expire.Config{MaxAge: time.Minute, Now: fixed(epoch)})
	entry := map[string]any{"message": "no timestamp here"}
	if !e.Allow(entry) {
		t.Fatal("expected entry with missing timestamp to be kept")
	}
}

func TestAllow_UnparsableTimestamp_IsKept(t *testing.T) {
	e, _ := expire.New(expire.Config{MaxAge: time.Minute, Now: fixed(epoch)})
	entry := map[string]any{"timestamp": "not-a-date"}
	if !e.Allow(entry) {
		t.Fatal("expected entry with bad timestamp to be kept")
	}
}

func TestAllow_RFC3339Nano_IsAccepted(t *testing.T) {
	e, _ := expire.New(expire.Config{MaxAge: time.Minute, Now: fixed(epoch)})
	entry := map[string]any{
		"timestamp": epoch.Add(-10 * time.Second).Format(time.RFC3339Nano),
	}
	if !e.Allow(entry) {
		t.Fatal("expected RFC3339Nano timestamp to be accepted")
	}
}

func TestAllow_CustomTimestampField(t *testing.T) {
	e, _ := expire.New(expire.Config{
		TimestampField: "ts",
		MaxAge:         time.Minute,
		Now:            fixed(epoch),
	})
	entry := map[string]any{
		"ts": epoch.Add(-2 * time.Minute).Format(time.RFC3339),
	}
	if e.Allow(entry) {
		t.Fatal("expected entry to be dropped via custom field")
	}
}

func TestNew_Defaults(t *testing.T) {
	e, err := expire.New(expire.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Field() != "timestamp" {
		t.Errorf("expected default field 'timestamp', got %q", e.Field())
	}
	if e.MaxAge() != 5*time.Minute {
		t.Errorf("expected default max age 5m, got %v", e.MaxAge())
	}
}

func TestAllow_ExactlyAtMaxAge_IsKept(t *testing.T) {
	e, _ := expire.New(expire.Config{MaxAge: time.Minute, Now: fixed(epoch)})
	entry := map[string]any{
		"timestamp": epoch.Add(-time.Minute).Format(time.RFC3339),
	}
	if !e.Allow(entry) {
		t.Fatal("expected entry exactly at max age boundary to be kept")
	}
}
