# 0006 — Server-side HA connection; deploy as a container first

- **Status:** Accepted (connection model) · Proposed (hosting target)
- **Date:** 2026-06-24
- **Driven by:** [0002](0002-server-side-rendering.md), [0005](0005-server-stack-go.md)

## Context

In the SSR model the server, not the browser, holds the Home Assistant
connection and credentials. Two things to settle: how the server talks to HA,
and where the server runs.

## Decision

**Connection (Accepted):** the server runs a Go **HA client** — WebSocket
subscription to `state_changed` for live state, REST for service calls — and
caches the latest entity states for rendering. The **long-lived token lives only
in the server environment**; it never reaches the browser. A **mock mode**
(mirroring the current repo's fabricated-entity behavior) supports offline
layout work.

**Hosting (Proposed — start simple):** run the server as a **plain Docker
container on pandaserver** first. Keep the **HA add-on** path (Supervisor
`ingress` + `SUPERVISOR_TOKEN` for tokenless HA access, plus HA-authenticated
access) as a later option once the app is proven.

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
