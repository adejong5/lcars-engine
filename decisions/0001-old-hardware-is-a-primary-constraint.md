# 0001 — Old-hardware browser support is a primary constraint

- **Status:** Accepted
- **Date:** 2026-06-24

## Context

The project exists to drive **old hardware** — wall-mounted tablets, smart
displays, and Raspberry Pis whose browsers can't run the modern Home Assistant
frontend. The aspiration is to support browsers **as old as ~2015**.

The audit in [`../COMPATIBILITY.md`](../COMPATIBILITY.md) showed the current
Svelte 5 + Vite stack cannot meet that. The hard floor is the **JavaScript
toolchain, not the CSS**:

- the only entry point is `<script type="module">` → native ES modules (~2017),
- Vite's default `build.target` emits `es2020` (~2020),
- Svelte 5's reactivity is `Proxy`-based and cannot be polyfilled,
- there is no `nomodule`/legacy bundle.

So the realistic floor for the existing stack is **~2020 (Chrome 87 / Safari
14)** — roughly five years above the goal. A literal-2015 browser can't load
the app at all.

## Decision

Treat **~2015-era browser support as a first-class, architecture-driving
constraint**, not an afterthought. Concretely:

- The current static Svelte stack is documented as a **~2020 functional
  baseline** and kept as a reference/demo (it is not abandoned).
- To reach materially older hardware we accept that a **different runtime
  architecture is required**, and pursue it deliberately (see
  [0002](0002-server-side-rendering.md)).
- Visual-fidelity features (CSS) must remain **flatten-able / fall-back-able**;
  browser-specific fallbacks live in our own layer, never in the vendored
  theme CSS (see `CLAUDE.md` → "Upstream template watch-outs").

## Consequences

- Justifies the move to server-side rendering and htmx
  ([0002](0002-server-side-rendering.md), [0003](0003-htmx-for-partial-updates.md)).
- Every new feature is evaluated against the support floor; bleeding-edge CSS
  from upstream template pulls is a recurring risk to re-audit.
- We must be honest about what each architecture actually supports rather than
  claiming "2015" loosely.

## Alternatives considered

- **Accept the ~2020 floor and do nothing.** Rejected: it abandons the oldest
  target devices, which are the entire reason the project exists.
