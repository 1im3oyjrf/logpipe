// Package pivot provides a log entry transformer that promotes a dynamic
// key/value pair into a top-level field.
//
// Given an entry such as:
//
//	{"key": "region", "value": "us-east-1", "level": "info"}
//
// Applying the Pivoter yields:
//
//	{"region": "us-east-1", "key": "region", "value": "us-east-1", "level": "info"}
//
// When DropSource is enabled the original KeyField and ValueField are removed
// from the output, leaving only the newly promoted field.
package pivot
