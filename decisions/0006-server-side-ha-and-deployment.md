# 0006 — Server-side HA connection; deploy as a container first

- **Status:** Accepted (connection model) · Proposed (hosting target)
- **Date:** 2026-06-24 (add-on readiness added 2026-06-26)
- **Driven by:** [0002](0002-server-side-rendering.md), [0005](0005-server-stack-go.md)

## Context

In the SSR model the server, not the browser, holds the Home Assistant
connection and credentials. Two things to settle: how the server talks to HA,
and where the server runs.

## Decision

**Connection (Accepted):** the server runs a Go **HA client** — WebSocket
subscription to `state_changed` for live state, and `call_service` over the same
WebSocket — caching the latest entity states for rendering. The **long-lived
token lives only in the server environment**; it never reaches the browser. A
**mock mode** (mirroring the current repo's fabricated-entity behavior) supports
offline layout work.

**Hosting (Proposed — start simple):** run the server as a **plain Docker
container on pandaserver** first. The **HA add-on** path is a later option once
the app is proven — but we keep it cheap by building the enablers in now (below).

**Add-on readiness — built in now (Accepted):** the goal is one binary that runs
both as a standalone container and as an add-on, serving identical content.

- **Ingress chosen** over an exposed port, for HA's authenticated proxy and the
  in-sidebar panel. The cost is one knob: HA serves the add-on under a
  per-request path prefix (`X-Ingress-Path`). The server reads that header
  (middleware, done) and templates will render `<base href="{prefix}/">` and use
  **relative URLs throughout**, so the same markup works at `/` (standalone) or
  under the prefix (add-on). No absolute paths anywhere.
- **Supervisor-aware connection (done):** when `ADDON` is set (auto-detected
  from `SUPERVISOR_TOKEN`), the client uses `ws://supervisor/core/websocket` with
  `SUPERVISOR_TOKEN` instead of `HA_HOST`/`HA_TOKEN`.
- **Multi-arch static binary (done):** the Dockerfile cross-compiles a
  CGO-free binary via buildx `TARGETOS`/`TARGETARCH` into a `scratch` image, so
  the same build serves every add-on arch (aarch64/armv7/amd64).
- **Env-driven config (done):** a later `run.sh` with `bashio` maps the add-on's
  `options.json` to env vars with no app change.

**Deferred:** the add-on manifest itself — `config.yaml` (slug, arch list,
`ingress: true`, `ingress_port`, `homeassistant_api: true`, options schema),
`build.yaml`, `run.sh`, and the add-on repository structure — is written once the
UI exists and the app is versioned.

## Consequences

- Retires the "`dist/` is a secret-carrier" problem entirely — no token in any
  bundle.
- Fits existing infrastructure (pandaserver already serves the current `dist/`
  via nginx).
- The add-on upgrade later adds ingress/auth and tokenless HA access without
  changing the app's core.
- Requires whichever HA install supports add-ons **only if** we take the add-on
  path; the container path is install-agnostic.

## Alternatives considered

- **HA add-on first.** More setup up front (config.yaml, ingress, Supervisor
  API) and depends on an HA OS / Supervised install. Deferred, not rejected.
- **Keep the client-side HA WebSocket.** Rejected: it leaks the token into the
  browser and requires a modern browser — the opposite of the project goal.
