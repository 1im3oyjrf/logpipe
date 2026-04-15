// Package health provides a simple health check mechanism for logpipe,
// reporting whether all configured sources are active and the pipeline
// is processing entries without errors.
package health

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Status represents the overall health state of the pipeline.
type Status int

const (
	StatusOK      Status = iota // All sources healthy
	StatusDegraded              // One or more sources have errors
	StatusUnknown               // Not yet evaluated
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusDegraded:
		return "DEGRADED"
	default:
		return "UNKNOWN"
	}
}

// SourceHealth holds the health state for a single named source.
type SourceHealth struct {
	Name      string
	Healthy   bool
	LastError error
	CheckedAt time.Time
}

// Checker tracks per-source health and reports overall pipeline status.
type Checker struct {
	mu      sync.RWMutex
	sources map[string]*SourceHealth
}

// New returns an initialised Checker with no sources registered.
func New() *Checker {
	return &Checker{
		sources: make(map[string]*SourceHealth),
	}
}

// Register adds a named source to the health tracker, initially healthy.
func (c *Checker) Register(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sources[name] = &SourceHealth{
		Name:      name,
		Healthy:   true,
		CheckedAt: time.Now(),
	}
}

// SetError marks a source as unhealthy with the given error.
func (c *Checker) SetError(name string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if sh, ok := c.sources[name]; ok {
		sh.Healthy = false
		sh.LastError = err
		sh.CheckedAt = time.Now()
	}
}

// SetHealthy clears any error and marks the source as healthy.
func (c *Checker) SetHealthy(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if sh, ok := c.sources[name]; ok {
		sh.Healthy = true
		sh.LastError = nil
		sh.CheckedAt = time.Now()
	}
}

// Overall returns StatusOK if all sources are healthy, StatusDegraded otherwise.
// Returns StatusUnknown when no sources are registered.
func (c *Checker) Overall() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.sources) == 0 {
		return StatusUnknown
	}
	for _, sh := range c.sources {
		if !sh.Healthy {
			return StatusDegraded
		}
	}
	return StatusOK
}

// Report writes a human-readable health summary to w.
func (c *Checker) Report(w io.Writer) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	fmt.Fprintf(w, "status: %s\n", c.Overall())
	for _, sh := range c.sources {
		if sh.Healthy {
			fmt.Fprintf(w, "  [OK]      %s\n", sh.Name)
		} else {
			fmt.Fprintf(w, "  [ERROR]   %s: %v\n", sh.Name, sh.LastError)
		}
	}
}
