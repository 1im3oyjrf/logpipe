package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// LogEntry represents a single parsed JSON log line.
type LogEntry struct {
	Source    string
	Raw       string
	Fields    map[string]interface{}
	Timestamp time.Time
}

// JSONReader reads lines from an io.Reader and parses them as JSON log entries.
type JSONReader struct {
	source string
	reader *bufio.Reader
}

// NewJSONReader creates a new JSONReader for the given source name and reader.
func NewJSONReader(source string, r io.Reader) *JSONReader {
	return &JSONReader{
		source: source,
		reader: bufio.NewReader(r),
	}
}

// ReadAll reads all available log entries from the reader and sends them to the provided channel.
// It stops when EOF is reached or an unrecoverable error occurs.
func (j *JSONReader) ReadAll(out chan<- LogEntry) error {
	for {
		line, err := j.reader.ReadString('\n')
		if len(line) > 0 {
			entry, parseErr := j.parseLine(line)
			if parseErr == nil {
				out <- entry
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reader error for source %q: %w", j.source, err)
		}
	}
}

func (j *JSONReader) parseLine(line string) (LogEntry, error) {
	fields := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return LogEntry{}, fmt.Errorf("invalid JSON: %w", err)
	}

	entry := LogEntry{
		Source: j.source,
		Raw:    line,
		Fields: fields,
	}

	// Attempt to extract a timestamp from common field names.
	for _, key := range []string{"time", "timestamp", "ts", "@timestamp"} {
		if v, ok := fields[key]; ok {
			if s, ok := v.(string); ok {
				if t, err := time.Parse(time.RFC3339, s); err == nil {
					entry.Timestamp = t
					break
				}
			}
		}
	}

	return entry, nil
}
