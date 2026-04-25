// Package sieve implements a probabilistic duplicate-suppression filter for
// structured log entries.
//
// A [Sieve] maintains a fixed-size boolean bit array. Each incoming entry has
// a chosen field value hashed to a slot index. If the slot is unoccupied the
// entry is forwarded and the slot is marked; if the slot is already occupied
// the entry is dropped as a probable duplicate.
//
// Because multiple distinct values can hash to the same slot, false positives
// (legitimate entries being dropped) are possible. The probability decreases
// as [Config.Slots] increases. There are no false negatives: an entry that
// has genuinely been seen before will always be dropped.
//
// Call [Sieve.Reset] to clear all slots and begin a new generation, for
// example at the start of each log-rotation window.
package sieve
