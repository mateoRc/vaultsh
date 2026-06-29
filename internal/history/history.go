package history

import "sync"

type Store struct {
	mu      sync.Mutex
	limit   int
	entries []string
}

func New(limit int) *Store {
	return &Store{limit: limit}
}

func (s *Store) Add(line string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = append(s.entries, line)
	if len(s.entries) > s.limit {
		s.entries = s.entries[len(s.entries)-s.limit:]
	}
}

func (s *Store) Entries() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries := make([]string, len(s.entries))
	copy(entries, s.entries)
	return entries
}
