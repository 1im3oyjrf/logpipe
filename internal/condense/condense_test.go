package condense_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/condense"
)

func entry(msg, level string) map[string]interface{} {
	return map[string]interface{}{"message": msg, "level": level}
}

func TestApply_FirstEntry_BufferedNotEmitted(t *testing.T) {
	c := condense.New(condense.Config{})
	out, ok := c.Apply(entry("hello", "info"))
	if ok || out != nil {
		t.Fatal("expected first entry to be buffered")
	}
}

func TestApply_DifferentEntry_FlushesFirst(t *testing.T) {
	c := condense.New(condense.Config{})
	c.Apply(entry("hello", "info"))
	out, ok := c.Apply(entry("world", "info"))
	if !ok {
		t.Fatal("expected flush on distinct entry")
	}
	if out["message"] != "hello" {
		t.Fatalf("expected flushed message=hello, got %v", out["message"])
	}
}

func TestApply_RepeatedEntry_NoRepeatField(t *testing.T) {
	c := condense.New(condense.Config{})
	c.Apply(entry("hello", "info"))
	out, ok := c.Apply(entry("hello", "info"))
	if ok {
		t.Fatal("expected repeat to be suppressed")
	}
	_ = out
}

func TestFlush_EmitsAccumulatedRepeat(t *testing.T) {
	c := condense.New(condense.Config{})
	c.Apply(entry("hello", "info"))
	c.Apply(entry("hello", "info"))
	c.Apply(entry("hello", "info"))
	out, ok := c.Flush()
	if !ok {
		t.Fatal("expected flush to emit entry")
	}
	if out["_repeat"] != 3 {
		t.Fatalf("expected _repeat=3, got %v", out["_repeat"])
	}
}

func TestFlush_SingleEntry_NoRepeatField(t *testing.T) {
	c := condense.New(condense.Config{})
	c.Apply(entry("hello", "info"))
	out, ok := c.Flush()
	if !ok {
		t.Fatal("expected flush to emit entry")
	}
	if _, exists := out["_repeat"]; exists {
		t.Fatal("expected no _repeat field for single entry")
	}
}

func TestFlush_Empty_ReturnsFalse(t *testing.T) {
	c := condense.New(condense.Config{})
	_, ok := c.Flush()
	if ok {
		t.Fatal("expected false on empty flush")
	}
}

func TestApply_CustomCountField(t *testing.T) {
	c := condense.New(condense.Config{CountField: "count"})
	c.Apply(entry("x", "warn"))
	c.Apply(entry("x", "warn"))
	out, _ := c.Flush()
	if out["count"] != 2 {
		t.Fatalf("expected count=2, got %v", out["count"])
	}
}

func TestApply_MaxRepeat_ForcesFlush(t *testing.T) {
	c := condense.New(condense.Config{MaxRepeat: 3})
	c.Apply(entry("loop", "debug"))
	c.Apply(entry("loop", "debug"))
	out, ok := c.Apply(entry("loop", "debug"))
	if !ok {
		t.Fatal("expected forced flush at MaxRepeat")
	}
	if out["_repeat"] != 3 {
		t.Fatalf("expected _repeat=3, got %v", out["_repeat"])
	}
}
