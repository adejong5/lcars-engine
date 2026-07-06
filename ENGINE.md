# ENGINE.md — Building LCARS dashboards with this engine

This is a guide for an AI assistant assembling dashboard pages from the
component blocks in `internal/render/templates/` and the page tables in
`internal/server/pages.go`. It assumes you have already read the repo layout
(server-rendered Go + htmx; the vendored theme; the cssgen compat build). This
document is about *how to design a page*: structuring information, picking
components for each kind of data, and knowing when a value should be a
read-only display versus an interactive control.

Ported from the original Svelte engine's ENGINE.md; the design guidance is
unchanged, the mechanics are this repo's.

---

## 1. What LCARS is

LCARS ("Library Computer Access/Retrieval System") is the fictional Star Trek
computer interface. Visually it is defined by a few strong conventions, and a
good dashboard leans into all of them:

- **Framing elbows.** Content sits inside rounded "elbow" frames — a thick
  coloured bar runs along an edge and bends 90° into a side rail. The bend is
  the signature shape.
- **Flat, saturated colour blocks.** No gradients on the chrome, no drop
  shadows. Colour is used categorically (a colour *means* something), not
  decoratively. The palette is fixed per theme (gold/butterscotch, reds/mars,
  blues/ice, african-violet, etc.) and exposed as CSS variables.
- **Dense, labelled readouts.** Numbers everywhere, in monospace-ish columns,
  often with cryptic alphanumeric tags. Information density is part of the
  aesthetic.
- **Rounded pill buttons and segmented bars.** Controls are stadium-shaped
  pills; meters are built from discrete segments.

The job of a dashboard built here is to take live Home Assistant data and
present it in this idiom: grouped into framed regions, colour-coded by meaning,
with the things you can *act on* clearly distinct from the things you only
*read*.

The class names and DOM structure of every component block reproduce the
original thelcars.com template so its theme stylesheet applies unmodified.
**Treat the markup as load-bearing** — never rename a class or restructure a
block's DOM to restyle it; override in `components.css` (our layer) instead.
`classic.css` is vendored byte-identical and never edited.

---

## 2. The information hierarchy

Design a page in two tiers: the **page frame** (tier 1) and the **sections**
inside its body (tier 2).

```
frame.html                        ← TIER 1: the page frame (chrome)
├─ banner                         ← identity — site / page name
├─ ctl-nav (under the banner)     ← quick controls for THIS page's system
├─ left panels (03-06…)           ← page-to-page navigation (sibling pages)
├─ main (dash.html)               ← TIER 2: columns of sections
│  ├─ bezel { bar "Group A" + content }   …primary section (elbow-framed)
│  │        { bar "Group B" + content }   …secondary heading inside a bezel
│  └─ bezel { … }
└─ footer                         ← status line (attribution — keep it)
```

### Tier 1 — the page frame

One per page; the generic `dash.html` renders every dashboard through it. The
`banner` answers "where am I"; the `footer` carries the required attribution.
Two parts of the frame are easy to confuse, so be deliberate:

- **The left panel column is navigation** — it links to *sibling pages* of the
  dashboard (`framePanels` marks the current page lit). This is the only place
  page-to-page navigation belongs.
- **The button array under the banner (`ctl-nav`) is for the current page, not
  navigation.** It holds **quick controls for *this* page's system** — the one
  or two actions you would reach for first on this screen (arm the alarm,
  exterior lights, run a speed test). They act on the page's own subsystem via
  the controls table; they are *not* links to other pages. Keep it to a few
  buttons — the long tail of controls belongs in `main`.

### Tier 2 — sections within main

The house rule for these dashboards: **primary sectioning uses bezels;
secondary sectioning uses bars.** Each top-level group on a page is an
elbow-framed bezel (`section` in the page table), and headings inside a bezel
are `lcars-text-bar` bars. Don't nest bezels; if a bezel needs more than a
couple of sub-headings, it should probably be two bezels.

The original guide's caution still applies in spirit: the emphasis a frame
carries comes from restraint. Prefer a handful of well-chosen bezels per page
over boxing every stray readout, and let bars carry the rhythm inside them.

### Using bezels to express grouping

The bezel classes (`top-bezel`/`bottom-bezel` + `bezel--left|right|both`) are
the elbow frame pieces: a coloured outer block with a black content area inset
by padding, so the colour shows as a top/bottom bar plus side rails.

- **Pick the side** with the rail pointing toward the structure it belongs to;
  left-railed frames read as "attached to the page's left spine" and are the
  default.
- **Colour carries meaning.** Set the section's colour to the category colour
  for the group (e.g. `var(--mars)` for security/warnings, `var(--ice)` for
  network). The inner area stays black.
- Geometry defaults (15px bar, 60px rail, 30px inner radius, outer = inner +
  bar) are tuned; prefer them for consistency.

---

## 3. Page anatomy (where things live)

A dashboard page is **data, not markup**: one entry in the `pages` table in
`internal/server/pages.go`.

- `pageDef` — slug, title, banner, `Controls` (the ctl-nav button ids), and
  `Cols`: columns of sections.
- `section` — a bezel: heading, `Kind` (`cells` | `bars` | `meters` |
  `cascade`), optional `PillCols` (1 = wide readout rows, 3 = dense arrays),
  and its `Elems`.
- `elem` — one live readout: an entity plus a `view` spec (`TileSpec`,
  `BarSpec`, `MeterSpec`) chosen by the constructor (`tile`/`bar`/`meter`).
- `controls` — the actionable buttons (see §5); pages reference them by id.

Routes, the SSE stream, and the fallback poll all derive from these tables.
Editing a page means editing `pages.go` only. The kiosk layout sizes every
page to fill a 1080p monitor with no scrolling (narrow screens fall back to
normal scrolling); check new sections fit at 1920×1080.

---

## 4. Component catalogue by purpose

### Structure & navigation
| Block / mechanism | Use for |
|---|---|
| `frame.html` | The page frame. Automatic for every page. |
| `lcars-bar` / heading bars | Titled divider introducing a content group inside a bezel. |
| bezel (`dash-sec`) | The elbow frame around each primary section. |
| left panels (`framePanels`) | Page-to-page navigation. |
| `ctl-nav` + `toggle` | The under-banner button row — quick controls for the current page's system, not navigation. |
| `pillbox` (`cells` kind) | A grid of pills; status readouts or (as buttons) quick toggles. |

### Read-only data displays
| Block | Best for |
|---|---|
| `cascade-cols` (`cascade` kind) | Dense, multi-entity telemetry — the animated number grid. The workhorse for "show me lots of sensor values at once." Cells carry per-cell formatting (and can carry threshold colouring). |
| `bar-row` (`bars` kind) | A handful of comparable values as horizontal bars (each with its own `Max`, `Unit`, `Color`). |
| `meter-block` (`meters` kind) | A value as a segmented vertical gauge (VU-meter style) — the house standard for single values against a range (temperatures, capacities). |
| `cell` (`cells` kind) | A labelled state readout pill (alarm state, HVAC mode, motion CLEAR/MOTION). |
| `image-frame` / `gallery` | A camera snapshot or image, framed; or a set of images. |

### Actionable controls
| Mechanism | Use for |
|---|---|
| `toggle` (controls table) | A single on/off-able entity as a stateful button; `hx-post` to `/action/{id}`, colour reflects state, SSE keeps it honest. |
| pill toggles | The same control rendered as a pill inside a pillbox — for arrays of quick toggles (all the lights). |
| `accordion` | Collapsible section to tuck away secondary controls/detail. |

---

## 5. Choosing a component from the Home Assistant entity

The single most important design decision is **read-only vs. actionable**, and
it follows almost entirely from the entity's domain (the part before the dot
in `switch.foo`).

### Rule of thumb

> A value is **read-only** if Home Assistant exposes no service that changes
> it. It is **actionable** if there is a service that sets/toggles/invokes it
> *and* the user should be able to trigger that from this screen.

Read-only entities must never be rendered as buttons — a button that can't do
anything is a lie. Conversely, an actionable entity rendered as plain text
wastes the control.

### Domain → treatment

**Read-only (display with cascade / meter / bar / cell):**

| Domain | Typical meaning | Suggested display |
|---|---|---|
| `sensor.*` | Numeric/text telemetry (temp, power, speed, counts) | cascade cell; `meter`/`bar` for a highlighted value |
| `binary_sensor.*` | on/off state you can't set (motion, door, presence) | `cell` with a `view.Map` (MOTION/CLEAR, OPEN/CLOSED) |
| `weather.*` | Forecast | text/cascade |
| `device_tracker.*`, `person.*` | home/away | status cell |
| `sun.*`, `update.*`, diagnostic sensors | state | cell / cascade |

**Actionable (render a control that calls a service):**

| Domain | Service to call | Suggested control |
|---|---|---|
| `switch.*` | `turn_on` / `turn_off` (toggle on current state) | toggle; colour reflects on/off |
| `light.*` | `turn_on` / `turn_off` | toggle (pill toggle in arrays) |
| `fan.*` | `turn_on`/`turn_off`, `set_percentage` | toggle + level control |
| `lock.*` | `lock.lock` / `lock.unlock` | toggle with explicit services (see note) |
| `cover.*` | `open_cover` / `close_cover` / `stop_cover` | up/down/stop buttons |
| `climate.*` | `set_hvac_mode`, `set_temperature` | on/off toggle now; mode-cycle later |
| `button.*` | `button.press` | momentary button |
| `script.*` / `scene.*` | `turn_on` | button |
| `input_boolean.*` | `toggle` | toggle |
| `input_number.*` / `number.*` | `set_value` | stepper/level control |
| `media_player.*` | `media_play_pause`, `volume_set`, … | transport buttons |

Controls live in the `controls` table; `handleAction` calls the service and
returns the updated button fragment.

### Notes & patterns

- **Always send the command unconditionally for state-ambiguous domains.** For
  a lock whose state may be `unknown`, call `lock.lock` regardless of the
  read-back state rather than guarding on `state === 'unlocked'` — the guard
  can wedge the control. Toggle domains (switch/light) are the exception: read
  current state to decide on vs off.
- **Reflect state in the control.** A toggle button looks different when its
  entity is on vs off (gold vs grey). The control both *acts* and *reads
  back* — SSE refreshes it when the state changes elsewhere (wall switch, HA
  app).
- **Mixed cluster.** A group often has both: read-only cells/cascade for the
  sensors plus toggles for the actionable entities of the same subsystem.
  Group them in one bezel.
- **High-level controls go up top.** Put this page's system's primary actions
  (the one or two things you'd reach for first) in the `ctl-nav` row; put the
  long tail of detail in `main`. The header is for acting on the current page,
  not for jumping to other pages — that's the left panels' job.

---

## 6. Home Assistant data

`ha.Source` is the store (`internal/ha`): `State(id)`, `All()`,
`CallService(domain, service, data, target)`, `Subscribe()` for change events.

- **Mock mode.** Standalone without HA credentials defaults to `MOCK=true`:
  any entity a page reads is fabricated and drifts. All layout work and the
  test suite run offline; live mode is reserved for the deployed container.
- **Numbers arrive as strings.** Parse via `view.Num(raw)`; it rejects
  `unknown`/`unavailable`/non-numeric input.
- **Formatting & thresholds** are first-class in the `view` package:
  `view.Fixed(n)`, `view.Map{...}` (state substitution), and the classifiers
  `view.Hot(limit)` / `view.Low(limit)` which return `view.AlertClass`
  (`font-red blink`) for out-of-range values. Wire them into `TileSpec`.
- **Unavailable states** render as `----` (`view.Placeholder`); mirror that
  convention in any custom display.
- **Live updates**: every element id is an SSE event name; the frame opens one
  `/sse` stream and htmx swaps fragments in place, with an infrequent
  `hx-get` poll as the reconciling fallback.

---

## 7. House idioms & gotchas

- **Element ids are identity.** `{slug}-{n}` ids are the HTML id, the SSE
  event name, and the poll path. An entity may appear on several pages with
  different presentations because events are keyed by element, not entity.
- **SSE fragments are single-line.** `FragmentHTML` flattens whitespace; keep
  fragment blocks free of significant whitespace (`<pre>` etc.).
- **Never reference by `url(#id)` in SVG.** `<base href>` re-roots fragment
  URLs, breaking gradients/filters on every non-base page. Precompute solid
  values server-side instead (see the level-meter).
- **The compat build is the product.** Pages are served with
  `*.compat.css` generated by `tools/cssgen` for ~2017 browsers. After any CSS
  change: `go run ./tools/cssgen`, and write `components.css` additions
  old-browser-safe where the passes can't infer context (e.g. sibling margins
  instead of flex `gap` when `display:flex` lives on another selector). The
  build fails if modern syntax leaks through.
- **Avoid the class name `active`** — the theme styles it globally (accordion
  behaviour).
- **Keep the footer attribution.** The template's terms require the
  TheLCARS.com credit line — don't remove it.
- **Verify at 1080p.** The kiosk layout must report
  `scrollHeight == 1080` at 1920×1080 (`tools/csscheck` has the Playwright
  setup; a quick page-fit probe lives in the dev scripts).

---

## 8. A minimal page skeleton

```go
{Slug: "climate", Title: "Climate", Banner: "EXAMPLE • CLIMATE",
    Controls: []string{"kitchen_hvac"},          // ctl-nav: this page's quick action
    Cols: [][]section{
        {{Heading: "Temperatures", Kind: "meters", Elems: []elem{
            meter("sensor.living_room_temp", roomTemp("Living")),
            meter("sensor.kitchen_temp", roomTemp("Kitchen")),
        }}},
        {{Heading: "Status", Kind: "cells", PillCols: 1, Elems: []elem{
            tile("climate.house", view.TileSpec{Label: "HVAC"}),
            tile("binary_sensor.window", view.TileSpec{Label: "Window", Format: opened}),
        }}},
    }},
```

The component demo at `/` is the living reference — every block is
demonstrated there. When you add a block or view spec to the engine, update
that page too.
