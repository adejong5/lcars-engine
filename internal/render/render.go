// Package render serves the LCARS UI: embedded html/template pages, a shared
// component library, and the static theme assets (CSS, fonts, htmx). Pages are
// rendered server-side; the browser receives plain HTML.
package render

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed templates/*.html templates/pages/*.html
var tmplFS embed.FS

//go:embed static
var staticFS embed.FS

// funcs are template helpers. The template.CSS-returning ones let pages set
// dynamic inline styles (widths, colours) without html/template refusing them.
var funcs = template.FuncMap{
	"widthpct": func(f float64) template.CSS { return template.CSS(fmt.Sprintf("width:%.2f%%", f*100)) },
	"clipinset": func(f float64) template.CSS {
		return template.CSS(fmt.Sprintf("clip-path:inset(0 %.2f%% 0 0)", (1-f)*100))
	},
	"css": func(s string) template.CSS { return template.CSS(s) },
}

// Renderer holds one parsed template set per page (each = frame + components +
// that page), plus the shared set used to render standalone fragments.
type Renderer struct {
	pages map[string]*template.Template
	base  *template.Template // frame + components; renders fragments like "cell"
}

// New parses the embedded templates.
func New(dev bool) (*Renderer, error) {
	base, err := template.New("base").Funcs(funcs).ParseFS(tmplFS, "templates/frame.html", "templates/components.html")
	if err != nil {
		return nil, err
	}

	pageFiles, err := fs.Glob(tmplFS, "templates/pages/*.html")
	if err != nil {
		return nil, err
	}
	pages := make(map[string]*template.Template, len(pageFiles))
	for _, pf := range pageFiles {
		clone, err := base.Clone()
		if err != nil {
			return nil, err
		}
		if _, err := clone.ParseFS(tmplFS, pf); err != nil {
			return nil, err
		}
		name := strings.TrimSuffix(path.Base(pf), ".html")
		pages[name] = clone
	}
	return &Renderer{pages: pages, base: base}, nil
}

// Page renders the named page (its "content") inside the frame.
func (r *Renderer) Page(w http.ResponseWriter, name string, data any) error {
	t, ok := r.pages[name]
	if !ok {
		return fmt.Errorf("render: unknown page %q", name)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.ExecuteTemplate(w, "frame", data)
}

// Fragment renders a named standalone fragment (e.g. "cell", "toggle"), used by
// the poll and action endpoints.
func (r *Renderer) Fragment(w http.ResponseWriter, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.base.ExecuteTemplate(w, name, data)
}

// FragmentHTML renders a named fragment to a single-line string for SSE event
// data (the SSE wire format is line-oriented, so embedded newlines are
// flattened to spaces — harmless in HTML).
func (r *Renderer) FragmentHTML(name string, data any) (string, error) {
	var b strings.Builder
	if err := r.base.ExecuteTemplate(&b, name, data); err != nil {
		return "", err
	}
	return strings.Join(strings.Fields(b.String()), " "), nil
}

// Static serves the embedded theme assets under /static/.
func Static() http.Handler {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err) // embedded; cannot fail at runtime
	}
	return http.StripPrefix("/static/", http.FileServer(http.FS(sub)))
}
