# 0007 — Live updates: htmx polling first, SSE later

- **Status:** Proposed (provisional default)
- **Date:** 2026-06-24
- **Driven by:** [0003](0003-htmx-for-partial-updates.md)

## Context

The dashboard needs live data and controls without full-page refreshes. With
htmx chosen, the transport options are: periodic **polling**, **Server-Sent
Events (SSE)**, or **WebSocket**.

## Decision (provisional)

**Start with htmx polling.** Tiles refresh on `hx-trigger="every Ns"`,
consolidated where possible into a single poller that returns **out-of-band
fragments** so many tiles update from one request. **Controls** are
`hx-post` forms that call the HA service and return the updated fragment (these
also work as plain `<form>` posts with no JS).

Treat **SSE** (`EventSource`, ~2011+) as the ready upgrade if polling proves too
chatty — the server pushes OOB fragments on HA `state_changed`. **WebSocket** is
held in reserve and only if a future need justifies it.

## Consequences

- **Maximum compatibility** to start: polling is plain XHR on a timer, supported
  everywhere htmx runs.
- Some redundant requests; mitigated by consolidating polls and tuning the
  interval.
- The SSE upgrade keeps the floor at ~2011 and reduces traffic without changing
  the template/fragment design.

## Alternatives considered

- **SSE first.** More efficient push and still old-browser-friendly (~2011), but
  slightly more server code and a connection to manage. Deferred as an upgrade,
  not the starting point.
- **WebSocket.** Lowest-latency but heaviest and unnecessary for periodic
  telemetry. Held in reserve.

## Status note

Marked **Proposed** because it is an agreed starting default we expect to
revisit (likely promoting SSE) once the first pages are running.
