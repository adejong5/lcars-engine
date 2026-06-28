package server

import (
	"math/rand/v2"
	"net/http"

	"github.com/adejong5/lcars-engine/internal/view"
)

// frameData is the common header data every page needs (used by frame.html).
type frameData struct {
	Title  string
	Banner string
	Base   string // ingress path prefix ("" when standalone)
}

// ── Operations dashboard ───────────────────────────────────────

// cell is one live readout on the ops page: its entity id (also the element id
// and SSE event name) and how to present it.
type cell struct {
	ID   string
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

type cellView struct {
	ID   string
	Tile view.Tile
}

func (s *Server) renderCell(c cell) cellView {
	st, ok := s.src.State(c.ID)
	return cellView{ID: c.ID, Tile: view.RenderTile(c.Spec, st.State, ok)}
}

type opsData struct {
	frameData
	Cells []cellView
}

func (s *Server) handleOps(w http.ResponseWriter, r *http.Request) {
	cells := make([]cellView, len(opsCells))
	for i, c := range opsCells {
		cells[i] = s.renderCell(c)
	}
	data := opsData{
		frameData: frameData{Title: "Operations", Banner: "HIGHBROOK • OPERATIONS", Base: BasePath(r)},
		Cells:     cells,
	}
	if err := s.render.Page(w, "ops", data); err != nil {
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
	if err := s.render.Fragment(w, "cell", s.renderCell(c).Tile); err != nil {
		s.log.Error("render cell", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// ── Controls (hx-post actions) ─────────────────────────────────

// control is a toggle button bound to an on/off entity.
type control struct {
	ID     string // action id (URL) and element identity
	Label  string
	Domain string
	Entity string
}

var controls = []control{
	{ID: "kitchen_spare", Label: "Kitchen Spare", Domain: "switch", Entity: "switch.kitchen_spare"},
}

func controlByID(id string) (control, bool) {
	for _, c := range controls {
		if c.ID == id {
			return c, true
		}
	}
	return control{}, false
}

type serverToggle struct {
	ID    string
	Label string
	State string
	On    bool
}

func (s *Server) toggleView(c control) serverToggle {
	on := false
	if st, ok := s.src.State(c.Entity); ok {
		on = st.State == "on"
	}
	return serverToggle{ID: c.ID, Label: c.Label, State: onoff(on), On: on}
}

// handleAction toggles a control's entity and returns the updated button. The
// service set is reflected optimistically; SSE/poll reconcile the real state.
func (s *Server) handleAction(w http.ResponseWriter, r *http.Request) {
	c, ok := controlByID(r.PathValue("id"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	st, _ := s.src.State(c.Entity)
	service := "turn_on"
	if st.State == "on" {
		service = "turn_off"
	}
	if err := s.src.CallService(c.Domain, service, nil, map[string]any{"entity_id": c.Entity}); err != nil {
		s.log.Error("call service", "entity", c.Entity, "err", err)
		http.Error(w, "service call failed", http.StatusBadGateway)
		return
	}
	on := service == "turn_on"
	tv := serverToggle{ID: c.ID, Label: c.Label, State: onoff(on), On: on}
	if err := s.render.Fragment(w, "toggle", tv); err != nil {
		s.log.Error("render toggle", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

func onoff(on bool) string {
	if on {
		return "ON"
	}
	return "OFF"
}

// ── Component demo (index) ─────────────────────────────────────

type cascadeCell struct {
	Row  int
	Text string
}

type navItem struct {
	Label string
	Href  string
}

type galleryImage struct {
	Src, Alt, Caption string
}

type galleryData struct {
	Thumbs bool
	Images []galleryImage
}

type imageFrameData struct {
	Title, Src, Alt string
}

type accordionData struct {
	Title, Body string
}

type indexData struct {
	frameData
	Gauges    []view.Gauge
	Bars      []view.Bar
	Toggle    serverToggle
	Meters    []view.Meter
	Cascade   [][]cascadeCell
	NavItems  []navItem
	Image     imageFrameData
	Gallery   galleryData
	Accordion accordionData
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	stops := []string{"var(--ice)", "var(--gold)", "var(--mars)"}
	data := indexData{
		frameData: frameData{Title: "Components", Banner: "LCARS • COMPONENTS", Base: BasePath(r)},
		Gauges: []view.Gauge{
			{Label: "Cold", Display: "58.0°F", Frac: 0.10, OK: true},
			{Label: "Comfort", Display: "71.0°F", Frac: 0.53, OK: true},
			{Label: "Hot", Display: "84.0°F", Frac: 0.97, OK: true},
		},
		Bars: []view.Bar{
			{Label: "Download", Display: "842 Mb/s", Frac: 0.84, Color: "var(--ice)"},
			{Label: "Upload", Display: "31 Mb/s", Frac: 0.31, Color: "var(--african-violet)"},
			{Label: "Ping", Display: "12 ms", Frac: 0.12, Color: "var(--gold)"},
		},
		Toggle: s.toggleView(controls[0]),
		Meters: []view.Meter{
			view.RenderMeter("m1", 25, 0, 100, 16, 3, stops),
			view.RenderMeter("m2", 60, 0, 100, 16, 3, stops),
			view.RenderMeter("m3", 90, 0, 100, 12, 4, stops),
		},
		Cascade:  demoCascade(),
		NavItems: []navItem{{Label: "Subnet", Href: "subnet"}, {Label: "Thermal", Href: "thermal"}, {Label: "Engr", Href: "engr"}},
		Image:    imageFrameData{Title: "Forward Sensor", Src: "static/placeholder.svg", Alt: "placeholder"},
		Gallery: galleryData{Thumbs: true, Images: []galleryImage{
			{Src: "static/placeholder.svg", Alt: "one", Caption: "Deck 1"},
			{Src: "static/placeholder.svg", Alt: "two", Caption: "Deck 2"},
		}},
		Accordion: accordionData{Title: "Diagnostics", Body: "Collapsible detail — secondary readouts tucked away until needed."},
	}
	if err := s.render.Page(w, "index", data); err != nil {
		s.log.Error("render index", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// demoCascade builds a sample number grid for the component demo.
func demoCascade() [][]cascadeCell {
	pattern := []int{1, 1, 2, 3, 3, 4, 5, 6, 7}
	cols := make([][]cascadeCell, 12)
	for c := range cols {
		col := make([]cascadeCell, len(pattern))
		for i, p := range pattern {
			col[i] = cascadeCell{Row: p, Text: digits(1 + rand.IntN(6))}
		}
		cols[c] = col
	}
	return cols
}

func digits(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = byte('0' + rand.IntN(10))
	}
	return string(s)
}
