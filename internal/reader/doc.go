// Package reader provides utilities for reading and parsing structured
// JSON log entries from arbitrary io.Reader sources.
//
// Usage:
//
//	file, _ := os.Open("app.log")
//	r := reader.NewJSONReader("app", file)
//	out := make(chan reader.LogEntry, 100)
//	go r.ReadAll(out)
//	for entry := range out {
//		fmt.Println(entry.Fields)
//	}
//
// Each LogEntry carries the raw line, parsed fields, source label, and
// an optional timestamp extracted from well-known field names
// ("time", "timestamp", "ts", "@timestamp").
//
// Lines that are not valid JSON are not silently dropped; they are emitted
// as LogEntry values with a nil Fields map and the parse error recorded in
// the Err field, allowing callers to decide how to handle malformed input.
package reader
