// Package ha holds the Home Assistant state model, an in-memory store, and the
// data sources that fill it: a mock provider (Phase 2) and, later, the live
// WebSocket client (Phase 3).
package ha

import "time"

// State is a single Home Assistant entity state, shaped to mirror the relevant
// fields of HA's REST/WebSocket state objects.
type State struct {
	EntityID    string         `json:"entity_id"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
	LastUpdated time.Time      `json:"last_updated"`
}
