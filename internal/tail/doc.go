// Package tail provides a file-tailing utility that follows a log file as new
// content is appended, delivering lines over a channel.
//
// Unlike external tail libraries, this implementation is intentionally minimal:
// it seeks to the end of the file on startup (so only new lines are emitted)
// and polls at a fixed interval when no data is available. This keeps the
// dependency tree small while covering the common logpipe use-case of watching
// a single actively-written log file.
//
// Typical usage:
//
//	tr := tail.New("/var/log/app.log")
//	go tr.Run(ctx)
//	for line := range tr.Lines() {
//		// process line
//	}
package tail
