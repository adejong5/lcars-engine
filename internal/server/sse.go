package server

import (
	"fmt"
	"net/http"
	"time"
)

// handleSSE streams per-entity fragment updates to the browser over Server-Sent
// Events. The htmx sse extension swaps each event into the element whose
// sse-swap matches the event name (the entity id). This is the primary live-
// update path; the page's infrequent hx-get poll is the fallback.
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	rc := http.NewResponseController(w)
	// Long-lived stream: clear the server's WriteTimeout for this connection.
	_ = rc.SetWriteDeadline(time.Time{})

	// Only entities that appear as live cells produce fragments.
	live := make(map[string]cell, len(opsCells))
	for _, c := range opsCells {
		live[c.ID] = c
	}

	changes, cancel := s.src.Subscribe()
	defer cancel()

	keepalive := time.NewTicker(25 * time.Second)
	defer keepalive.Stop()
	ctx := r.Context()

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
			c, ok := live[id]
			if !ok {
				continue
			}
			frag, err := s.render.CellHTML(s.renderCell(c).Tile)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", id, frag)
			if rc.Flush() != nil {
				return
			}
		}
	}
}
