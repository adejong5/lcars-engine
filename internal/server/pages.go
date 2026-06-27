package server

import (
	"net/http"

	"github.com/adejong5/lcars-engine/internal/view"
)

// pageData is the view model passed to the page templates. The frame uses
// Title/Banner/Base; each page adds its own fields.
type pageData struct {
	Title  string
	Banner string
	Base   string // ingress path prefix ("" when standalone)

	Temps   []view.Gauge
	Switch  view.Tile
	Climate view.Tile
}

// handleOps renders the Operations page from current entity state.
func (s *Server) handleOps(w http.ResponseWriter, r *http.Request) {
	if err := s.render.Page(w, s.buildOps(BasePath(r))); err != nil {
		s.log.Error("render ops", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

func (s *Server) buildOps(base string) pageData {
	tile := func(id string, sp view.TileSpec) view.Tile {
		st, ok := s.src.State(id)
		return view.RenderTile(sp, st.State, ok)
	}
	gauge := func(id string, sp view.GaugeSpec) view.Gauge {
		st, ok := s.src.State(id)
		return view.RenderGauge(sp, st.State, ok)
	}
	return pageData{
		Title:   "Operations",
		Banner:  "HIGHBROOK • OPERATIONS",
		Base:    base,
		Switch:  tile("switch.kitchen_spare", view.TileSpec{Label: "Kitchen Spare"}),
		Climate: tile("climate.kitchen", view.TileSpec{Label: "Kitchen HVAC"}),
		Temps: []view.Gauge{
			gauge("sensor.kitchen_current_temperature", view.GaugeSpec{Label: "Kitchen", Min: 55, Max: 85, Unit: "°F", Decimals: 1}),
			gauge("sensor.spare_temperature_temperature", view.GaugeSpec{Label: "Spare", Min: 55, Max: 85, Unit: "°F", Decimals: 1}),
		},
	}
}
