// Package render serves the LCARS UI: embedded html/template pages and the
// static theme assets (CSS, fonts, htmx). Pages are rendered server-side; the
// browser receives plain HTML (Phase 7).
package render

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed templates/*.html
var tmplFS embed.FS

//go:embed static
var staticFS embed.FS

// Renderer holds the parsed page templates.
type Renderer struct {
	tmpl *template.Template
}

// New parses the embedded templates. dev is reserved for future live-reload;
// templates are embedded, so this parses once.
func New(dev bool) (*Renderer, error) {
	t, err := template.New("render").ParseFS(tmplFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Renderer{tmpl: t}, nil
}

// Page renders a full page: the frame with the page's "content" filled in.
func (r *Renderer) Page(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.tmpl.ExecuteTemplate(w, "frame", data)
}

// Cell renders one live cell's inner fragment (label + value) to the writer,
// used by the fallback-poll endpoint.
func (r *Renderer) Cell(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.tmpl.ExecuteTemplate(w, "cell", data)
}

// CellHTML renders one live cell's inner fragment to a string, used by the SSE
// endpoint (single line, safe to put on a data: line).
func (r *Renderer) CellHTML(data any) (string, error) {
	var b strings.Builder
	if err := r.tmpl.ExecuteTemplate(&b, "cell", data); err != nil {
		return "", err
	}
	return b.String(), nil
}

// Static serves the embedded theme assets under /static/.
func Static() http.Handler {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err) // embedded; cannot fail at runtime
	}
	return http.StripPrefix("/static/", http.FileServer(http.FS(sub)))
}
