package sampler_test

import (
	"testing"

	"github.com/your-org/logpipe/internal/reader"
	"github.com/your-org/logpipe/internal/sampler"
)

func entry(msg string) reader.Entry {
	return reader.Entry{Message: msg, Level: "info"}
}

func TestNew_RateOne_KeepsAll(t *testing.T) {
	s := sampler.New(1)
	for i := 0; i < 20; i++ {
		if !s.Keep(entry("x")) {
			t.Fatalf("rate=1 should keep every entry, dropped at i=%d", i)
		}
	}
	if s.Dropped() != 0 {
		t.Errorf("expected 0 dropped, got %d", s.Dropped())
	}
}

func TestNew_RateZero_TreatedAsOne(t *testing.T) {
	s := sampler.New(0)
	if s.Rate() != 1 {
		t.Errorf("expected rate normalised to 1, got %d", s.Rate())
	}
}

func TestKeep_RateTen_KeepsOneInTen(t *testing.T) {
	s := sampler.New(10)
	kept := 0
	const total = 100
	for i := 0; i < total; i++ {
		if s.Keep(entry("x")) {
			kept++
		}
	}
	if kept != 10 {
		t.Errorf("expected 10 kept entries, got %d", kept)
	}
	if s.Dropped() != 90 {
		t.Errorf("expected 90 dropped, got %d", s.Dropped())
	}
}

func TestKeep_RateTwo_AlternatesKeepDrop(t *testing.T) {
	s := sampler.New(2)
	results := make([]bool, 6)
	for i := range results {
		results[i] = s.Keep(entry("x"))
	}
	expected := []bool{true, false, true, false, true, false}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("index %d: expected %v, got %v", i, expected[i], got)
		}
	}
}

func TestDropped_AccumulatesAcrossCalls(t *testing.T) {
	s := sampler.New(5)
	for i := 0; i < 25; i++ {
		s.Keep(entry("x"))
	}
	if s.Dropped() != 20 {
		t.Errorf("expected 20 dropped, got %d", s.Dropped())
	}
}
