// Package output handles rendering of log entries to the terminal.
package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Level colors for log level highlighting.
var (
	levelColors = map[string]*color.Color{
		"debug":   color.New(color.FgCyan),
		"info":    color.New(color.FgGreen),
		"warn":    color.New(color.FgYellow),
		"warning": color.New(color.FgYellow),
		"error":   color.New(color.FgRed),
		"fatal":   color.New(color.FgRed, color.Bold),
	}
	defaultColor = color.New(color.FgWhite)
	keyColor     = color.New(color.FgBlue)
	timeColor    = color.New(color.FgMagenta)
)

// Options controls formatter behaviour.
type Options struct {
	// NoColor disables ANSI color output.
	NoColor bool
	// TimeFormat is the Go time layout used to render timestamps.
	// Defaults to time.RFC3339 when empty.
	TimeFormat string
	// Source is an optional label printed before each line (e.g. filename).
	Source string
}

// Formatter writes human-readable log lines to an io.Writer.
type Formatter struct {
	w    io.Writer
	opts Options
}

// New creates a Formatter that writes to w.
func New(w io.Writer, opts Options) *Formatter {
	if opts.NoColor {
		color.NoColor = true
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.RFC3339
	}
	return &Formatter{w: w, opts: opts}
}

// Write renders a single log entry as a formatted line.
// entry is expected to be a map decoded from a JSON log line.
func (f *Formatter) Write(entry map[string]interface{}) {
	ts := f.extractTime(entry)
	level := f.extractLevel(entry)
	msg := f.extractMessage(entry)

	lc := levelColor(level)

	var sb strings.Builder
	if f.opts.Source != "" {
		sb.WriteString(keyColor.Sprintf("[%s] ", f.opts.Source))
	}
	sb.WriteString(timeColor.Sprintf("%s ", ts))
	sb.WriteString(lc.Sprintf("%-7s ", strings.ToUpper(level)))
	sb.WriteString(defaultColor.Sprint(msg))

	// Append remaining fields.
	for k, v := range entry {
		switch k {
		case "time", "ts", "timestamp", "level", "lvl", "msg", "message":
			continue
		}
		sb.WriteString(keyColor.Sprintf(" %s=", k))
		sb.WriteString(fmt.Sprintf("%v", v))
	}

	fmt.Fprintln(f.w, sb.String())
}

func (f *Formatter) extractTime(e map[string]interface{}) string {
	for _, k := range []string{"time", "ts", "timestamp"} {
		if v, ok := e[k]; ok {
			switch t := v.(type) {
			case string:
				return t
			case float64:
				return time.Unix(int64(t), 0).Format(f.opts.TimeFormat)
			}
		}
	}
	return time.Now().Format(f.opts.TimeFormat)
}

func (f *Formatter) extractLevel(e map[string]interface{}) string {
	for _, k := range []string{"level", "lvl"} {
		if v, ok := e[k]; ok {
			if s, ok := v.(string); ok {
				return strings.ToLower(s)
			}
		}
	}
	return "info"
}

func (f *Formatter) extractMessage(e map[string]interface{}) string {
	for _, k := range []string{"msg", "message"} {
		if v, ok := e[k]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// levelColor returns the color associated with the given log level.
// Falls back to defaultColor for unrecognised levels.
func levelColor(level string) *color.Color {
	if c, ok := levelColors[level]; ok {
		return c
	}
	return defaultColor
}
