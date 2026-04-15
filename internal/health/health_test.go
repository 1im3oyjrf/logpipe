package health_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/health"
)

func TestNew_NoSources_StatusUnknown(t *testing.T) {
	c := health.New()
	if got := c.Overall(); got != health.StatusUnknown {
		t.Errorf("expected StatusUnknown, got %s", got)
	}
}

func TestRegister_SingleSource_StatusOK(t *testing.T) {
	c := health.New()
	c.Register("stdin")
	if got := c.Overall(); got != health.StatusOK {
		t.Errorf("expected StatusOK after register, got %s", got)
	}
}

func TestSetError_MarksDegraded(t *testing.T) {
	c := health.New()
	c.Register("file1.log")
	c.Register("file2.log")
	c.SetError("file1.log", errors.New("file not found"))
	if got := c.Overall(); got != health.StatusDegraded {
		t.Errorf("expected StatusDegraded, got %s", got)
	}
}

func TestSetHealthy_ClearsError(t *testing.T) {
	c := health.New()
	c.Register("app.log")
	c.SetError("app.log", errors.New("read error"))
	c.SetHealthy("app.log")
	if got := c.Overall(); got != health.StatusOK {
		t.Errorf("expected StatusOK after recovery, got %s", got)
	}
}

func TestSetError_UnknownSource_NoOp(t *testing.T) {
	c := health.New()
	c.Register("known.log")
	// Should not panic for unregistered source
	c.SetError("unknown.log", errors.New("oops"))
	if got := c.Overall(); got != health.StatusOK {
		t.Errorf("expected StatusOK, got %s", got)
	}
}

func TestReport_ContainsSourceNames(t *testing.T) {
	c := health.New()
	c.Register("alpha.log")
	c.Register("beta.log")
	c.SetError("beta.log", errors.New("timeout"))

	var buf bytes.Buffer
	c.Report(&buf)
	out := buf.String()

	if !strings.Contains(out, "alpha.log") {
		t.Errorf("report missing alpha.log: %s", out)
	}
	if !strings.Contains(out, "beta.log") {
		t.Errorf("report missing beta.log: %s", out)
	}
	if !strings.Contains(out, "timeout") {
		t.Errorf("report missing error message: %s", out)
	}
	if !strings.Contains(out, "DEGRADED") {
		t.Errorf("report missing DEGRADED status: %s", out)
	}
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		s    health.Status
		want string
	}{
		{health.StatusOK, "OK"},
		{health.StatusDegraded, "DEGRADED"},
		{health.StatusUnknown, "UNKNOWN"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Status.String() = %q, want %q", got, tc.want)
		}
	}
}
