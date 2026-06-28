package ha

import (
	"context"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
)

// MockSource fabricates plausible entity data for offline development: any
// entity read is lazily seeded, and a background loop drifts a subset of values
// on each tick. Ported from the original ha.svelte.js mock behavior.
type MockSource struct {
	store *Store
}

// NewMock returns a mock source with an empty store.
func NewMock() *MockSource {
	return &MockSource{store: NewStore()}
}

// State lazily fabricates an unknown entity, then returns it (ok always true).
func (m *MockSource) State(id string) (State, bool) {
	if st, ok := m.store.Get(id); ok {
		return st, true
	}
	st := State{
		EntityID:    id,
		State:       mockSeedValue(id),
		Attributes:  map[string]any{},
		LastUpdated: time.Now().UTC(),
	}
	m.store.Set(st)
	return st, true
}

// All returns every fabricated state, sorted by entity_id.
func (m *MockSource) All() []State { return m.store.All() }

// Subscribe forwards to the underlying store.
func (m *MockSource) Subscribe() (<-chan string, func()) { return m.store.Subscribe() }

// CallService applies common on/off-style services to the targeted entities so
// the mock UI reflects control actions offline. Unrecognized services are
// accepted as no-ops.
func (m *MockSource) CallService(domain, service string, data, target map[string]any) error {
	for _, id := range targetEntities(target) {
		st, _ := m.State(id) // ensure seeded
		switch service {
		case "turn_on":
			st.State = "on"
		case "turn_off":
			st.State = "off"
		case "toggle":
			if st.State == "on" {
				st.State = "off"
			} else {
				st.State = "on"
			}
		case "lock":
			st.State = "locked"
		case "unlock":
			st.State = "unlocked"
		case "open_cover":
			st.State = "open"
		case "close_cover":
			st.State = "closed"
		default:
			continue
		}
		st.LastUpdated = time.Now().UTC()
		m.store.Set(st)
	}
	return nil
}

// targetEntities extracts entity_id(s) from a service target, which HA allows
// to be a single string or a list.
func targetEntities(target map[string]any) []string {
	if target == nil {
		return nil
	}
	switch v := target["entity_id"].(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case []any:
		out := make([]string, 0, len(v))
		for _, e := range v {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// Run drifts a random subset (~35%) of known entities on each interval until
// the context is cancelled.
func (m *MockSource) Run(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			m.tick()
		}
	}
}

func (m *MockSource) tick() {
	for _, st := range m.store.All() {
		if rand.Float64() < 0.35 {
			st.State = mockDriftValue(st.EntityID, st.State)
			st.LastUpdated = time.Now().UTC()
			m.store.Set(st)
		}
	}
}

// mockSeedValue returns an initial value for a freshly seen entity: on/off for
// binary sensors, a whole number otherwise.
func mockSeedValue(id string) string {
	if strings.HasPrefix(id, "binary_sensor.") {
		if rand.IntN(2) == 0 {
			return "off"
		}
		return "on"
	}
	return strconv.Itoa(2 + rand.IntN(9998))
}

// mockDriftValue returns the next value for an entity: binary sensors flip
// occasionally; numeric values drift +/-5%; anything else passes through.
func mockDriftValue(id, prev string) string {
	if strings.HasPrefix(id, "binary_sensor.") {
		if rand.Float64() > 0.8 {
			if prev == "on" {
				return "off"
			}
			return "on"
		}
		return prev
	}
	base, err := strconv.ParseFloat(prev, 64)
	if err != nil {
		return prev
	}
	drifted := base * (1 + (rand.Float64()*0.1 - 0.05))
	return strconv.FormatFloat(roundTo(drifted, 1), 'f', -1, 64)
}

func roundTo(v float64, decimals int) float64 {
	p := math.Pow(10, float64(decimals))
	return math.Round(v*p) / p
}
