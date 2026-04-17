package fieldmap_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/fieldmap"
)

func TestApply_NoRules_PassesThrough(t *testing.T) {
	m := fieldmap.New(fieldmap.Config{})
	in := map[string]any{"msg": "hello", "level": "info"}
	out := m.Apply(in)
	if out["msg"] != "hello" || out["level"] != "info" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestApply_RenamesField(t *testing.T) {
	m := fieldmap.New(fieldmap.Config{
		Rules: map[string]string{"msg": "message"},
	})
	out := m.Apply(map[string]any{"msg": "hello", "level": "info"})
	if _, ok := out["msg"]; ok {
		t.Fatal("old key should be absent")
	}
	if out["message"] != "hello" {
		t.Fatalf("expected message=hello, got %v", out["message"])
	}
	if out["level"] != "info" {
		t.Fatal("unrelated field should be preserved")
	}
}

func TestApply_CaseInsensitiveSource(t *testing.T) {
	m := fieldmap.New(fieldmap.Config{
		Rules: map[string]string{"MSG": "message"},
	})
	out := m.Apply(map[string]any{"msg": "hi"})
	if out["message"] != "hi" {
		t.Fatalf("expected case-insensitive match, got %v", out)
	}
}

func TestApply_DropUnmapped(t *testing.T) {
	m := fieldmap.New(fieldmap.Config{
		Rules:        map[string]string{"msg": "message"},
		DropUnmapped: true,
	})
	out := m.Apply(map[string]any{"msg": "hello", "extra": "drop me"})
	if _, ok := out["extra"]; ok {
		t.Fatal("unmapped field should be dropped")
	}
	if out["message"] != "hello" {
		t.Fatal("mapped field should be present")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	m := fieldmap.New(fieldmap.Config{
		Rules: map[string]string{"msg": "message"},
	})
	in := map[string]any{"msg": "hello"}
	m.Apply(in)
	if _, ok := in["msg"]; !ok {
		t.Fatal("original map must not be mutated")
	}
}

func TestRules_ReturnsCopy(t *testing.T) {
	m := fieldmap.New(fieldmap.Config{
		Rules: map[string]string{"msg": "message"},
	})
	r := m.Rules()
	r["injected"] = "bad"
	if _, ok := m.Rules()["injected"]; ok {
		t.Fatal("Rules() should return an independent copy")
	}
}

func TestApply_MultipleRenames(t *testing.T) {
	m := fieldmap.New(fieldmap.Config{
		Rules: map[string]string{
			"msg":   "message",
			"lvl":   "level",
			"ts":    "timestamp",
		},
	})
	out := m.Apply(map[string]any{"msg": "hello", "lvl": "warn", "ts": "2024-01-01"})
	if out["message"] != "hello" || out["level"] != "warn" || out["timestamp"] != "2024-01-01" {
		t.Fatalf("expected all fields renamed, got %v", out)
	}
	for _, old := range []string{"msg", "lvl", "ts"} {
		if _, ok := out[old]; ok {
			t.Fatalf("old key %q should be absent", old)
		}
	}
}
