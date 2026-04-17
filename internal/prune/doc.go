// Package prune provides a Pruner that removes unwanted fields from structured
// log entries before they are forwarded downstream.
//
// Fields are matched case-insensitively so that callers do not need to
// normalise field names before configuring the pruner.
//
// Usage:
//
//	p := prune.New(prune.Config{
//		Fields: []string{"secret", "token", "password"},
//	})
//	clean := p.Apply(entry)
package prune
