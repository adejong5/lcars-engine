# 0007 — Live updates: SSE primary, infrequent polling fallback

- **Status:** Accepted
- **Date:** 2026-06-24 (resolved 2026-06-28)
- **Driven by:** [0003](0003-htmx-for-partial-updates.md)

## Context

The dashboard needs live data and controls without full-page refreshes. With
htmx chosen, the transport options are: periodic **polling**, **Server-Sent
Events (SSE)**, or **WebSocket**. This ADR was originally "polling first, SSE
later" as a provisional default; it is now resolved to ship both at once, with
SSE as the primary path.

## Decision

**SSE is the primary update path; an infrequent poll is the fallback.**

- The server holds the one HA WebSocket and broadcasts entity changes
  internally (the store's subscribe mechanism). A `GET /sse` endpoint streams
  per-entity HTML fragments; the htmx **sse extension** swaps each into the
  element whose `sse-swap` matches the event name (the entity id).
- Each live cell **also** has `hx-get="cells/{id}" hx-trigger="every 60s"` — an
  **infrequent fallback poll** that reconciles anything SSE missed (a dropped
  connection, a dropped notification under load, or a client without SSE). The
  same fragment renderer serves both paths.
- **Controls** are `hx-post` forms that call the HA service and return the
  updated fragment (and degrade to plain `<form>` posts with no JS).

Because SSE delivers promptly, the poll is a safety net, not the mechanism — so
its interval is deliberately long (60s) to keep traffic negligible.

## Consequences

- Prompt updates with little traffic; the slow poll guarantees eventual
  consistency without hammering the server.
- The SSE endpoint must clear the server's `WriteTimeout` per-connection
  (`http.ResponseController`) and emit periodic keepalive comments so proxies
  don't close the idle stream.
- `EventSource` is ~2011+, comfortably within the support floor; the fallback
  poll covers anything older or any dropped stream.

## Alternatives considered

- **Polling only (the original provisional default).** Simpler, maximal
  compatibility, but either chatty (short interval) or laggy (long interval).
  Superseded: SSE gives promptness while the long poll retains the
  compatibility safety net.
- **WebSocket to the browser.** Lowest-latency but heaviest, reintroduces a
  browser-side socket, and is unnecessary for periodic telemetry. Held in
  reserve.
