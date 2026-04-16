// Package multiwriter provides a fan-out io.Writer that broadcasts each
// write to a dynamic set of underlying writers.
//
// It is safe for concurrent use. Targets can be added or removed at
// runtime without interrupting in-progress writes.
//
// Typical use in logpipe is to route formatted log output to multiple
// simultaneous sinks — e.g. stdout, a file, and a network connection —
// without coupling the pipeline to any specific destination.
package multiwriter
