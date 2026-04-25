// Package hover implements a non-destructive tap stage for the logpipe
// pipeline. It forwards every log entry to the next stage unmodified while
// simultaneously invoking a caller-supplied observer function for each entry.
//
// Typical uses include:
//
//   - Feeding a metrics counter without forking the pipeline.
//   - Attaching a debug printer during development.
//   - Triggering side-effects (e.g. alerting) on specific entries.
//
// The observer is called synchronously in the same goroutine as the pipeline
// read loop. Observers that may block should dispatch to a buffered channel
// or goroutine internally to avoid stalling the pipeline.
package hover
