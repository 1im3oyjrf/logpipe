package output

import (
	"bytes"
	"strings"
	"testing"
)

func newNoColorFormatter(source string) (*Formatter, *bytes.Buffer) {
	var buf bytes.Buffer
	f := New(&buf, Options{NoColor: true, Source: source})
	return f, &buf
}

func TestFormatter_BasicEntry(t *testing.T) {
	f, buf := newNoColorFormatter("")
	f.Write(map[string]interface{}{
		"level": "info",
		"msg":   "server started",
		"time":  "2024-01-01T00:00:00Z",
	})
	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO in output, got: %s", out)
	}
	if !strings.Contains(out, "server started") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestFormatter_ErrorLevel(t *testing.T) {
	f, buf := newNoColorFormatter("")
	f.Write(map[string]interface{}{
		"level": "error",
		"msg":   "something broke",
	})
	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", out)
	}
}

func TestFormatter_ExtraFields(t *testing.T) {
	f, buf := newNoColorFormatter("")
	f.Write(map[string]interface{}{
		"level":  "debug",
		"msg":    "query executed",
		"time":   "2024-06-01T12:00:00Z",
		"query":  "SELECT 1",
		"dur_ms": 42,
	})
	out := buf.String()
	if !strings.Contains(out, "query=") {
		t.Errorf("expected extra field 'query' in output, got: %s", out)
	}
	if !strings.Contains(out, "dur_ms=") {
		t.Errorf("expected extra field 'dur_ms' in output, got: %s", out)
	}
}

func TestFormatter_SourceLabel(t *testing.T) {
	f, buf := newNoColorFormatter("app.log")
	f.Write(map[string]interface{}{
		"level": "warn",
		"msg":   "low memory",
	})
	out := buf.String()
	if !strings.Contains(out, "app.log") {
		t.Errorf("expected source label in output, got: %s", out)
	}
}

func TestFormatter_UnixTimestamp(t *testing.T) {
	f, buf := newNoColorFormatter("")
	f.Write(map[string]interface{}{
		"level": "info",
		"msg":   "tick",
		"ts":    float64(0), // Unix epoch
	})
	out := buf.String()
	if !strings.Contains(out, "1970") {
		t.Errorf("expected formatted unix timestamp in output, got: %s", out)
	}
}

func TestFormatter_MissingLevelDefaultsToInfo(t *testing.T) {
	f, buf := newNoColorFormatter("")
	f.Write(map[string]interface{}{
		"message": "no level field",
	})
	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected default INFO level, got: %s", out)
	}
}
