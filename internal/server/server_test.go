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
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/ops", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /ops status %d", rr.Code)
	}
	body := rr.Body.String()
	for _, s := range []string{`class="banner"`, "lcars-text-bar", "TheLCARS.com", `href="static/classic.compat.css"`} {
		if !strings.Contains(body, s) {
			t.Fatalf("ops page missing %q", s)
		}
	}
}

func TestIndexDemoRenders(t *testing.T) {
	h := newTestServer()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("GET / status %d", rr.Code)
	}
	body := rr.Body.String()
	for _, s := range []string{`<base href="/">`, `class="pill-gauge"`, `class="bar-chart"`, "components.compat.css"} {
		if !strings.Contains(body, s) {
			t.Fatalf("demo page missing %q", s)
		}
	}
}

func TestIngressBase(t *testing.T) {
	h := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Ingress-Path", "/api/hassio_ingress/xyz/")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if !strings.Contains(rr.Body.String(), `<base href="/api/hassio_ingress/xyz/">`) {
		t.Fatal("ingress base href not applied")
	}
}

func TestActionToggle(t *testing.T) {
	h := newTestServer()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/action/kitchen_spare", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("POST /action/kitchen_spare status %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Kitchen Spare") || !strings.Contains(body, "hx-post=") {
		t.Fatalf("action did not return the toggle fragment: %s", body)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/action/nope", nil))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown action status %d, want 404", rr.Code)
	}
}

func TestCellFragment(t *testing.T) {
	h := newTestServer()

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/cells/switch.kitchen_spare", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /cells/{id} status %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Kitchen Spare") {
		t.Fatalf("cell fragment missing label: %s", rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/cells/sensor.not_a_cell", nil))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown cell status %d, want 404", rr.Code)
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
