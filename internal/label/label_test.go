package label_test

import (
	"testing"

	"github.com/logpipe/internal/label"
)

func base() map[string]any {
	return map[string]any{
		"level":   "error",
		"message": "disk full on /var",
		"service": "storage",
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	l := label.New(nil)
	out := l.Apply(base())
	if out["level"] != "error" {
		t.Fatalf("expected level=error, got %v", out["level"])
	}
}

func TestApply_MatchingRule_InjectsLabels(t *testing.T) {
	rules := []label.Rule{
		{Field: "level", Value: "error", Labels: map[string]string{"team": "oncall"}},
	}
	l := label.New(rules)
	out := l.Apply(base())
	if out["team"] != "oncall" {
		t.Fatalf("expected team=oncall, got %v", out["team"])
	}
}

func TestApply_NonMatchingRule_NoLabel(t *testing.T) {
	rules := []label.Rule{
		{Field: "level", Value: "warn", Labels: map[string]string{"team": "oncall"}},
	}
	l := label.New(rules)
	out := l.Apply(base())
	if _, ok := out["team"]; ok {
		t.Fatal("expected no team label")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	rules := []label.Rule{
		{Field: "message", Value: "DISK", Labels: map[string]string{"alert": "disk"}},
	}
	l := label.New(rules)
	out := l.Apply(base())
	if out["alert"] != "disk" {
		t.Fatalf("expected alert=disk, got %v", out["alert"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	rules := []label.Rule{
		{Field: "level", Value: "error", Labels: map[string]string{"team": "oncall"}},
	}
	l := label.New(rules)
	orig := base()
	l.Apply(orig)
	if _, ok := orig["team"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestApply_MultipleRules_AllMatching(t *testing.T) {
	rules := []label.Rule{
		{Field: "level", Value: "error", Labels: map[string]string{"severity": "high"}},
		{Field: "service", Value: "storage", Labels: map[string]string{"domain": "infra"}},
	}
	l := label.New(rules)
	out := l.Apply(base())
	if out["severity"] != "high" {
		t.Fatalf("expected severity=high")
	}
	if out["domain"] != "infra" {
		t.Fatalf("expected domain=infra")
	}
}
