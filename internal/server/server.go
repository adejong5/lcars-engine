// Package server wires the HTTP routes and middleware for the LCARS engine.
//
// Phase 1 exposes only /healthz. Later phases add the HA-backed JSON debug API
// (Phase 2), service-call actions (Phase 4), and the rendered pages (Phase 7+).
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/adejong5/lcars-engine/internal/config"
)

// Server holds shared dependencies for the HTTP handlers.
type Server struct {
	cfg config.Config
	log *slog.Logger
	mux *http.ServeMux
}

// New builds a Server with its routes registered.
func New(cfg config.Config, log *slog.Logger) *Server {
	s := &Server{cfg: cfg, log: log, mux: http.NewServeMux()}
	s.routes()
	return s
}

// Handler returns the root handler with middleware applied.
func (s *Server) Handler() http.Handler {
	return s.logRequests(s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealth)
}

// handleHealth reports liveness and the current mode. Terminal-verifiable:
//
//	curl -s localhost:8080/healthz
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"mock":   s.cfg.Mock,
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// logRequests emits one structured line per request: method, path, status, dur.
func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		s.log.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"dur", time.Since(start).String(),
		)
	})
}

// statusRecorder captures the response status for logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
