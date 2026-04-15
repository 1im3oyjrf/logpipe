// Package aggregator implements time-window based aggregation of structured
// log entries. Entries are grouped by a configurable field (e.g. "level" or
// "service") and counted within a rolling duration window.
//
// Usage:
//
//	agg := aggregator.New("level", 30*time.Second)
//	agg.Add(entry)
//	buckets := agg.Snapshot()
//	for _, b := range buckets {
//		fmt.Printf("%s: %d occurrences\n", b.Key, b.Count)
//	}
//
// Buckets are automatically reset when the first entry in a bucket falls
// outside the configured window, ensuring counts reflect only recent activity.
package aggregator
