// Package alert implements threshold-based alerting for logpipe.
//
// An Evaluator is configured with one or more Rules. Each Rule specifies
// a log field/value pair to watch, a count threshold, and a sliding time
// window. When the number of matching observations within the window reaches
// the threshold, an Alert is written to the configured io.Writer.
//
// Example usage:
//
//	rule := alert.Rule{
//		Name:      "too-many-errors",
//		Field:     "level",
//		Value:     "error",
//		Level:     alert.LevelError,
//		Threshold: 10,
//		Window:    time.Minute,
//	}
//	ev := alert.New([]alert.Rule{rule}, os.Stderr)
//	ev.Observe(entry.Fields)
package alert
