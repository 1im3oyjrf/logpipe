package sieve_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/sieve"
)

func base(msg string) map[string]any {
	return map[string]any{"message": msg, "level": "info"}
}

func TestAllow_FirstOccurrence_IsAllowed(t *testing.T) {
	s := sieve.New(sieve.Config{})
	if !s.Allow(base("hello")) {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestAllow_SecondOccurrence_IsDropped(t *testing.T) {
	s := sieve.New(sieve.Config{})
	s.Allow(base("hello"))
	if s.Allow(base("hello")) {
		t.Fatal("expected second occurrence to be dropped")
	}
}

func TestAllow_DifferentValues_BothAllowed(t *testing.T) {
	s := sieve.New(sieve.Config{Slots: 4096})
	if !s.Allow(base("alpha")) {
		t.Fatal("expected alpha to be allowed")
	}
	if !s.Allow(base("beta")) {
		t.Fatal("expected beta to be allowed")
	}
}

func TestAllow_CaseInsensitive_TreatsAsSame(t *testing.T) {
	s := sieve.New(sieve.Config{CaseInsensitive: true})
	s.Allow(base("Hello"))
	if s.Allow(base("hello")) {
		t.Fatal("expected case-insensitive duplicate to be dropped")
	}
}

func TestAllow_CaseSensitive_TreatsAsDifferent(t *testing.T) {
	s := sieve.New(sieve.Config{CaseInsensitive: false, Slots: 4096})
	s.Allow(base("Hello"))
	if !s.Allow(base("hello")) {
		t.Fatal("expected case-sensitive entries to be treated as distinct")
	}
}

func TestReset_ClearsSlots(t *testing.T) {
	s := sieve.New(sieve.Config{})
	s.Allow(base("hello"))
	s.Reset()
	if !s.Allow(base("hello")) {
		t.Fatal("expected entry to be allowed after reset")
	}
}

func TestAllow_CustomField_UsedForHashing(t *testing.T) {
	s := sieve.New(sieve.Config{Field: "request_id"})
	e1 := map[string]any{"request_id": "abc-123", "level": "info"}
	e2 := map[string]any{"request_id": "abc-123", "level": "warn"}
	s.Allow(e1)
	if s.Allow(e2) {
		t.Fatal("expected same request_id to be dropped")
	}
}

func TestAllow_MissingField_EmptyStringHashed(t *testing.T) {
	s := sieve.New(sieve.Config{Field: "trace_id"})
	e := map[string]any{"message": "no trace"}
	if !s.Allow(e) {
		t.Fatal("expected missing field (empty string) to be allowed on first occurrence")
	}
	if s.Allow(e) {
		t.Fatal("expected missing field duplicate to be dropped")
	}
}

func TestNew_DefaultSlots_Applied(t *testing.T) {
	s := sieve.New(sieve.Config{})
	if s == nil {
		t.Fatal("expected non-nil sieve")
	}
}
