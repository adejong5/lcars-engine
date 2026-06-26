package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIngressBasePath(t *testing.T) {
	s := &Server{}
	var seen string
	h := s.ingress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = BasePath(r)
	}))

	// Behind HA ingress: trailing slash is trimmed.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Ingress-Path", "/api/hassio_ingress/abc123/")
	h.ServeHTTP(httptest.NewRecorder(), req)
	if seen != "/api/hassio_ingress/abc123" {
		t.Fatalf("base path = %q", seen)
	}

	// Direct (standalone): no header → empty.
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if seen != "" {
		t.Fatalf("direct base path should be empty, got %q", seen)
	}
}
