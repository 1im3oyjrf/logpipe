// Package checkpoint provides persistent offset tracking for log sources,
// allowing logpipe to resume tailing from where it left off after a restart.
package checkpoint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Store holds file read offsets keyed by source path.
type Store struct {
	mu      sync.Mutex
	path    string
	offsets map[string]int64
}

// New loads an existing checkpoint file from disk, or returns an empty Store
// if the file does not yet exist.
func New(path string) (*Store, error) {
	s := &Store{
		path:    path,
		offsets: make(map[string]int64),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &s.offsets); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns the last saved offset for the given source, or 0 if unknown.
func (s *Store) Get(source string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.offsets[source]
}

// Set updates the in-memory offset for the given source.
func (s *Store) Set(source string, offset int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offsets[source] = offset
}

// Flush writes the current offsets atomically to disk.
func (s *Store) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(s.offsets, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
