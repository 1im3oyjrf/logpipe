// Package split provides the Splitter processor, which expands a single log
// entry containing an array-valued field into multiple entries — one per
// element of the array.
//
// # Usage
//
//	s, err := split.New(split.Config{
//		Field:       "tags",   // array field to expand
//		TargetField: "tag",    // optional: field name on each child entry
//	})
//	if err != nil { … }
//
//	children := s.Apply(entry)  // returns []Entry
//
// If the named field is absent or not a []any slice the original entry is
// returned unchanged inside a one-element slice so callers can always range
// over the result unconditionally.
package split
