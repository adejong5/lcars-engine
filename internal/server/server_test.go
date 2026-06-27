package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adejong5/lcars-engine/internal/config"
	"github.com/adejong5/lcars-engine/internal/ha"
)

func newTestServer() http.Handler {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(config.Config{Mock: true}, log, ha.NewMock()).Handler()
}

func TestCallServiceTogglesMockState(t *testing.T) {
	h := newTestServer()

	body := `{"domain":"switch","service":"turn_on","target":{"entity_id":"switch.test"}}`
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/call", strings.NewReader(body)))
	if rr.Code != http.StatusOK {
		t.Fatalf("POST /api/call status %d, want 200; body=%s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/state/switch.test", nil))
	if !strings.Contains(rr.Body.String(), `"state":"on"`) {
		t.Fatalf("entity not turned on: %s", rr.Body.String())
	}
}

func TestCallServiceValidation(t *testing.T) {
	h := newTestServer()

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/call", strings.NewReader(`{}`)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("missing domain/service: status %d, want 400", rr.Code)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/call", strings.NewReader(`not json`)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("bad JSON: status %d, want 400", rr.Code)
	}
}

func TestHealthz(t *testing.T) {
	h := newTestServer()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"status":"ok"`) {
		t.Fatalf("healthz: %d %s", rr.Code, rr.Body.String())
	}
}
