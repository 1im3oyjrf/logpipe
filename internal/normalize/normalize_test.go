package normalize_test

import (
	"testing"

	"github.com/your-org/logpipe/internal/normalize"
)

func base() map[string]any {
	return map[string]any{
		"Message": "hello",
		"Level":   "WARN",
		"ts":      "2024-01-01T00:00:00Z",
	}
}

func TestApply_NoFieldMap_PassesThrough(t *testing.T) {
	n := normalize.New(normalize.Config{})
	out := n.Apply(base())
	if out["Message"] != "hello" {
		t.Fatalf("expected Message to survive, got %v", out)
	}
}

func TestApply_FieldMap_RenamesKey(t *testing.T) {
	n := normalize.New(normalize.Config{
		FieldMap: map[string]string{"message": "msg"},
	})
	out := n.Apply(base())
	if _, ok := out["Message"]; ok {
		t.Fatal("old key should be gone")
	}
	if out["msg"] != "hello" {
		t.Fatalf("expected msg=hello, got %v", out["msg"])
	}
}

func TestApply_LevelNormalisedToLower(t *testing.T) {
	n := normalize.New(normalize.Config{
		FieldMap:   map[string]string{"level": "level"},
		LevelField: "level",
	})
	out := n.Apply(map[string]any{"level": "ERROR"})
	if out["level"] != "error" {
		t.Fatalf("expected error, got %v", out["level"])
	}
}

func TestApply_CustomLevelField(t *testing.T) {
	n := normalize.New(normalize.Config{LevelField: "severity"})
	out := n.Apply(map[string]any{"severity": "CRITICAL"})
	if out["severity"] != "critical" {
		t.Fatalf("expected critical, got %v", out["severity"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	n := normalize.New(normalize.Config{
		FieldMap: map[string]string{"message": "msg"},
	})
	orig := base()
	_ = n.Apply(orig)
	if _, ok := orig["Message"]; !ok {
		t.Fatal("original entry was mutated")
	}
}

func TestApplyLevel_ReturnsNormalisedLevel(t *testing.T) {
	n := normalize.New(normalize.Config{})
	level := n.ApplyLevel(map[string]any{"level": "INFO"})
	if level != "info" {
		t.Fatalf("expected info, got %q", level)
	}
}

func TestApplyLevel_MissingField_ReturnsEmpty(t *testing.T) {
	n := normalize.New(normalize.Config{})
	level := n.ApplyLevel(map[string]any{"msg": "no level here"})
	if level != "" {
		t.Fatalf("expected empty, got %q", level)
	}
}
