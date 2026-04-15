// Package metrics provides lightweight, thread-safe counters for tracking
// logpipe pipeline activity at runtime.
//
// A Metrics instance records lines read, matched, dropped, and errors
// encountered during log processing. Callers obtain an immutable point-in-time
// view via Snapshot, which also exposes derived rates such as MatchRate and
// DropRate.
//
// The Reporter type wraps a Metrics instance and periodically writes human-
// readable summaries to any io.Writer, making it easy to surface live stats
// to stderr or a log file without coupling the pipeline to a specific sink.
package metrics
