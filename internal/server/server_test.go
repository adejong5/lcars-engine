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
	"github.com/adejong5/lcars-engine/internal/render"
)

func newTestServer() http.Handler {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	rnd, err := render.New(false)
	if err != nil {
		panic(err)
	}
	return New(config.Config{Mock: true, DebugAPI: true}, log, ha.NewMock(), rnd).Handler()
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

func TestOpsPageRenders(t *testing.T) {
	h := newTestServer()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("GET / status %d", rr.Code)
	}
	body := rr.Body.String()
	for _, s := range []string{`<base href="/">`, `class="banner"`, "lcars-text-bar", "TheLCARS.com", `href="static/classic.css"`} {
		if !strings.Contains(body, s) {
			t.Fatalf("page missing %q", s)
		}
	}
}

func TestOpsPageIngressBase(t *testing.T) {
	h := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Ingress-Path", "/api/hassio_ingress/xyz/")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if !strings.Contains(rr.Body.String(), `<base href="/api/hassio_ingress/xyz/">`) {
		t.Fatal("ingress base href not applied")
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
