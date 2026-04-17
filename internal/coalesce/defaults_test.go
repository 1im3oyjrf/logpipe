package coalesce_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/coalesce"
)

func TestNew_EmptyConfig_NoRules(t *testing.T) {
	tr := coalesce.New(coalesce.Config{})
	if tr == nil {
		t.Fatal("expected non-nil transformer")
	}
}

func TestNew_RulesStored(t *testing.T) {
	cfg := coalesce.Config{
		Rules: []coalesce.Rule{
			{Sources: []string{"a", "b"}, Target: "c"},
			{Sources: []string{"x"}, Target: "y", KeepSources: true},
		},
	}
	tr := coalesce.New(cfg)
	rules := coalesce.RulesOf(tr)
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Target != "c" {
		t.Fatalf("expected first target=c, got %s", rules[0].Target)
	}
	if rules[1].KeepSources != true {
		t.Fatal("expected KeepSources=true on second rule")
	}
}
