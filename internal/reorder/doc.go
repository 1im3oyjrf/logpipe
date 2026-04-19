// Package reorder implements a log-entry transformation stage that promotes
// a configured list of field keys to the front of each entry's field map.
//
// # Usage
//
// Construct a Reorder with a Config that lists the fields to promote:
//
//	r := reorder.New(reorder.Config{
//		Fields:          []string{"level", "msg", "ts"},
//		CaseInsensitive: true,
//	})
//
// Apply it to each entry as it flows through the pipeline:
//
//	out := r.Apply(entry)
//
// Fields that appear in Config.Fields are written first (in the declared
// order); all remaining fields follow in their original iteration order.
// Fields named in Config.Fields that are absent from the entry are silently
// skipped. The original entry is never mutated.
package reorder
