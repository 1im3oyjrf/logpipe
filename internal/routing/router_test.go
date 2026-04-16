package routing_test

import (
	"testing"

	"github.com/user/logpipe/internal/reader"
	"github.com/user/logpipe/internal/routing"
)

func makeEntry(fields map[string]any) reader.Entry {
	return reader.Entry{
		Message: "test message",
		Level:   "info",
		Fields:  fields,
	}
}

func TestDispatch_MatchingRule_SendsToTarget(t *testing.T) {
	rules := []routing.Rule{
		{Field: "level", Value: "error", Target: "errors"},
	}
	r := routing.New(rules, 8)
	defer r.Close()

	entry := makeEntry(map[string]any{"level": "error"})
	r.Dispatch(entry)

	ch := r.Channel("errors")
	if ch == nil {
		t.Fatal("expected errors channel to exist")
	}
	select {
	case got := <-ch:
		if got.Level != entry.Level {
			t.Errorf("expected level %q, got %q", entry.Level, got.Level)
		}
	default:
		t.Error("expected entry in errors channel")
	}
}

func TestDispatch_NoMatch_SendsToDefault(t *testing.T) {
	rules := []routing.Rule{
		{Field: "level", Value: "error", Target: "errors"},
	}
	r := routing.New(rules, 8)
	defer r.Close()

	entry := makeEntry(map[string]any{"level": "info"})
	r.Dispatch(entry)

	ch := r.Channel("default")
	select {
	case got := <-ch:
		if got.Message != entry.Message {
			t.Errorf("unexpected message: %q", got.Message)
		}
	default:
		t.Error("expected entry in default channel")
	}
}

func TestDispatch_CaseInsensitiveMatch(t *testing.T) {
	rules := []routing.Rule{
		{Field: "level", Value: "WARN", Target: "warnings"},
	}
	r := routing.New(rules, 8)
	defer r.Close()

	r.Dispatch(makeEntry(map[string]any{"level": "warn"}))

	ch := r.Channel("warnings")
	select {
	case <-ch:
		// ok
	default:
		t.Error("expected entry in warnings channel")
	}
}

func TestDispatch_MissingField_FallsToDefault(t *testing.T) {
	rules := []routing.Rule{
		{Field: "service", Value: "auth", Target: "auth"},
	}
	r := routing.New(rules, 8)
	defer r.Close()

	r.Dispatch(makeEntry(map[string]any{"level": "info"}))

	ch := r.Channel("default")
	select {
	case <-ch:
		// ok
	default:
		t.Error("expected entry in default channel")
	}
}

func TestChannel_UnknownName_ReturnsNil(t *testing.T) {
	r := routing.New(nil, 8)
	defer r.Close()

	if ch := r.Channel("nonexistent"); ch != nil {
		t.Error("expected nil for unknown channel")
	}
}

func TestTargets_ContainsDefault(t *testing.T) {
	r := routing.New(nil, 8)
	defer r.Close()

	for _, name := range r.Targets() {
		if name == "default" {
			return
		}
	}
	t.Error("expected default to be in targets")
}

func TestTargets_ContainsRuleTargets(t *testing.T) {
	rules := []routing.Rule{
		{Field: "level", Value: "error", Target: "errors"},
		{Field: "service", Value: "auth", Target: "auth"},
	}
	r := routing.New(rules, 8)
	defer r.Close()

	targets := r.Targets()
	want := []string{"errors", "auth", "default"}
	for _, name := range want {
		found := false
		for _, t := range targets {
			if t == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected target %q to be in targets list", name)
		}
	}
}
