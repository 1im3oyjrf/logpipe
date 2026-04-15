// Package source manages named log input sources and provides a
// Multiplexer that fans in entries from multiple concurrent readers
// into a single unified stream.
//
// Typical usage:
//
//	srcA := source.New("api", fileA)
//	srcB := source.New("worker", fileB)
//	mux  := source.NewMultiplexer(srcA, srcB)
//
//	for entry := range mux.Stream(ctx) {
//		fmt.Println(entry.Source, entry.Fields)
//	}
//
// Each source is read in its own goroutine; the output channel is
// closed automatically once all sources reach EOF or the context is
// cancelled.
package source
