package ha

// Source is the read interface the HTTP layer depends on. The mock provider
// (Phase 2) and the live HA client (Phase 3) both implement it, so handlers are
// agnostic to where state comes from. CallService is added in Phase 4.
type Source interface {
	// State returns the current state for an entity. The mock provider
	// fabricates unknown entities on demand (ok always true); the live client
	// returns ok=false for entities HA has not reported.
	State(id string) (State, bool)

	// All returns every known state, sorted by entity_id.
	All() []State

	// CallService invokes a Home Assistant service (e.g. switch.turn_on). The
	// live client sends it over the WebSocket; the mock applies it optimistically
	// to its own store. data and target may be nil.
	CallService(domain, service string, data, target map[string]any) error
}
