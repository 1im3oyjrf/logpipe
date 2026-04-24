// Package extract implements a log-entry processor that promotes nested
// fields to the top level using dot-separated key paths.
//
// Given an entry such as:
//
//	{"metadata": {"request_id": "abc", "user_id": 42}}
//
// Configuring the path "metadata.request_id" produces:
//
//	{"metadata": {"request_id": "abc", "user_id": 42}, "metadata.request_id": "abc"}
//
// When DropSource is enabled the leaf key is removed from the nested map
// after extraction, leaving the parent map intact for any remaining keys.
//
// Field matching is case-sensitive by default; set CaseInsensitive to true
// to match regardless of capitalisation.
package extract
