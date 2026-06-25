# LCARS Engine

Server-rendered Star Trek **LCARS** dashboards for Home Assistant, built to run
on **old hardware** — wall tablets, smart displays, and Raspberry Pis whose
browsers can't run the modern Home Assistant frontend.

This repository is an **htmx + Go** rebuild of the original Svelte project. The
server holds the Home Assistant connection and renders plain HTML; the browser
only swaps HTML fragments via htmx — no client framework runtime and no
ES-module build, so the supported-browser floor drops far below the ~2020 wall
of the Svelte version.

> **Status:** early — scaffolding in progress. The architecture and the
> reasoning behind every major choice are recorded as Architecture Decision
> Records in [`decisions/`](decisions/).

## Why a rewrite

See [`decisions/`](decisions/): old-hardware support is the driving constraint
([0001](decisions/0001-old-hardware-is-a-primary-constraint.md)), which led to
server-side rendering ([0002](decisions/0002-server-side-rendering.md)), htmx
for partial updates without full refreshes
([0003](decisions/0003-htmx-for-partial-updates.md)), a fresh project that
harvests the original's theme CSS and markup
([0004](decisions/0004-new-project-harvest-assets.md)), and a Go server
([0005](decisions/0005-server-stack-go.md)).

## Related

- **[`lcars-engine-js`](https://github.com/adejong5/lcars-engine-js)** — the
  original Svelte 5 + Vite implementation this project descends from. It remains
  the LCARS markup/CSS reference and hosts the live mock demo on GitHub Pages.
