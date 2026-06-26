package view

import "testing"

func TestRenderTileFormatted(t *testing.T) {
	sp := TileSpec{Label: "CPU", Unit: "%", Format: Round, Class: Hot(85)}

	tile := RenderTile(sp, "90.4", true)
	if !tile.OK || tile.Value != "90%" {
		t.Fatalf("value = %q ok=%v, want 90%%", tile.Value, tile.OK)
	}
	if tile.Class != AlertClass {
		t.Fatalf("class = %q, want alert", tile.Class)
	}
	if tile.Label != "CPU" {
		t.Fatalf("label = %q", tile.Label)
	}

	calm := RenderTile(sp, "12", true)
	if calm.Value != "12%" || calm.Class != "" {
		t.Fatalf("calm tile = %+v", calm)
	}
}

func TestRenderTileUnavailable(t *testing.T) {
	sp := TileSpec{Label: "CPU", Unit: "%", Format: Round}
	for _, raw := range []string{"", "unknown", "unavailable"} {
		tile := RenderTile(sp, raw, true)
		if tile.OK || tile.Value != Placeholder {
			t.Fatalf("RenderTile(%q) = %+v, want placeholder", raw, tile)
		}
	}
	// unknown entity (not in store)
	if tile := RenderTile(sp, "90", false); tile.OK || tile.Value != Placeholder {
		t.Fatalf("unknown entity = %+v, want placeholder", tile)
	}
}

func TestRenderTilePassthrough(t *testing.T) {
	// no Format → raw value (e.g. a switch state) shown as-is
	tile := RenderTile(TileSpec{Label: "Spare"}, "on", true)
	if tile.Value != "on" || !tile.OK {
		t.Fatalf("passthrough = %+v", tile)
	}
}
