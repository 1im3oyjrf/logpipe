// Package pipeline wires together all logpipe components — source
// multiplexing, JSON reading, filtering, and formatted output — into a
// single cohesive processing pipeline.
//
// Usage:
//
//	cfg := &config.Config{
//		Sources: []config.Source{{Reader: os.Stdin, Label: "stdin"}},
//		Pattern: "error",
//		Level:   "warn",
//		NoColor: false,
//		Output:  os.Stdout,
//	}
//
//	if err := pipeline.Run(ctx, cfg); err != nil {
//		log.Fatal(err)
//	}
//
// Run blocks until the context is cancelled or all sources are exhausted.
package pipeline
