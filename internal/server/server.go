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
	"github.com/adejong5/lcars-engine/internal/ha"
)

// Server holds shared dependencies for the HTTP handlers.
type Server struct {
	cfg config.Config
	log *slog.Logger
	src ha.Source
	mux *http.ServeMux
}

// New builds a Server with its routes registered.
func New(cfg config.Config, log *slog.Logger, src ha.Source) *Server {
	s := &Server{cfg: cfg, log: log, src: src, mux: http.NewServeMux()}
	s.routes()
	return s
}

// Handler returns the root handler with middleware applied. ingress is
// outermost so the X-Ingress-Path prefix is available to every handler.
func (s *Server) Handler() http.Handler {
	return s.ingress(s.logRequests(s.mux))
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealth)
	s.mux.HandleFunc("GET /api/states", s.handleStates)
	s.mux.HandleFunc("GET /api/state/{id}", s.handleState)
	s.mux.HandleFunc("POST /api/call", s.handleCall)
}

// handleCall invokes a Home Assistant service. Terminal-verifiable:
//
//	curl -X POST localhost:8080/api/call \
//	  -d '{"domain":"switch","service":"turn_on","target":{"entity_id":"switch.kitchen_spare"}}'
func (s *Server) handleCall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain      string         `json:"domain"`
		Service     string         `json:"service"`
		ServiceData map[string]any `json:"service_data"`
		Target      map[string]any `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid JSON body"})
		return
	}
	if req.Domain == "" || req.Service == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "domain and service are required"})
		return
	}
	if err := s.src.CallService(req.Domain, req.Service, req.ServiceData, req.Target); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"domain":  req.Domain,
		"service": req.Service,
		"target":  req.Target,
	})
}

// handleStates returns every known entity state. Terminal-verifiable:
//
//	curl -s localhost:8080/api/states
func (s *Server) handleStates(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.src.All())
}

// handleState returns one entity. In mock mode unknown ids are fabricated on
// read (always 200); the live client (Phase 3) returns 404 for unknown ids.
//
//	curl -s localhost:8080/api/state/sensor.cpu
func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, ok := s.src.State(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "unknown entity", "entity_id": id})
		return
	}
	writeJSON(w, http.StatusOK, st)
}

// handleHealth reports liveness and the current mode. Terminal-verifiable:
//
//	curl -s localhost:8080/healthz
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"mock":   s.cfg.Mock,
		"addon":  s.cfg.Addon,
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
