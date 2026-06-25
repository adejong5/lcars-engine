# 0002 — Move dashboards to server-side rendering

- **Status:** Accepted
- **Date:** 2026-06-24
- **Driven by:** [0001](0001-old-hardware-is-a-primary-constraint.md)

## Context

The hard browser floor is the client JavaScript toolchain (ES-module entry
point, Vite's es2020 target, Svelte 5's `Proxy` runtime). No amount of CSS work
lowers it. To serve browsers below ~2020 we have to remove the framework
runtime and ESM-only delivery **from the client**.

## Decision

Render the dashboards **on the server**. A backend holds the Home Assistant
connection, renders the LCARS HTML from current entity state, and the browser
receives plain HTML. This eliminates all three hard floors at once: there is no
ES-module entry point, no es2020 bundle, and no `Proxy`-based client runtime to
require.

## Consequences

- **Requires a backend** (new operational component; see
  [0005](0005-server-stack-go.md) and
  [0006](0006-server-side-ha-and-deployment.md)).
- Home Assistant is talked to **server-side**, so the long-lived token lives in
  the server environment and never reaches the browser — this retires the
  "`dist/` is a secret-carrier" problem from the current static builds.
- We lose Svelte's client-side fine-grained reactivity; updates become
  server-rendered fragments. Acceptable for a status dashboard.
- The theme CSS can now be **flattened at build time** (PostCSS to inline
  `var()`, expand `clamp()`, add `:has()`/`gap` fallbacks) because we control a
  build step — something the current live-consumed stylesheets can't do cleanly.
- CSS features still gate **visual fidelity** (they don't stop the page
  loading); flattening addresses that separately.

## Alternatives considered

- **Keep the SPA and add `@vitejs/plugin-legacy` + polyfills.** Rejected:
  Svelte 5's `Proxy` runtime can't be polyfilled, so the floor stays ~2017 at
  best, for a lot of build complexity and little reach.
- **Pure no-JS server rendering with `<meta refresh>`.** Considered and kept as
  the conceptual ultimate-compatibility fallback, but rejected as the primary
  approach because of full-region reload flicker and loss of partial updates.
  See [0003](0003-htmx-for-partial-updates.md).
