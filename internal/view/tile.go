package view

// Tile is a single labelled readout — a data-cascade cell or a status pill.
type Tile struct {
	Label string // static label text
	Value string // formatted value (+unit), or Placeholder when unavailable
	Class string // extra CSS classes (e.g. alert styling)
	OK    bool   // false when the state is missing/unavailable
}

// TileSpec describes how to present one entity as a Tile.
type TileSpec struct {
	Label  string
	Unit   string
	Format Formatter  // nil → raw value passthrough
	Class  Classifier // nil → no extra classes
}

// RenderTile builds a Tile from a raw state string. known is whether the entity
// exists in the store. Pure: callers pass the raw string, so this needs no HA
// access. Works for non-numeric states too (e.g. switch on/off, climate mode).
func RenderTile(sp TileSpec, raw string, known bool) Tile {
	if !known || unavailable(raw) {
		return Tile{Label: sp.Label, Value: Placeholder, OK: false}
	}
	val := raw
	if sp.Format != nil {
		val = sp.Format(raw)
	}
	val += sp.Unit

	var cls string
	if sp.Class != nil {
		cls = sp.Class(raw)
	}
	return Tile{Label: sp.Label, Value: val, Class: cls, OK: true}
}
