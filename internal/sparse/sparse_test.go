package sparse_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/sparse"
)

func base(level, msg string) sparse.Entry {
	return sparse.Entry{"level": level, "message": msg}
}

func TestNew_NegativeEvery_ReturnsError(t *testing.T) {
	_, err := sparse.New(sparse.Config{Every: -1})
	if err == nil {
		t.Fatal("expected error for negative Every")
	}
}

func TestNew_ZeroEvery_TreatedAsOne(t *testing.T) {
	s, err := sparse.New(sparse.Config{Every: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Every entry should pass when effective rate is 1.
	for i := 0; i < 5; i++ {
		if !s.Allow(base("info", "msg")) {
			t.Errorf("entry %d should pass with Every=1", i)
		}
	}
}

func TestAllow_EveryOne_KeepsAll(t *testing.T) {
	s, _ := sparse.New(sparse.Config{Every: 1})
	for i := 0; i < 10; i++ {
		if !s.Allow(base("info", "msg")) {
			t.Errorf("entry %d should be kept", i)
		}
	}
}

func TestAllow_EveryThree_KeepsOneInThree(t *testing.T) {
	s, _ := sparse.New(sparse.Config{Every: 3})
	kept := 0
	for i := 0; i < 9; i++ {
		if s.Allow(base("info", "msg")) {
			kept++
		}
	}
	if kept != 3 {
		t.Errorf("expected 3 kept entries, got %d", kept)
	}
}

func TestAllow_DifferentBuckets_IndependentCounters(t *testing.T) {
	s, _ := sparse.New(sparse.Config{Every: 2})
	// First call per bucket should always pass.
	if !s.Allow(base("info", "a")) {
		t.Error("first info entry should pass")
	}
	if !s.Allow(base("error", "b")) {
		t.Error("first error entry should pass")
	}
	// Second call in each bucket should be dropped.
	if s.Allow(base("info", "c")) {
		t.Error("second info entry should be dropped")
	}
	if s.Allow(base("error", "d")) {
		t.Error("second error entry should be dropped")
	}
}

func TestAllow_CaseInsensitive_TreatsAsSameBucket(t *testing.T) {
	s, _ := sparse.New(sparse.Config{Every: 2, CaseInsensitive: true})
	// "INFO" and "info" should share the same counter.
	if !s.Allow(base("INFO", "first")) {
		t.Error("first entry should pass")
	}
	if s.Allow(base("info", "second")) {
		t.Error("second entry (same bucket, case-folded) should be dropped")
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	s, _ := sparse.New(sparse.Config{Every: 2})
	s.Allow(base("info", "a")) // counter → 1 (pass)
	s.Allow(base("info", "b")) // counter → 2 (drop)
	s.Reset()
	if !s.Allow(base("info", "c")) {
		t.Error("after reset, first entry should pass again")
	}
}

func TestAllow_DefaultField_IsLevel(t *testing.T) {
	// No Field specified — should default to "level".
	s, _ := sparse.New(sparse.Config{Every: 3})
	pass := 0
	for i := 0; i < 6; i++ {
		if s.Allow(sparse.Entry{"level": "warn", "msg": "x"}) {
			pass++
		}
	}
	if pass != 2 {
		t.Errorf("expected 2 passes in 6 calls with Every=3, got %d", pass)
	}
}
