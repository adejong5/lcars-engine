package server

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"

	"github.com/adejong5/lcars-engine/internal/view"
)

// frameData is the common data every page needs (used by frame.html): titles,
// the left-column primary navigation panels, and the page's functional
// controls (the button array top-right under the banner).
type frameData struct {
	Title    string
	Banner   string
	Site     string // installation name in the 01 panel (siteLabel)
	Base     string // ingress path prefix ("" when standalone)
	Live     bool   // open the SSE stream (dashboard pages)
	Panels   []panelLink
	Controls []serverToggle
}

// panelLink is one left-column nav panel. N is the theme panel number (3-9);
// filler panels have no Href and render inert.
type panelLink struct {
	N    int
	Code string // panel text, e.g. "04-THERMAL"
	Href string
	Here bool // this is the current page
}

// framePanels builds the left-column nav: one panel per dashboard page (03-06)
// plus inert fillers (07-09), marking the current page.
func framePanels(current string) []panelLink {
	var ps []panelLink
	n := 3
	for _, p := range pages {
		ps = append(ps, panelLink{
			N:    n,
			Code: fmt.Sprintf("%02d-%s", n, strings.ToUpper(p.Title)),
			Href: p.Slug,
			Here: p.Slug == current,
		})
		n++
	}
	for ; n <= 9; n++ {
		ps = append(ps, panelLink{N: n, Code: fmt.Sprintf("%02d", n)})
	}
	return ps
}

// ── Dashboard pages ────────────────────────────────────────────
//
// The four pages are data: each is columns of bezel-framed sections, each
// section a list of entity-backed elements presented as tiles, segmented
// meters, bar rows, or a data cascade. Page-level Controls are the functional
// buttons in the top-right array. Editing a page means editing these tables —
// handlers, SSE, and the fallback poll all derive from them.

// elemKind selects the presentation (and the fragment template) for an element.
type elemKind string

const (
	kindTile    elemKind = "cell"
	kindBar     elemKind = "bar-row"
	kindMeter   elemKind = "meter-block"
	kindReadout elemKind = "readout"
)

// elem is one live readout: an entity and how to present it. Exactly one of
// Tile/Bar/Meter/Readout is used, per Kind. (Cascade sections use Tile specs
// for their cell formatting.)
type elem struct {
	Entity  string
	Kind    elemKind
	Tile    view.TileSpec
	Bar     view.BarSpec
	Meter   view.MeterSpec
	Readout view.ReadoutSpec
}

func tile(entity string, sp view.TileSpec) elem {
	return elem{Entity: entity, Kind: kindTile, Tile: sp}
}
func bar(entity string, sp view.BarSpec) elem { return elem{Entity: entity, Kind: kindBar, Bar: sp} }
func meter(entity string, sp view.MeterSpec) elem {
	return elem{Entity: entity, Kind: kindMeter, Meter: sp}
}
func readout(entity string, sp view.ReadoutSpec) elem {
	return elem{Entity: entity, Kind: kindReadout, Readout: sp}
}

// section is one group on a page. A read-only readout is bar-headed and lists
// Elems (Kind "cells"/"bars"/"meters"/"cascade"). A control cluster is a bezel:
// its Heading is the top-piece header and its Groups are the content
// sub-sections — one per bezel piece, the last in the bottom (base) piece and
// any earlier ones in middle pieces. PillCols sets the pill grid (1 = wide
// rows, 3 = dense).
type section struct {
	Heading  string
	Kind     string
	PillCols int
	Color    string     // bezel category colour ("" = theme default)
	Elems    []elem     // readout sections
	Groups   [][]string // control-bezel content groups (control IDs per piece)
}

type pageDef struct {
	Slug, Title, Banner string
	Controls            []string    // control IDs: the page's functional buttons
	Cols                [][]section // columns of primary (bezel) sections
}

var (
	motion   = view.Map(map[string]string{"on": "MOTION", "off": "CLEAR"})
	opened   = view.Map(map[string]string{"on": "OPEN", "off": "CLOSED"})
	upToDate = view.Map(map[string]string{"on": "AVAILABLE", "off": "CURRENT"})
	fix1     = view.Fixed(1)
	fix0     = view.Fixed(0)
	stops    = []string{"var(--ice)", "var(--gold)", "var(--mars)"}
	// alarmed flags an open door/window in red (binary_sensor "on").
	alarmed = view.Map2(map[string]string{"on": view.AlertClass})
	hotCPU  = view.Hot(160) // server temps (°F): red above 160
	busy    = view.Hot(85)  // usage % : red above 85
)

// roomTemp presents a room temperature as a segmented meter.
func roomTemp(label string) view.MeterSpec {
	return view.MeterSpec{Label: label, Min: 50, Max: 90, Unit: "°F", Decimals: 1, Segments: 16, Stops: stops}
}

// pctMeter presents a 0-100% value as a segmented meter.
func pctMeter(label string) view.MeterSpec {
	return view.MeterSpec{Label: label, Min: 0, Max: 100, Unit: "%", Decimals: 0, Segments: 16, Stops: stops}
}

// battBar presents a battery level as a bar row.
func battBar(label string) view.BarSpec {
	return view.BarSpec{Label: label, Max: 100, Unit: "%", Decimals: 0, Color: "var(--ice)"}
}

// pages and controls are the site's dashboard definitions. They are filled in
// by pages_local.go (gitignored): entity ids, room names, and page layouts are
// personal to the household and stay out of the repo. Without a local file the
// server builds with no dashboard pages — the component demo at "/" still
// serves. See ENGINE.md for how to compose pages.
var pages []pageDef

// siteLabel names the installation in the frame (the 01 panel); local
// definitions override it.
var siteLabel = "LCARS"

// cascadeRows is how many cells stack in one data-cascade column.
const cascadeRows = 3

// liveElem is an element with its page-scoped identity: the HTML element id,
// SSE event name, and fallback-poll path segment are all ID.
type liveElem struct {
	ID   string
	Elem elem
}

// Derived indexes: element by id (poll), elements by entity (SSE fan-out),
// cascade sections by id and by entity, controls by entity, pages by slug.
var (
	elemByID      map[string]liveElem
	elemsByEntity map[string][]liveElem
	cascadeByID   map[string][]elem
	cascadesByEnt map[string][]string
	controlsByEnt map[string]control
	pageBySlug    map[string]pageDef
)

func init() { registerPages() }

// registerPages (re)builds the derived indexes from the pages/controls tables.
// Idempotent: pages_local.go's init calls it again after filling the tables
// (package init order between the two files is not relied upon).
func registerPages() {
	elemByID = map[string]liveElem{}
	elemsByEntity = map[string][]liveElem{}
	cascadeByID = map[string][]elem{}
	cascadesByEnt = map[string][]string{}
	controlsByEnt = map[string]control{}
	pageBySlug = map[string]pageDef{}
	for _, p := range pages {
		pageBySlug[p.Slug] = p
		n, nc := 0, 0
		for _, col := range p.Cols {
			for _, sec := range col {
				if sec.Kind == "cascade" {
					id := fmt.Sprintf("%s-casc%d", p.Slug, nc)
					nc++
					cascadeByID[id] = sec.Elems
					for _, e := range sec.Elems {
						cascadesByEnt[e.Entity] = append(cascadesByEnt[e.Entity], id)
					}
					continue
				}
				for _, e := range sec.Elems {
					le := liveElem{ID: fmt.Sprintf("%s-%d", p.Slug, n), Elem: e}
					elemByID[le.ID] = le
					elemsByEntity[e.Entity] = append(elemsByEntity[e.Entity], le)
					n++
				}
			}
		}
	}
	for _, c := range controls {
		controlsByEnt[c.Entity] = c
	}
}

// elemView carries one rendered element to the dash template; Kind picks the
// block, and the matching field is populated.
type elemView struct {
	ID      string
	Kind    elemKind
	Tile    view.Tile
	Bar     view.Bar
	Meter   view.MeterBlock
	Readout view.Readout
}

func (s *Server) renderElem(le liveElem) elemView {
	st, ok := s.src.State(le.Elem.Entity)
	ev := elemView{ID: le.ID, Kind: le.Elem.Kind}
	switch le.Elem.Kind {
	case kindBar:
		ev.Bar = view.RenderBar(le.Elem.Bar, st.State, ok)
	case kindMeter:
		ev.Meter = view.RenderMeterBlock(le.ID, le.Elem.Meter, st.State, ok)
	case kindReadout:
		ev.Readout = view.RenderReadout(le.Elem.Readout, st.State, ok)
	default:
		ev.Tile = view.RenderTile(le.Elem.Tile, st.State, ok)
	}
	return ev
}

// fragData returns the template data matching an elemView's fragment kind.
func (ev elemView) fragData() any {
	switch ev.Kind {
	case kindBar:
		return ev.Bar
	case kindMeter:
		return ev.Meter
	case kindReadout:
		return ev.Readout
	default:
		return ev.Tile
	}
}

// renderCascade formats a cascade section's entities into data-cascade columns
// (cascadeRows cells per column, colour row classes cycled diagonally).
func (s *Server) renderCascade(id string) [][]cascadeCell {
	elems := cascadeByID[id]
	var cols [][]cascadeCell
	for i := 0; i < len(elems); i += cascadeRows {
		var col []cascadeCell
		for j := i; j < i+cascadeRows && j < len(elems); j++ {
			st, ok := s.src.State(elems[j].Entity)
			t := view.RenderTile(elems[j].Tile, st.State, ok)
			col = append(col, cascadeCell{Row: (i/cascadeRows+j-i)%7 + 1, Text: t.Value, Cls: t.Class})
		}
		cols = append(cols, col)
	}
	return cols
}

type sectionView struct {
	Heading  string
	Kind     string
	PillCols int
	Color    string
	ID       string // cascade sections: the swap target id
	Elems    []elemView
	Cascade  [][]cascadeCell
	// Control bezel: each group of toggles renders in one bezel piece. BotGroup
	// is the last group (bottom/base piece); MidGroups are the earlier ones
	// (middle pieces). A non-empty BotGroup marks the section as a bezel.
	MidGroups [][]serverToggle
	BotGroup  []serverToggle
}

type dashData struct {
	frameData
	Cols [][]sectionView
}

// handleDash renders any registered page; the slug is the request path.
func (s *Server) handleDash(w http.ResponseWriter, r *http.Request) {
	p, ok := pageBySlug[strings.TrimPrefix(r.URL.Path, "/")]
	if !ok {
		http.NotFound(w, r)
		return
	}
	data := dashData{frameData: frameData{
		Title: p.Title, Banner: p.Banner, Site: siteLabel, Base: BasePath(r), Live: true,
		Panels: framePanels(p.Slug), Controls: s.pageControls(p),
	}}
	n, nc := 0, 0
	for _, col := range p.Cols {
		var cv []sectionView
		for _, sec := range col {
			sv := sectionView{Heading: sec.Heading, Kind: sec.Kind, PillCols: sec.PillCols, Color: sec.Color}
			for gi, group := range sec.Groups {
				var ts []serverToggle
				for _, id := range group {
					if c, ok := controlByID(id); ok {
						ts = append(ts, s.toggleView(c))
					}
				}
				if gi == len(sec.Groups)-1 {
					sv.BotGroup = ts // last group → bottom (base) piece
				} else {
					sv.MidGroups = append(sv.MidGroups, ts)
				}
			}
			if sec.Kind == "cascade" {
				sv.ID = fmt.Sprintf("%s-casc%d", p.Slug, nc)
				nc++
				sv.Cascade = s.renderCascade(sv.ID)
			} else {
				for range sec.Elems {
					sv.Elems = append(sv.Elems, s.renderElem(elemByID[fmt.Sprintf("%s-%d", p.Slug, n)]))
					n++
				}
			}
			cv = append(cv, sv)
		}
		data.Cols = append(data.Cols, cv)
	}
	if err := s.render.Page(w, "dash", data); err != nil {
		s.log.Error("render page", "slug", p.Slug, "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// pageControls resolves a page's control IDs to live toggle views.
func (s *Server) pageControls(p pageDef) []serverToggle {
	var ts []serverToggle
	for _, id := range p.Controls {
		if c, ok := controlByID(id); ok {
			ts = append(ts, s.toggleView(c))
		}
	}
	return ts
}

// handleCell returns one live element's fragment (the fallback-poll target).
// Cascade section ids resolve to the whole cascade body.
func (s *Server) handleCell(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, ok := cascadeByID[id]; ok {
		if err := s.render.Fragment(w, "cascade-cols", s.renderCascade(id)); err != nil {
			s.log.Error("render cascade", "err", err)
			http.Error(w, "render error", http.StatusInternalServerError)
		}
		return
	}
	le, ok := elemByID[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	ev := s.renderElem(le)
	if err := s.render.Fragment(w, string(ev.Kind), ev.fragData()); err != nil {
		s.log.Error("render cell", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// ── Controls (hx-post actions) ─────────────────────────────────

// control is a stateful button bound to an actionable entity. SvcOn/SvcOff
// default to turn_on/turn_off; ActiveWhen defaults to isActive. Domains whose
// services aren't a turn_on/off pair (locks) override them — see ENGINE.md §5.
type control struct {
	ID         string // action id (URL) and element identity
	Label      string
	Domain     string
	Entity     string
	SvcOn      string
	SvcOff     string
	ActiveWhen func(string) bool
}

// lightCtl is the common case: a light as a toggle.
func lightCtl(id, label, entity string) control {
	return control{ID: id, Label: label, Domain: "light", Entity: entity}
}

var controls []control // site-local, set by pages_local.go

// services returns the control's on/off services (defaulted).
func (c control) services() (on, off string) {
	on, off = c.SvcOn, c.SvcOff
	if on == "" {
		on = "turn_on"
	}
	if off == "" {
		off = "turn_off"
	}
	return
}

// active applies the control's state predicate (defaulted to isActive).
func (c control) active(state string) bool {
	if c.ActiveWhen != nil {
		return c.ActiveWhen(state)
	}
	return isActive(state)
}

func controlByID(id string) (control, bool) {
	for _, c := range controls {
		if c.ID == id {
			return c, true
		}
	}
	return control{}, false
}

// isActive reports whether a toggleable entity is on. Switches and lights are
// literally "on"; climate reports its mode (heat_cool, off, …), so anything
// but off/unknown counts as running.
func isActive(state string) bool {
	switch state {
	case "off", "unknown", "unavailable", "":
		return false
	}
	return true
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
		on = c.active(st.State)
	}
	return serverToggle{ID: c.ID, Label: c.Label, State: onoff(on), On: on}
}

// handleAction toggles a control's entity and returns the updated button
// (?s=pill returns the pill-toggle variant). The service set is reflected
// optimistically; SSE/poll reconcile the real state.
func (s *Server) handleAction(w http.ResponseWriter, r *http.Request) {
	c, ok := controlByID(r.PathValue("id"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	svcOn, svcOff := c.services()
	st, _ := s.src.State(c.Entity)
	service := svcOn
	if c.active(st.State) {
		service = svcOff
	}
	if err := s.src.CallService(c.Domain, service, nil, map[string]any{"entity_id": c.Entity}); err != nil {
		s.log.Error("call service", "entity", c.Entity, "err", err)
		http.Error(w, "service call failed", http.StatusBadGateway)
		return
	}
	on := service == svcOn
	frag := "toggle"
	if r.URL.Query().Get("s") == "pill" {
		frag = "pill-toggle"
	}
	tv := serverToggle{ID: c.ID, Label: c.Label, State: onoff(on), On: on}
	if err := s.render.Fragment(w, frag, tv); err != nil {
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
	Cls  string // threshold class (e.g. view.AlertClass)
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
	Readouts  []view.Readout
	Toggle    serverToggle
	Meters    []view.Meter
	Cascade   [][]cascadeCell
	NavItems  []navItem
	Image     imageFrameData
	Gallery   galleryData
	Accordion accordionData
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := indexData{
		frameData: frameData{Title: "Components", Banner: "LCARS • COMPONENTS", Site: siteLabel, Base: BasePath(r), Panels: framePanels("")},
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
		Readouts: []view.Readout{
			view.RenderReadout(view.ReadoutSpec{Label: "Deck 2", Min: 50, Max: 90, Unit: "°F", Decimals: 0, Stops: stops}, "63", true),
			view.RenderReadout(view.ReadoutSpec{Label: "Core", Min: 0, Max: 100, Unit: "%", Decimals: 0, Stops: stops}, "88", true),
			view.RenderReadout(view.ReadoutSpec{Label: "Aux", Min: 0, Max: 100, Unit: "%", Decimals: 0, Stops: stops}, "unknown", true),
		},
		Meters: []view.Meter{
			view.RenderMeter("m1", 25, 0, 100, 16, 3, stops),
			view.RenderMeter("m2", 60, 0, 100, 16, 3, stops),
			view.RenderMeter("m3", 90, 0, 100, 12, 4, stops),
		},
		Cascade:  demoCascade(),
		NavItems: demoNav(),
		Image:    imageFrameData{Title: "Forward Sensor", Src: "static/placeholder.svg", Alt: "placeholder"},
		Gallery: galleryData{Thumbs: true, Images: []galleryImage{
			{Src: "static/placeholder.svg", Alt: "one", Caption: "Deck 1"},
			{Src: "static/placeholder.svg", Alt: "two", Caption: "Deck 2"},
		}},
		Accordion: accordionData{Title: "Diagnostics", Body: "Collapsible detail — secondary readouts tucked away until needed."},
	}
	// The live toggle demo needs a defined control; without local definitions
	// the section renders a hint instead.
	if len(controls) > 0 {
		data.Toggle = s.toggleView(controls[0])
	}
	if err := s.render.Page(w, "index", data); err != nil {
		s.log.Error("render index", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// demoNav links the demo's nav-button row to the registered pages (empty when
// no local pages are defined).
func demoNav() []navItem {
	var items []navItem
	for _, p := range pages {
		items = append(items, navItem{Label: p.Title, Href: p.Slug})
	}
	return items
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
