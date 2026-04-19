package tag_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/tag"
)

func base() map[string]any {
	return map[string]any{"level": "error", "msg": "boom"}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	tr := tag.New(tag.Config{})
	out := tr.Apply(base())
	if _, ok := out["tags"]; ok {
		t.Fatal("expected no tags field")
	}
}

func TestApply_MatchingRule_InjectsTags(t *testing.T) {
	tr := tag.New(tag.Config{
		Rules: []tag.Rule{{Field: "level", Value: "error", Tags: []string{"alert", "critical"}}},
	})
	out := tr.Apply(base())
	tags, ok := out["tags"].([]string)
	if !ok {
		t.Fatal("expected tags field")
	}
	if len(tags) != 2 || tags[0] != "alert" || tags[1] != "critical" {
		t.Fatalf("unexpected tags: %v", tags)
	}
}

func TestApply_NonMatchingRule_NoTags(t *testing.T) {
	tr := tag.New(tag.Config{
		Rules: []tag.Rule{{Field: "level", Value: "info", Tags: []string{"quiet"}}},
	})
	out := tr.Apply(base())
	if _, ok := out["tags"]; ok {
		t.Fatal("expected no tags field")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	tr := tag.New(tag.Config{
		CaseInsensitive: true,
		Rules:           []tag.Rule{{Field: "level", Value: "ERROR", Tags: []string{"hi"}}},
	})
	out := tr.Apply(base())
	tags, ok := out["tags"].([]string)
	if !ok || len(tags) != 1 || tags[0] != "hi" {
		t.Fatalf("expected tag hi, got %v", out["tags"])
	}
}

func TestApply_CustomTargetField(t *testing.T) {
	tr := tag.New(tag.Config{
		TargetField: "labels",
		Rules:       []tag.Rule{{Field: "level", Value: "error", Tags: []string{"x"}}},
	})
	out := tr.Apply(base())
	if _, ok := out["labels"]; !ok {
		t.Fatal("expected labels field")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	tr := tag.New(tag.Config{
		Rules: []tag.Rule{{Field: "level", Value: "error", Tags: []string{"x"}}},
	})
	in := base()
	_ = tr.Apply(in)
	if _, ok := in["tags"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestApply_MultipleRules_MergesTags(t *testing.T) {
	tr := tag.New(tag.Config{
		Rules: []tag.Rule{
			{Field: "level", Value: "error", Tags: []string{"a"}},
			{Field: "msg", Value: "boom", Tags: []string{"b"}},
		},
	})
	out := tr.Apply(base())
	tags, ok := out["tags"].([]string)
	if !ok || len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", out["tags"])
	}
}
