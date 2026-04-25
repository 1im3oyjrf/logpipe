package cap2_test

import (
	"testing"
	"time"

	"github.com/user/logpipe/internal/cap2"
)

func entry(level, msg string) cap2.Entry {
	return cap2.Entry{"level": level, "message": msg}
}

func TestAllow_ZeroMax_ForwardsAll(t *testing.T) {
	c := cap2.New(cap2.Config{Field: "level", Max: 0})
	for i := 0; i < 100; i++ {
		if !c.Allow(entry("info", "msg")) {
			t.Fatal("expected all entries to be forwarded when Max is zero")
		}
	}
}

func TestAllow_WithinMax_Passes(t *testing.T) {
	c := cap2.New(cap2.Config{Field: "level", Max: 3, Window: time.Minute})
	for i := 0; i < 3; i++ {
		if !c.Allow(entry("error", "boom")) {
			t.Fatalf("entry %d should have been allowed", i)
		}
	}
}

func TestAllow_ExceedsMax_Drops(t *testing.T) {
	c := cap2.New(cap2.Config{Field: "level", Max: 2, Window: time.Minute})
	c.Allow(entry("warn", "a"))
	c.Allow(entry("warn", "b"))
	if c.Allow(entry("warn", "c")) {
		t.Fatal("third entry should have been dropped")
	}
}

func TestAllow_DifferentValues_IndependentCounters(t *testing.T) {
	c := cap2.New(cap2.Config{Field: "level", Max: 1, Window: time.Minute})
	if !c.Allow(entry("info", "a")) {
		t.Fatal("first info entry should be allowed")
	}
	if !c.Allow(entry("error", "b")) {
		t.Fatal("first error entry should be allowed independently")
	}
	if c.Allow(entry("info", "c")) {
		t.Fatal("second info entry should be dropped")
	}
}

func TestAllow_CaseInsensitive_TreatsAsSame(t *testing.T) {
	c := cap2.New(cap2.Config{Field: "level", Max: 1, Window: time.Minute, CaseInsensitive: true})
	c.Allow(entry("INFO", "a"))
	if c.Allow(entry("info", "b")) {
		t.Fatal("case-insensitive match should share the same counter")
	}
}

func TestAllow_EmptyField_SharedCounter(t *testing.T) {
	c := cap2.New(cap2.Config{Field: "", Max: 2, Window: time.Minute})
	c.Allow(entry("info", "a"))
	c.Allow(entry("error", "b"))
	if c.Allow(entry("warn", "c")) {
		t.Fatal("third entry of any kind should be dropped under shared counter")
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	c := cap2.New(cap2.Config{Field: "level", Max: 1, Window: time.Minute})
	c.Allow(entry("info", "a"))
	c.Reset()
	if !c.Allow(entry("info", "b")) {
		t.Fatal("entry should be allowed after reset")
	}
}

func TestNew_DefaultWindow_Applied(t *testing.T) {
	c := cap2.New(cap2.Config{Max: 5})
	if c == nil {
		t.Fatal("expected non-nil Capper")
	}
}
