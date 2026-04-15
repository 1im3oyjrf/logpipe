// Package routing provides log entry routing based on level or field rules.
// Entries can be dispatched to named output channels depending on configured
// routing rules, enabling fan-out to multiple destinations.
package routing

import (
	"fmt"
	"strings"
	"sync"

	"github.com/user/logpipe/internal/reader"
)

// Rule defines a single routing condition and its target channel name.
type Rule struct {
	// Field is the JSON field to inspect (e.g. "level", "service").
	Field string
	// Value is the expected value (case-insensitive match).
	Value string
	// Target is the named destination channel for matched entries.
	Target string
}

// Router dispatches log entries to named channels based on routing rules.
// Entries that match no rule are sent to the "default" channel.
type Router struct {
	mu       sync.RWMutex
	rules    []Rule
	channels map[string]chan reader.Entry
}

// New creates a Router with the given rules and per-channel buffer size.
func New(rules []Rule, bufSize int) *Router {
	if bufSize <= 0 {
		bufSize = 64
	}
	channels := make(map[string]chan reader.Entry)
	channels["default"] = make(chan reader.Entry, bufSize)
	for _, r := range rules {
		if _, ok := channels[r.Target]; !ok {
			channels[r.Target] = make(chan reader.Entry, bufSize)
		}
	}
	return &Router{
		rules:    rules,
		channels: channels,
	}
}

// Channel returns the named read-only output channel, or nil if not found.
func (r *Router) Channel(name string) <-chan reader.Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, ok := r.channels[name]
	if !ok {
		return nil
	}
	return ch
}

// Targets returns all registered channel names.
func (r *Router) Targets() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.channels))
	for name := range r.channels {
		names = append(names, name)
	}
	return names
}

// Dispatch sends the entry to the first matching rule's channel, or "default".
func (r *Router) Dispatch(entry reader.Entry) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rule := range r.rules {
		val, ok := entry.Fields[rule.Field]
		if !ok {
			continue
		}
		if strings.EqualFold(fmt.Sprintf("%v", val), rule.Value) {
			select {
			case r.channels[rule.Target] <- entry:
			default:
			}
			return
		}
	}
	select {
	case r.channels["default"] <- entry:
	default:
	}
}

// Close closes all output channels.
func (r *Router) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ch := range r.channels {
		close(ch)
	}
}
