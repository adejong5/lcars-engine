# Design decisions

Architecture Decision Records (ADRs) for this project. Each file captures one
significant decision: the context that forced it, the choice made, the
consequences, and the alternatives that were rejected and why. They are
append-only history — when a later decision changes an earlier one, add a new
ADR that supersedes it rather than rewriting the old one.

**Status legend:** `Accepted` (decided and in effect) · `Proposed` (a
provisional default we agreed to start with but expect to revisit) ·
`Superseded by NNNN`.

| # | Decision | Status |
|---|---|---|
| [0001](0001-old-hardware-is-a-primary-constraint.md) | Old-hardware browser support is a primary constraint | Accepted |
| [0002](0002-server-side-rendering.md) | Move dashboards to server-side rendering | Accepted |
| [0003](0003-htmx-for-partial-updates.md) | Use htmx (1.x) for partial updates | Accepted |
| [0004](0004-new-project-harvest-assets.md) | Build a new project; harvest assets (don't convert in place) | Accepted |
| [0005](0005-server-stack-go.md) | Server stack: Go (net/http + html/template) | Accepted |
| [0006](0006-server-side-ha-and-deployment.md) | Server-side HA connection; deploy as a container first | Accepted / Proposed |
| [0007](0007-live-updates-polling-first.md) | Live updates: htmx polling first, SSE later | Proposed |

Related background lives in [`../COMPATIBILITY.md`](../COMPATIBILITY.md) (the
browser-support audit) and [`../ENGINE.md`](../ENGINE.md) (dashboard design
guide for the current Svelte stack).
