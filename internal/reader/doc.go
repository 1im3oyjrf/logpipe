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
package reader
