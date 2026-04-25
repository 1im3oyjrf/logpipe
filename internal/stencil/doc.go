// Package stencil stamps a fixed set of key-value fields onto every log entry
// that flows through the pipeline.
//
// A Stencil is constructed from a Config that declares the fields to inject,
// whether existing fields may be overwritten, and whether key matching should
// be case-insensitive.
//
// Usage:
//
//	s, err := stencil.New(stencil.Config{
//		Fields:    map[string]any{"env": "production", "region": "eu-west-1"},
//		Overwrite: false,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	enriched := s.Apply(entry)
package stencil
