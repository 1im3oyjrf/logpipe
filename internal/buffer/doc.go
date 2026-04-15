// Package buffer implements a thread-safe fixed-size ring buffer for retaining
// recent log entries in memory.
//
// The Ring type stores up to N entries in chronological order. When the buffer
// is full, the oldest entry is silently overwritten by the newest, ensuring
// constant memory usage regardless of log volume.
//
// Typical usage:
//
//	buf := buffer.New(500)
//
//	// push entries as they arrive from the pipeline
//	buf.Push(entry)
//
//	// retrieve a point-in-time snapshot for display or export
//	recent := buf.Snapshot()
//
// All methods are safe for concurrent use.
package buffer
