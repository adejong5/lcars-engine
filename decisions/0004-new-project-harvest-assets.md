# 0004 — Build a new project; harvest assets (don't convert in place)

- **Status:** Accepted
- **Date:** 2026-06-24
- **Driven by:** [0002](0002-server-side-rendering.md), [0003](0003-htmx-for-partial-updates.md)

## Context

With server-side rendering + htmx decided, should we convert the existing
Svelte repo in place, or start a fresh project?

The key fact: a Svelte SPA (client renders, client opens the HA WebSocket) and
an htmx SSR app (server renders, server holds the HA connection) **share no
executable code**. The `.svelte` components get reimplemented as server
templates either way.

## Decision

**Start a new project and harvest** from the existing Svelte engine rather than
migrating it.

**Naming (settled 2026-06-25):** the htmx rewrite takes the canonical
**`lcars-engine`** slug (a fresh repo + working dir), and the original Svelte
project is renamed to **`lcars-engine-js`** — the `-js` marks that the
JavaScript/Svelte client runtime is its defining trait. The Svelte project keeps
its history and its GitHub Pages mock demo on the renamed slug.

What transfers cleanly into the new project:

- the vendored theme CSS (`shared/assets/*.css`) — reused as-is, later flattened,
- fonts, sounds, MDI icon usage,
- the **exact LCARS markup** — each `shared/components/*.svelte` is read as the
  spec for the HTML its server template must emit (the class contract is the
  portable value),
- `ENGINE.md` design rules and `COMPATIBILITY.md`,
- the entity mappings from the private dashboards.

Repos mirror the existing public/private split:

- **`lcars-engine`** (new, canonical slug) — the htmx/Go server engine: vendors
  the theme CSS + assets, holds the backend and templates. Shareable/public, no
  secrets. The decision ADRs now live here.
- **`lcars-engine-js`** (the renamed original) — the Svelte reference:
  markup/CSS/design source of truth and the GitHub Pages mock demo.
- **deploy config + HA token** stay private (env / deployments repo).

## Consequences

- `lcars-engine-js` keeps working as documentation and as the only thing that
  can run on GitHub Pages (an htmx server can't).
- No mixing of a Vite/Svelte toolchain with the Go server in one tree; no two
  parallel rendering systems to maintain in one repo.
- Components must be re-authored as templates — but that work is unavoidable
  regardless of repo choice.

## Alternatives considered

- **Convert the Svelte repo in place.** Rejected: it would **kill the GitHub
  Pages demo** (Pages is static; an htmx server can't run there), mix two
  incompatible runtimes and build systems in one repo, and save no component
  work (nothing executable transfers anyway).
