package server

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestIndexDemoRenders(t *testing.T) {
	h := newTestServer()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("GET / status %d", rr.Code)
	}
	body := rr.Body.String()
	for _, s := range []string{`<base href="/">`, `class="pill-gauge"`, `class="bar-chart"`, `class="readout"`, `href="static/classic.compat.css"`, "components.compat.css"} {
		if !strings.Contains(body, s) {
			t.Fatalf("demo page missing %q", s)
		}
	}
}

// TestThemeOverlay: THEME adds the colour overlay stylesheet after the base
// pair; unset serves the classic base only.
func TestThemeOverlay(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	rnd, err := render.New(false)
	if err != nil {
		t.Fatal(err)
	}
	h := New(config.Config{Mock: true, Theme: "nemesis-blue"}, log, ha.NewMock(), rnd).Handler()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(rr.Body.String(), `href="static/theme-nemesis-blue.css"`) {
		t.Fatal("theme overlay link missing")
	}

	rr = httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if strings.Contains(rr.Body.String(), "static/theme-") {
		t.Fatal("theme overlay linked with no THEME set")
	}

	// the overlay files themselves are served
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/static/theme-nemesis-blue.css", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("theme stylesheet not served: %d", rr.Code)
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

// TestSSEStreamsThroughMiddleware guards against the logging middleware's
// ResponseWriter wrapper breaking http.ResponseController: without Unwrap, the
// stream never flushes and the client hangs with 0 bytes received. Uses a real
// server so the flush actually crosses a socket.
func TestSSEStreams(t *testing.T) {
	ts := httptest.NewServer(newTestServer())
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/sse", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /sse: %v", err)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("SSE Content-Type = %q", ct)
	}

	// The handler writes ": connected" and flushes immediately; if flushing were
	// broken this read would block until the context deadline instead.
	line, err := bufio.NewReader(resp.Body).ReadString('\n')
	if err != nil {
		t.Fatalf("reading first SSE flush: %v", err)
	}
	if !strings.Contains(line, "connected") {
		t.Fatalf("first SSE line = %q, want the connected comment", line)
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
