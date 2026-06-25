# 0005 — Server stack: Go (net/http + html/template)

- **Status:** Accepted
- **Date:** 2026-06-24 (revised 2026-06-25)
- **Driven by:** [0004](0004-new-project-harvest-assets.md)

## Context

The new server ([0004](0004-new-project-harvest-assets.md)) needs a language and
template engine. Candidates: Python (FastAPI + Jinja2), Go (net/http +
html/template), Node (Express/Fastify + a template engine).

This was first decided as **Python**, on the strength of the Home Assistant
ecosystem (HA is Python) and AI-authorability. It was **revised to Go** once we
committed to also shipping a **Home Assistant add-on** and made **portability** a
goal — we want one artifact that runs well as a plain container *and* as a lean
add-on on any host, regardless of where HA ends up. That reweights the decision:

- **Footprint and multi-arch builds become first-class.** An add-on is a Docker
  container the Supervisor runs, and it must ship for aarch64/armv7/amd64. Go
  cross-compiles to a **single static binary** per arch (`FROM scratch`, ~10–20
  MB image, ~10 MB RAM, fast Supervisor installs). Python/Node images are larger
  and slower to build on armv7 (native wheels/modules).
- **The Supervisor proxy neutralizes Python's main advantage.** Inside an add-on
  you reach HA over plain HTTP/WebSocket via `SUPERVISOR_TOKEN`, so the rich
  Python HA client libraries aren't needed — a thin client in any language
  suffices. That was the core of the original Python rationale, and the add-on
  context removes it.

The current HA host is a powerful machine, so footprint isn't a problem *today* —
but **portability/future-proofing wins**: a Go binary stays lean everywhere,
which matches the project's whole minimalism thesis (the same reason we chose
htmx on the client, [0003](0003-htmx-for-partial-updates.md)).

## Decision

**Go, with the standard library: `net/http` for serving and `html/template`
for rendering.** The server holds the HA WebSocket connection, caches entity
state, renders LCARS fragments from `html/template`, and ships as one portable
static binary — wrapped unchanged as either a pandaserver container or an HA
add-on ([0006](0006-server-side-ha-and-deployment.md)).

## Consequences

- One small static binary, trivially cross-compiled for every add-on arch;
  minimal image, RAM, and attack surface; fast Supervisor install/update.
- `html/template` reproduces the LCARS class-contract; partials map to the
  Svelte components read as specs ([0004](0004-new-project-harvest-assets.md)).
- The HA WebSocket/REST client is hand-written Go (small surface: subscribe to
  `state_changed`, call services) — more glue than Python's libraries, but
  modest for this scope.
- Slightly less "AI-authorable" than Python, accepted as a fair trade for
  portability and footprint.

## Alternatives considered

- **Python (FastAPI + Jinja2).** The original choice: best HA ecosystem,
  highest AI-authorability, Jinja2 maps cleanly to the markup. Rejected once
  add-on portability raised footprint/multi-arch weight and the Supervisor proxy
  removed the need for Python's HA libraries. Still the fallback if hand-written
  Go HA glue proves more costly than expected.
- **Node (Express/Fastify + templates).** Stays in JavaScript but keeps a JS
  toolchain we're leaving and is heavier than Go. Rejected.
