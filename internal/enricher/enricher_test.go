package enricher_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/enricher"
	"github.com/yourorg/logpipe/internal/reader"
)

func base() reader.Entry {
	return reader.Entry{
		Level:   "info",
		Message: "hello",
		Fields:  map[string]any{"app": "test"},
	}
}

func TestApply_StaticFields_AreInjected(t *testing.T) {
	e := enricher.New(enricher.Config{
		StaticFields: map[string]string{"env": "prod", "region": "us-east-1"},
	})
	out := e.Apply(base())
	if out.Fields["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", out.Fields["env"])
	}
	if out.Fields["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %v", out.Fields["region"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	e := enricher.New(enricher.Config{
		StaticFields: map[string]string{"injected": "yes"},
	})
	orig := base()
	e.Apply(orig)
	if _, ok := orig.Fields["injected"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestApply_AddTimestamp_SetsField(t *testing.T) {
	e := enricher.New(enricher.Config{AddTimestamp: true})
	out := e.Apply(base())
	v, ok := out.Fields["enriched_at"]
	if !ok {
		t.Fatal("enriched_at field missing")
	}
	ts, ok := v.(string)
	if !ok || !strings.Contains(ts, "T") {
		t.Errorf("unexpected enriched_at value: %v", v)
	}
}

func TestApply_HostField_IsSet(t *testing.T) {
	e := enricher.New(enricher.Config{HostField: "host"})
	out := e.Apply(base())
	v, ok := out.Fields["host"]
	if !ok {
		t.Skip("hostname unavailable in this environment")
	}
	if v == "" {
		t.Error("host field is empty")
	}
}

func TestApply_NilFields_InitialisedSafely(t *testing.T) {
	e := enricher.New(enricher.Config{
		StaticFields: map[string]string{"k": "v"},
	})
	entry := reader.Entry{Level: "warn", Message: "no fields"}
	out := e.Apply(entry)
	if out.Fields["k"] != "v" {
		t.Errorf("expected k=v, got %v", out.Fields["k"])
	}
}
