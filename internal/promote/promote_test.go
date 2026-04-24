package promote_test

import (
	"testing"

	"logpipe/internal/promote"
)

var base = map[string]any{
	"level":   "info",
	"message": "hello",
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := promote.New(promote.Config{})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestApply_ScalarField_PromotedToTarget(t *testing.T) {
	p, _ := promote.New(promote.Config{Field: "level", DropSource: true})
	entry := map[string]any{"level": "warn", "msg": "hi"}
	out := p.Apply(entry)
	if out["level"] != "warn" {
		t.Fatalf("expected level=warn, got %v", out["level"])
	}
}

func TestApply_MapField_MergedToTopLevel(t *testing.T) {
	p, _ := promote.New(promote.Config{Field: "meta", DropSource: true})
	entry := map[string]any{
		"msg":  "test",
		"meta": map[string]any{"host": "srv1", "env": "prod"},
	}
	out := p.Apply(entry)
	if out["host"] != "srv1" {
		t.Errorf("expected host=srv1, got %v", out["host"])
	}
	if out["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", out["env"])
	}
	if _, ok := out["meta"]; ok {
		t.Error("expected meta to be removed")
	}
}

func TestApply_DropSourceFalse_KeepsOriginal(t *testing.T) {
	p, _ := promote.New(promote.Config{Field: "meta", DropSource: false})
	entry := map[string]any{
		"msg":  "test",
		"meta": map[string]any{"host": "srv1"},
	}
	out := p.Apply(entry)
	if _, ok := out["meta"]; !ok {
		t.Error("expected meta to be retained when DropSource=false")
	}
	if out["host"] != "srv1" {
		t.Errorf("expected host=srv1, got %v", out["host"])
	}
}

func TestApply_MissingField_PassesThrough(t *testing.T) {
	p, _ := promote.New(promote.Config{Field: "nonexistent", DropSource: true})
	entry := map[string]any{"msg": "hello"}
	out := p.Apply(entry)
	if out["msg"] != "hello" {
		t.Errorf("expected msg to be preserved")
	}
	if len(out) != 1 {
		t.Errorf("expected no extra keys, got %v", out)
	}
}

func TestApply_CaseInsensitiveField_Matches(t *testing.T) {
	p, _ := promote.New(promote.Config{Field: "META", DropSource: true, CaseInsensitive: true})
	entry := map[string]any{
		"msg":  "test",
		"meta": map[string]any{"region": "us-east"},
	}
	out := p.Apply(entry)
	if out["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %v", out["region"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	p, _ := promote.New(promote.Config{Field: "meta", DropSource: true})
	original := map[string]any{
		"msg":  "test",
		"meta": map[string]any{"host": "srv1"},
	}
	p.Apply(original)
	if _, ok := original["meta"]; !ok {
		t.Error("original entry should not be mutated")
	}
}

func TestApply_DefaultTarget_UsesLeafSegment(t *testing.T) {
	p, _ := promote.New(promote.Config{Field: "context"})
	entry := map[string]any{"context": "request-id-42"}
	out := p.Apply(entry)
	// scalar: target defaults to "context"
	if out["context"] != "request-id-42" {
		t.Errorf("expected context=request-id-42, got %v", out["context"])
	}
}
