// Package promote provides a log entry transformer that lifts a nested or
// embedded field to the top level of the entry.
//
// When the target field holds a map value, all of its key-value pairs are
// merged directly into the parent entry. When the value is a scalar (string,
// number, bool, etc.) it is assigned to a configurable top-level key that
// defaults to the leaf segment of the source field path.
//
// The original field can optionally be retained by setting DropSource to false
// in the Config. Case-insensitive field matching is also supported.
//
// Example usage:
//
//	p, err := promote.New(promote.Config{
//		Field:      "metadata",
//		DropSource: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	out := p.Apply(entry)
package promote
