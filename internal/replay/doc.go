// Package replay provides historical log replay from an in-memory ring buffer.
//
// A Manager wraps a fixed-capacity ring buffer. As new log entries arrive they
// are recorded via Record; once the buffer is full the oldest entries are
// silently evicted. At any point a caller may request a replay of the buffered
// entries, optionally filtered by a grep pattern, by calling Replay which
// returns a channel that is closed once all matching entries have been emitted.
//
// Typical usage:
//
//	mgr := replay.NewManager(500)
//
//	// record entries as they flow through the pipeline
//	mgr.Record(entry)
//
//	// replay the last 500 entries matching "error"
//	ch := mgr.Replay(ctx, replay.Config{Pattern: "error"})
//	for e := range ch {
//		fmt.Println(e.Message)
//	}
package replay
