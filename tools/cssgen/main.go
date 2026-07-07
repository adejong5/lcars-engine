// Command cssgen validates that the served stylesheets stay old-browser
// compatible. Both classic.compat.css (the theme, TheLCARS.com-derived) and
// components.compat.css (our component layer) are maintained sources, hand-
// authored directly in ~2017-approved CSS; this tool fails the build if a
// modern feature sneaks in.
//
// History: cssgen began as a generator converting the vendored upstream theme
// and a modern-authored components.css into the compat files (one adaptation
// pass per feature: text-box, :has, logical properties, flex-gap, clamp).
// Both sources were dropped once the project went feature-beyond the upstream
// (decisions/0008) — the converged compat output became the source, and only
// the tripwire remains. The pass implementations live in git history
// (`git log -- tools/cssgen`) if a conversion is ever needed again.
//
// Target baseline: ~2017 (CSS-grid era), with an eventual stretch to ~2015.
//
//	go run ./tools/cssgen
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// checked are the maintained compat stylesheets: every rule must satisfy the
// old-browser baseline.
var checked = []string{
	"internal/render/static/classic.compat.css",
	"internal/render/static/components.compat.css",
}

func main() {
	bad := false
	for _, src := range checked {
		raw, err := os.ReadFile(src)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read:", err)
			os.Exit(1)
		}
		if leftover := modernLeftovers(string(raw)); len(leftover) > 0 {
			fmt.Fprintf(os.Stderr, "%s: modern CSS not in the old-browser baseline: %s\n",
				src, strings.Join(leftover, ", "))
			bad = true
			continue
		}
		fmt.Println(src, " ok")
	}
	if bad {
		fmt.Fprintln(os.Stderr, "\nrewrite for ~2017 browsers: clamp/min/max -> @media breakpoints;")
		fmt.Fprintln(os.Stderr, ":has() -> restructure or class from the server; logical props ->")
		fmt.Fprintln(os.Stderr, "physical; flex gap -> child margins; grid gap -> grid-gap; text-box -> padding bias")
		os.Exit(1)
	}
}

var (
	cssComment = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	minMaxFn   = regexp.MustCompile(`(^|[^-A-Za-z])(min|max)\(`)
	insetProp  = regexp.MustCompile(`\binset[ \t]*:`)
	// gap/row-gap/column-gap: flex containers need margin fallbacks and grid
	// containers the grid-gap spelling on old engines. The prefix class
	// excludes grid-gap / grid-row-gap / grid-column-gap.
	gapProp = regexp.MustCompile(`(^|[;{\s])(gap|row-gap|column-gap)[ \t]*:`)
)

// modernLeftovers scans comment-stripped CSS for features newer than the
// ~2017 baseline and returns what it found.
func modernLeftovers(css string) []string {
	bare := cssComment.ReplaceAllString(css, "")
	var found []string
	for _, probe := range []string{"text-box:", ":has(", "clamp(",
		"margin-block", "margin-inline", "padding-block", "padding-inline",
		"border-block", "border-inline", "block-size", "inline-size"} {
		if strings.Contains(bare, probe) {
			found = append(found, probe)
		}
	}
	if m := minMaxFn.FindStringSubmatch(bare); m != nil {
		found = append(found, m[2]+"(")
	}
	if insetProp.MatchString(bare) {
		found = append(found, "inset:")
	}
	if m := gapProp.FindStringSubmatch(bare); m != nil {
		found = append(found, m[2]+":")
	}
	return found
}
