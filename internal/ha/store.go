package ha

import (
	"sort"
	"sync"
)

// Store is a concurrent-safe in-memory cache of entity states, keyed by
// entity_id. It is the shared sink that the mock provider and the live HA
// client write into, and that the HTTP handlers read from. It also broadcasts
// entity changes to subscribers (the SSE endpoint).
type Store struct {
	mu      sync.RWMutex
	m       map[string]State
	subs    map[int]chan string // subscriber id → channel of changed entity_ids
	nextSub int
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

// Set inserts or replaces a state (keyed by its EntityID) and, if the state
// value changed, notifies subscribers with the entity_id.
func (s *Store) Set(st State) {
	s.mu.Lock()
	old, existed := s.m[st.EntityID]
	s.m[st.EntityID] = st
	changed := !existed || old.State != st.State
	var targets []chan string
	if changed && len(s.subs) > 0 {
		targets = make([]chan string, 0, len(s.subs))
		for _, ch := range s.subs {
			targets = append(targets, ch)
		}
	}
	s.mu.Unlock()

	for _, ch := range targets {
		select {
		case ch <- st.EntityID:
		default: // subscriber lagging; it reconciles via the fallback poll
		}
	}
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

// Subscribe returns a channel that receives the entity_id of each changed
// entity, plus a cancel func to unsubscribe. The channel is buffered and sends
// are dropped if a subscriber falls behind (clients reconcile via the fallback
// poll). The channel is never closed, to avoid racing a concurrent Set.
func (s *Store) Subscribe() (<-chan string, func()) {
	ch := make(chan string, 64)
	s.mu.Lock()
	if s.subs == nil {
		s.subs = make(map[int]chan string)
	}
	id := s.nextSub
	s.nextSub++
	s.subs[id] = ch
	s.mu.Unlock()

	return ch, func() {
		s.mu.Lock()
		delete(s.subs, id)
		s.mu.Unlock()
	}
}
