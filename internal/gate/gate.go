// Package gate provides a conditional pass-through filter that forwards
// log entries only when a named field satisfies a configured comparison.
package gate

import (
	"fmt"
	"strings"

	"github.com/logpipe/internal/parser"
)

// Op is a comparison operator.
type Op string

const (
	OpEq  Op = "eq"
	OpNeq Op = "neq"
	OpGt  Op = "gt"
	OpLt  Op = "lt"
)

// Config holds the gate configuration.
type Config struct {
	Field           string
	Op              Op
	Value           string
	CaseInsensitive bool
}

// Gate forwards entries that satisfy the configured condition.
type Gate struct {
	cfg Config
}

// New returns a Gate or an error if the configuration is invalid.
func New(cfg Config) (*Gate, error) {
	if strings.TrimSpace(cfg.Field) == "" {
		return nil, fmt.Errorf("gate: field must not be empty")
	}
	switch cfg.Op {
	case OpEq, OpNeq, OpGt, OpLt:
	default:
		return nil, fmt.Errorf("gate: unknown operator %q", cfg.Op)
	}
	return &Gate{cfg: cfg}, nil
}

// Allow returns true when the entry passes the gate condition.
func (g *Gate) Allow(entry map[string]any) bool {
	v := parser.GetString(entry, g.cfg.Field)
	actual := v
	expected := g.cfg.Value
	if g.cfg.CaseInsensitive {
		actual = strings.ToLower(actual)
		expected = strings.ToLower(expected)
	}
	switch g.cfg.Op {
	case OpEq:
		return actual == expected
	case OpNeq:
		return actual != expected
	case OpGt:
		return actual > expected
	case OpLt:
		return actual < expected
	}
	return false
}
