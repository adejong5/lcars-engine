package ha

import (
	"sort"
	"sync"
)

// Store is a concurrent-safe in-memory cache of entity states, keyed by
// entity_id. It is the shared sink that the mock provider and (Phase 3) the
// live HA client write into, and that the HTTP handlers read from.
type Store struct {
	mu sync.RWMutex
	m  map[string]State
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{m: make(map[string]State)}
}

// Get returns the state for an entity and whether it is present.
func (s *Store) Get(id string) (State, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.m[id]
	return st, ok
}

// Has reports whether an entity is present.
func (s *Store) Has(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[id]
	return ok
}

// Set inserts or replaces a state (keyed by its EntityID).
func (s *Store) Set(st State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[st.EntityID] = st
}

// All returns every state as a snapshot, sorted by entity_id for stable output.
func (s *Store) All() []State {
	s.mu.RLock()
	out := make([]State, 0, len(s.m))
	for _, st := range s.m {
		out = append(out, st)
	}
	s.mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].EntityID < out[j].EntityID })
	return out
}
