# 0008 — Drop the CSS generation pipeline; the compat build is the source

- **Status:** Accepted
- **Date:** 2026-07-06

## Context

Until now the TheLCARS.com "classic" theme was vendored byte-identical as
`classic.css`, our component layer was authored with modern conveniences in
`components.css`, and `tools/cssgen` generated the served `*.compat.css` from
both (ADR 0001's fallback-layer policy). Keeping the pristine theme copy had
one purpose: re-syncing with upstream template releases.

We are now feature-beyond the upstream — bezel control clusters, kiosk
sizing, segmented meters, readout stacks, the compat centering fixes — and a
future upstream pull would be a redesign, not a merge. Meanwhile every CSS
change ran through the generate step, and the served files were one step
removed from what anyone edited.

## Decision

- Delete the vendored `classic.css` and our modern-authored `components.css`.
  **The served `classic.compat.css` and `components.compat.css` are the
  maintained sources**, hand-edited directly in old-browser-approved CSS.
- Every rule must satisfy the ~2017 baseline: `@media` breakpoints instead of
  `clamp()`/`min()`/`max()`, no `:has()`, physical properties instead of
  logical, sibling margins instead of flex `gap`, `grid-gap` instead of grid
  `gap`, no `text-box`.
- `tools/cssgen` becomes a pure validator: it scans both files for
  modern-feature leftovers and fails the build on a hit. The adaptation-pass
  implementations were deleted (recoverable from git history if a conversion
  is ever needed again).
- Attribution stands: the theme remains derived from the LCARS Inspired
  Website Template by Jim Robertus (TheLCARS.com); LICENSE.md, the file
  header, and the in-page credit line are kept.
- The csscheck harness keeps its self-contained checks (clamp probes, the
  panel-2 elbow table vs a live `max()` rule); the checks that compared
  compat output against a modern original (textbox, flexgap, logical) are
  retired with the originals.
- The theme's colour vocabulary (`--<name>` variables and the
  `font-/background-/button-/bullet-*` utility families) is kept and treated
  as API — other themes may restyle it; the demo page documents it.

## Consequences

- CSS edits are simpler (edit the served file → validate) but must be written
  compat-first; the validator fails the build if modern CSS sneaks in.
- Re-syncing with future TheLCARS.com releases is no longer supported; any
  upstream idea gets ported by hand.
- The generated-era `classic.css`/`components.css` and the cssgen passes
  remain available in git history (`dd22639` and earlier) for archaeology.
- ADR 0001's "fallbacks never in the vendored files" rule is superseded:
  there are no vendored files, and fixes are made in place.
