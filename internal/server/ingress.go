package server

import (
	"context"
	"net/http"
	"strings"
)

// ingress support for running behind Home Assistant's authenticated reverse
// proxy. HA serves the add-on under a per-request path prefix and passes it in
// the X-Ingress-Path header. Templates render a <base href="{BasePath}/"> and
// use relative URLs throughout, so the same markup works whether the app is
// mounted under that prefix (add-on) or at / (standalone container).

type ctxKey int

const basePathKey ctxKey = iota

// ingress captures X-Ingress-Path into the request context. It is the outermost
// middleware so the prefix is available to every handler/template.
func (s *Server) ingress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bp := strings.TrimRight(r.Header.Get("X-Ingress-Path"), "/")
		r = r.WithContext(context.WithValue(r.Context(), basePathKey, bp))
		next.ServeHTTP(w, r)
	})
}

// BasePath returns the ingress path prefix for this request, or "" when the app
// is reached directly (standalone container).
func BasePath(r *http.Request) string {
	if v, ok := r.Context().Value(basePathKey).(string); ok {
		return v
	}
	return ""
}
