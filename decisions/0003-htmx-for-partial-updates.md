# 0003 — Use htmx (1.x) for partial updates

- **Status:** Accepted
- **Date:** 2026-06-24
- **Driven by:** [0001](0001-old-hardware-is-a-primary-constraint.md), [0002](0002-server-side-rendering.md)

## Context

Server-side rendering is decided. The remaining question is how the browser
updates **without full-page refreshes**. A full `<meta refresh>` reload is
maximally compatible but flickers and reloads the whole region. We want partial
DOM updates with the broadest old-browser support and the least client weight.
Compared: htmx, Turbo (Hotwire), Unpoly, and no-JS meta-refresh.

## Decision

Use **htmx, pinned to the 1.x line**. htmx makes a background XHR on a trigger,
the server returns an **HTML fragment**, and htmx swaps it into a target
element — no navigation, no asset re-fetch, no flicker.

- It is a single dependency-free **classic script (~14 kB)** with no build and
  no ES-module requirement — so it sidesteps the client-JS floor from
  [0002](0002-server-side-rendering.md).
- **1.x has the best old-browser reach** (~2014+, and the htmx project
  deliberately maintains 1.x for legacy hardware). 2.x dropped the oldest
  browsers, so we pin 1.x.
- Live data uses **polling / out-of-band swaps** to start, with SSE as an
  upgrade path (see [0007](0007-live-updates-polling-first.md)).

## Consequences

- Adds a small JS dependency: the floor rises from "any browser" (pure no-JS)
  to **~2014**, which is within the old-hardware target.
- The server must return HTML fragments, not just whole pages; controls become
  fragment-returning endpoints.
- Lowest learning curve and least payload of the candidates — good for a weak
  Pi CPU.

## Alternatives considered

- **Turbo (Hotwire).** Turbo Streams is the nicest server-*push* model, but the
  browser floor is highest (~2018–20, ES2017 assumptions), it leans on a
  bundler/importmap, and it's Rails-flavored. Rejected for the compatibility
  goal.
- **Unpoly.** Most batteries-included (polling, overlays, validation,
  transitions) and historically broad browser support, but the largest payload,
  ships its own CSS, and is more than a status dashboard needs. Rejected, but
  the closest runner-up.
- **No-JS `<meta refresh>` / iframe refresh.** Zero client JS, universal
  support, and naturally tokenless — but full-region reloads flicker and lose
  partial updates. Rejected as primary; retained as the ultimate-compatibility
  fallback if a target device can't run htmx 1.x.
