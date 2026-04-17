package remap_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/remap"
)

func TestNew_EmptyConfig_NoRules(t *testing.T) {
	r := remap.New(remap.Config{})
	if len(remap.RulesOf(r)) != 0 {
		t.Fatal("expected no rules")
	}
}

func TestNew_RulesStored(t *testing.T) {
	cfg := remap.Config{
		Rules: []remap.Rule{
			{Field: "level", From: "WARN", To: "warning"},
			{Field: "status", From: "ERR", To: "error"},
		},
	}
	r := remap.New(cfg)
	if len(remap.RulesOf(r)) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(remap.RulesOf(r)))
	}
}

func TestNew_RulesAreCopied(t *testing.T) {
	rules := []remap.Rule{
		{Field: "level", From: "WARN", To: "warning"},
	}
	cfg := remap.Config{Rules: rules}
	r := remap.New(cfg)
	rules[0].To = "mutated"
	if remap.RulesOf(r)[0].To == "mutated" {
		t.Fatal("rules slice was not copied")
	}
}
