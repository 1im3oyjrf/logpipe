// Package output provides terminal rendering for structured log entries.
//
// The Formatter type accepts a decoded JSON log entry (map[string]interface{})
// and writes a human-readable, optionally colourised line to any io.Writer.
//
// Well-known fields — time/ts/timestamp, level/lvl, msg/message — are given
// dedicated positions in the output line. All remaining fields are appended as
// key=value pairs so that no information is lost.
//
// Colour output relies on github.com/fatih/color and is automatically
// suppressed when the NO_COLOR environment variable is set or when
// Options.NoColor is true.
//
// Example usage:
//
//	f := output.New(os.Stdout, output.Options{Source: "app.log"})
//	f.Write(entry) // entry is a map[string]interface{} from the JSON reader
package output
