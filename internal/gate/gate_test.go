package gate_test

import (
	"testing"

	"github.com/logpipe/internal/gate"
)

func base() map[string]any {
	return map[string]any{"level": "error", "status": "500"}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := gate.New(gate.Config{Field: "", Op: gate.OpEq, Value: "x"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_UnknownOp_ReturnsError(t *testing.T) {
	_, err := gate.New(gate.Config{Field: "level", Op: "contains", Value: "x"})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestAllow_Eq_Match(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "level", Op: gate.OpEq, Value: "error"})
	if !g.Allow(base()) {
		t.Fatal("expected entry to pass")
	}
}

func TestAllow_Eq_NoMatch(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "level", Op: gate.OpEq, Value: "info"})
	if g.Allow(base()) {
		t.Fatal("expected entry to be blocked")
	}
}

func TestAllow_Neq_Match(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "level", Op: gate.OpNeq, Value: "info"})
	if !g.Allow(base()) {
		t.Fatal("expected entry to pass")
	}
}

func TestAllow_CaseInsensitive_Match(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "level", Op: gate.OpEq, Value: "ERROR", CaseInsensitive: true})
	if !g.Allow(base()) {
		t.Fatal("expected case-insensitive match")
	}
}

func TestAllow_MissingField_Neq_Passes(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "missing", Op: gate.OpNeq, Value: "x"})
	if !g.Allow(base()) {
		t.Fatal("missing field should return empty string, which != x")
	}
}

func TestAllow_Gt_Match(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "status", Op: gate.OpGt, Value: "400"})
	if !g.Allow(base()) {
		t.Fatal("expected 500 > 400")
	}
}

func TestAllow_Lt_NoMatch(t *testing.T) {
	g, _ := gate.New(gate.Config{Field: "status", Op: gate.OpLt, Value: "400"})
	if g.Allow(base()) {
		t.Fatal("expected 500 < 400 to fail")
	}
}
