// Package label implements rule-based labelling for structured log entries.
//
// A Labeller holds a set of Rules. Each Rule specifies a field name, a
// substring value to match against (case-insensitively), and a map of
// key/value labels to inject into the entry when the rule fires.
//
// Multiple rules may match a single entry; all matching labels are applied.
// The original entry map is never modified — Apply always returns a shallow
// copy with the additional label fields merged in.
//
// Example:
//
//	rules := []label.Rule{
//		{Field: "level", Value: "error", Labels: map[string]string{"team": "oncall"}},
//	}
//	l := label.New(rules)
//	labelled := l.Apply(entry)
package label
