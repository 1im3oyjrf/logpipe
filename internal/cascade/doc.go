// Package cascade provides a rule-based field injection processor.
//
// A Cascade evaluates an ordered list of rules against each log entry.
// Each rule specifies a source field and expected value; when the entry's
// field matches (case-insensitively), a target field is set to a configured
// value.
//
// By default StopOnFirst is true, meaning only the first matching rule is
// applied. Setting StopOnFirst to false causes all matching rules to be
// evaluated and their target fields to be injected.
//
// Example usage:
//
//	c, err := cascade.New(cascade.Config{
//	    Rules: []cascade.Rule{
//	        {Field: "level", Value: "error", Target: "priority", Set: "high"},
//	        {Field: "level", Value: "warn",  Target: "priority", Set: "medium"},
//	    },
//	    StopOnFirst: true,
//	})
//	out := c.Apply(entry)
package cascade
