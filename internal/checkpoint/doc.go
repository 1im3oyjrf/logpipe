// Package checkpoint implements persistent offset tracking for logpipe sources.
//
// When logpipe tails one or more files it records the byte offset of the last
// line successfully processed. On restart the stored offsets are read back so
// that each source resumes exactly where it left off, preventing both duplicate
// delivery and missed lines.
//
// # Usage
//
//	store, err := checkpoint.New("/var/lib/logpipe/checkpoint.json")
//	if err != nil { ... }
//
//	// Read the last saved position for a source.
//	offset := store.Get("/var/log/app.log")
//
//	// After processing a line, advance the offset.
//	store.Set("/var/log/app.log", newOffset)
//
//	// Flush periodically via Manager.
//	mgr := checkpoint.NewManager(store, 5*time.Second)
//	mgr.Run(ctx)
package checkpoint
