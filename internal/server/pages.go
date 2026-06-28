package server

import (
	"net/http"

	"github.com/adejong5/lcars-engine/internal/view"
)

// cell is one live readout on a page: its entity id (also the element id and
// SSE event name) and how to present it. opsCells defines the Operations page.
type cell struct {
	ID   string // entity id
	Spec view.TileSpec
}

var opsCells = []cell{
	{ID: "switch.kitchen_spare", Spec: view.TileSpec{Label: "Kitchen Spare"}},
	{ID: "climate.kitchen", Spec: view.TileSpec{Label: "Kitchen HVAC"}},
	{ID: "sensor.kitchen_current_temperature", Spec: view.TileSpec{Label: "Kitchen", Unit: "°F", Format: view.Fixed(1)}},
	{ID: "sensor.spare_temperature_temperature", Spec: view.TileSpec{Label: "Spare", Unit: "°F", Format: view.Fixed(1)}},
}

func cellByID(id string) (cell, bool) {
	for _, c := range opsCells {
		if c.ID == id {
			return c, true
		}
	}
	return cell{}, false
}

// cellView is a cell paired with its current rendered tile, for the templates.
type cellView struct {
	ID   string
	Tile view.Tile
}

type pageData struct {
	Title  string
	Banner string
	Base   string // ingress path prefix ("" when standalone)
	Cells  []cellView
}

func (s *Server) renderCell(c cell) cellView {
	st, ok := s.src.State(c.ID)
	return cellView{ID: c.ID, Tile: view.RenderTile(c.Spec, st.State, ok)}
}

// handleOps renders the Operations page from current entity state.
func (s *Server) handleOps(w http.ResponseWriter, r *http.Request) {
	cells := make([]cellView, len(opsCells))
	for i, c := range opsCells {
		cells[i] = s.renderCell(c)
	}
	data := pageData{Title: "Operations", Banner: "HIGHBROOK • OPERATIONS", Base: BasePath(r), Cells: cells}
	if err := s.render.Page(w, data); err != nil {
		s.log.Error("render ops", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// handleCell returns one live cell's fragment (the fallback-poll target).
func (s *Server) handleCell(w http.ResponseWriter, r *http.Request) {
	c, ok := cellByID(r.PathValue("id"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	if err := s.render.Cell(w, s.renderCell(c).Tile); err != nil {
		s.log.Error("render cell", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}
