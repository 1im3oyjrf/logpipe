package ceiling_test

import (
	"testing"
	"time"

	"github.com/your-org/logpipe/internal/ceiling"
)

func entry() map[string]interface{} {
	return map[string]interface{}{"msg": "hello", "level": "info"}
}

func TestAllow_WithinMax_Passes(t *testing.T) {
	c := ceiling.New(ceiling.Config{Max: 3, Window: time.Second})
	for i := 0; i < 3; i++ {
		if !c.Allow(entry()) {
			t.Fatalf("expected entry %d to pass", i)
		}
	}
}

func TestAllow_ExceedsMax_Drops(t *testing.T) {
	c := ceiling.New(ceiling.Config{Max: 2, Window: time.Second})
	c.Allow(entry())
	c.Allow(entry())
	if c.Allow(entry()) {
		t.Fatal("expected third entry to be dropped")
	}
	if c.Dropped() != 1 {
		t.Fatalf("expected dropped=1, got %d", c.Dropped())
	}
}

func TestAllow_WindowExpiry_ResetsCount(t *testing.T) {
	c := ceiling.New(ceiling.Config{Max: 1, Window: 20 * time.Millisecond})
	if !c.Allow(entry()) {
		t.Fatal("first entry should pass")
	}
	if c.Allow(entry()) {
		t.Fatal("second entry should be dropped")
	}
	time.Sleep(30 * time.Millisecond)
	if !c.Allow(entry()) {
		t.Fatal("entry after window expiry should pass")
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	c := ceiling.New(ceiling.Config{Max: 1, Window: time.Second})
	c.Allow(entry())
	c.Allow(entry()) // dropped
	if c.Dropped() != 1 {
		t.Fatal("expected 1 dropped before reset")
	}
	c.Reset()
	if c.Dropped() != 0 {
		t.Fatal("expected 0 dropped after reset")
	}
	if !c.Allow(entry()) {
		t.Fatal("entry after reset should pass")
	}
}

func TestNew_PanicOnZeroMax(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for Max=0")
		}
	}()
	ceiling.New(ceiling.Config{Max: 0})
}

func TestApply_ChannelFiltering(t *testing.T) {
	c := ceiling.New(ceiling.Config{Max: 2, Window: time.Second})
	in := make(chan map[string]interface{}, 5)
	for i := 0; i < 5; i++ {
		in <- entry()
	}
	close(in)
	out := c.Apply(in)
	var count int
	for range out {
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 entries through channel, got %d", count)
	}
}
