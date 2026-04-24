package cascade_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/cascade"
)

func base() map[string]any {
	return map[string]any{
		"level":   "error",
		"message": "something went wrong",
		"service": "api",
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	c, err := cascade.New(cascade.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := c.Apply(base())
	if out["level"] != "error" {
		t.Errorf("expected level=error, got %v", out["level"])
	}
}

func TestApply_MatchingRule_InjectsTarget(t *testing.T) {
	c, _ := cascade.New(cascade.Config{
		Rules: []cascade.Rule{
			{Field: "level", Value: "error", Target: "priority", Set: "high"},
		},
		StopOnFirst: true,
	})
	out := c.Apply(base())
	if out["priority"] != "high" {
		t.Errorf("expected priority=high, got %v", out["priority"])
	}
}

func TestApply_NonMatchingRule_NoInjection(t *testing.T) {
	c, _ := cascade.New(cascade.Config{
		Rules: []cascade.Rule{
			{Field: "level", Value: "warn", Target: "priority", Set: "medium"},
		},
		StopOnFirst: true,
	})
	out := c.Apply(base())
	if _, ok := out["priority"]; ok {
		t.Error("expected no priority field to be injected")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	c, _ := cascade.New(cascade.Config{
		Rules: []cascade.Rule{
			{Field: "level", Value: "ERROR", Target: "priority", Set: "high"},
		},
		StopOnFirst: true,
	})
	out := c.Apply(base())
	if out["priority"] != "high" {
		t.Errorf("expected priority=high, got %v", out["priority"])
	}
}

func TestApply_StopOnFirst_OnlyFirstRuleApplied(t *testing.T) {
	c, _ := cascade.New(cascade.Config{
		Rules: []cascade.Rule{
			{Field: "level", Value: "error", Target: "priority", Set: "high"},
			{Field: "level", Value: "error", Target: "flag", Set: "yes"},
		},
		StopOnFirst: true,
	})
	out := c.Apply(base())
	if out["priority"] != "high" {
		t.Errorf("expected priority=high, got %v", out["priority"])
	}
	if _, ok := out["flag"]; ok {
		t.Error("expected flag not to be set when StopOnFirst=true")
	}
}

func TestApply_AllRules_MultipleTargetsSet(t *testing.T) {
	c, _ := cascade.New(cascade.Config{
		Rules: []cascade.Rule{
			{Field: "level", Value: "error", Target: "priority", Set: "high"},
			{Field: "level", Value: "error", Target: "flag", Set: "yes"},
		},
		StopOnFirst: false,
	})
	out := c.Apply(base())
	if out["priority"] != "high" {
		t.Errorf("expected priority=high, got %v", out["priority"])
	}
	if out["flag"] != "yes" {
		t.Errorf("expected flag=yes, got %v", out["flag"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	c, _ := cascade.New(cascade.Config{
		Rules: []cascade.Rule{
			{Field: "level", Value: "error", Target: "priority", Set: "high"},
		},
		StopOnFirst: true,
	})
	orig := base()
	c.Apply(orig)
	if _, ok := orig["priority"]; ok {
		t.Error("original entry must not be mutated")
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := cascade.New(cascade.Config{
		Rules: []cascade.Rule{
			{Field: "", Value: "error", Target: "priority", Set: "high"},
		},
	})
	if err == nil {
		t.Error("expected error for empty field")
	}
}

func TestNew_EmptyTarget_ReturnsError(t *testing.T) {
	_, err := cascade.New(cascade.Config{
		Rules: []cascade.Rule{
			{Field: "level", Value: "error", Target: "", Set: "high"},
		},
	})
	if err == nil {
		t.Error("expected error for empty target")
	}
}
