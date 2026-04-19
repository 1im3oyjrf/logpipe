// Package sequence provides a Sequencer that stamps each log entry with a
// monotonically increasing integer field.
//
// This is useful for preserving the original ingestion order of entries when
// they are later sorted, batched, or routed through concurrent pipelines.
//
// Usage:
//
//	seq := sequence.New(sequence.Config{Field: "_seq"})
//	outEntry := seq.Apply(inEntry)
//
// The counter is safe for concurrent use.
package sequence
