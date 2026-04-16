// Package schema provides JSON log entry validation against
// a configurable set of required and optional fields.
package schema

import (
	"fmt"
	"strings"
)

// Config holds the schema validation configuration.
type Config struct {
	// RequiredFields lists field names that must be present in every log entry.
	RequiredFields []string
	// AllowUnknown controls whether fields not listed in KnownFields are permitted.
	AllowUnknown bool
	// KnownFields is the exhaustive list of permitted field names when AllowUnknown is false.
	KnownFields []string
}

// Validator checks log entries against a schema definition.
type Validator struct {
	cfg        Config
	required   map[string]struct{}
	known      map[string]struct{}
}

// New creates a Validator from the provided Config.
func New(cfg Config) *Validator {
	v := &Validator{cfg: cfg}

	v.required = make(map[string]struct{}, len(cfg.RequiredFields))
	for _, f := range cfg.RequiredFields {
		v.required[strings.ToLower(f)] = struct{}{}
	}

	v.known = make(map[string]struct{}, len(cfg.KnownFields))
	for _, f := range cfg.KnownFields {
		v.known[strings.ToLower(f)] = struct{}{}
	}

	return v
}

// Validate checks the provided map of fields against the schema.
// It returns a non-nil error describing the first violation found.
func (v *Validator) Validate(entry map[string]any) error {
	for req := range v.required {
		found := false
		for k := range entry {
			if strings.ToLower(k) == req {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("schema: required field %q is missing", req)
		}
	}

	if !v.cfg.AllowUnknown && len(v.known) > 0 {
		for k := range entry {
			if _, ok := v.known[strings.ToLower(k)]; !ok {
				return fmt.Errorf("schema: unknown field %q", k)
			}
		}
	}

	return nil
}
