// Package tag provides a Tagger that evaluates a set of rules against log
// entries and injects a list of string tags into a configurable target field.
//
// Rules are evaluated in order; all matching rules contribute their tags so
// the final tag list may contain tags from multiple rules. The original entry
// is never mutated — Apply always returns a new map when tags are added.
//
// Example:
//
//	tr := tag.New(tag.Config{
//		TargetField: "tags",
//		Rules: []tag.Rule{
//			{Field: "level", Value: "error", Tags: []string{"alert"}},
//		},
//	})
//	out := tr.Apply(entry)
package tag
