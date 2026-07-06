package server

import (
	"fmt"
	"net/http"
	"time"
)

// handleSSE streams per-element fragment updates to the browser over Server-
// Sent Events. Event names are element ids (an entity can appear on several
// pages with different presentations, so events are keyed by element, not
// entity). The htmx sse extension swaps each event into the element whose
// sse-swap matches; events for other pages' elements simply match nothing.
// This is the primary live-update path; the page's infrequent hx-get poll is
// the fallback.
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	rc := http.NewResponseController(w)
	// Long-lived stream: clear the server's WriteTimeout for this connection.
	_ = rc.SetWriteDeadline(time.Time{})

	changes, cancel := s.src.Subscribe()
	defer cancel()

	keepalive := time.NewTicker(25 * time.Second)
	defer keepalive.Stop()
	ctx := r.Context()

	emit := func(event, frag string) bool {
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, frag)
		return rc.Flush() == nil
	}

	fmt.Fprint(w, ": connected\n\n")
	_ = rc.Flush()

	for {
		select {
		case <-ctx.Done():
			return
		case <-keepalive.C:
			fmt.Fprint(w, ": keepalive\n\n")
			if rc.Flush() != nil {
				return
			}
		case id := <-changes:
			for _, le := range elemsByEntity[id] {
				ev := s.renderElem(le)
				frag, err := s.render.FragmentHTML(string(ev.Kind), ev.fragData())
				if err != nil {
					continue
				}
				if !emit(le.ID, frag) {
					return
				}
			}
			// A cascade re-renders whole when any constituent entity changes.
			for _, cid := range cascadesByEnt[id] {
				frag, err := s.render.FragmentHTML("cascade-cols", s.renderCascade(cid))
				if err != nil {
					continue
				}
				if !emit(cid, frag) {
					return
				}
			}
			// Controls can change from outside (wall switch, HA app): refresh
			// the button. A control may render as a ctl-nav toggle or a section
			// pill-toggle; emit both variants (the absent one matches nothing).
			if c, ok := controlsByEnt[id]; ok {
				tv := s.toggleView(c)
				if frag, err := s.render.FragmentHTML("toggle", tv); err == nil {
					if !emit("ctl-"+c.ID, frag) {
						return
					}
				}
				if frag, err := s.render.FragmentHTML("pill-toggle", tv); err == nil {
					if !emit("pill-"+c.ID, frag) {
						return
					}
				}
			}
		}
	}
}
