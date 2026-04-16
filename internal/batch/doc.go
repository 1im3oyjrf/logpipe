// Package batch provides a time-and-size-based batcher for log entries.
//
// A Batcher reads from an input channel of reader.Entry values and groups
// them into slices (batches). A batch is emitted when either:
//
//   - the number of buffered entries reaches MaxSize, or
//   - MaxWait duration elapses since the last flush.
//
// This is useful for downstream consumers that benefit from processing
// multiple entries at once, such as bulk writers or aggregators.
package batch
