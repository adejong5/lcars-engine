# ENGINE.md — Building LCARS dashboards with this component library

This is a guide for an AI assistant assembling dashboard sites from the
components in `shared/`. It assumes you have already read `CLAUDE.md` (build
commands, the public/private two-repo split, theme system). This document is
about *how to design a site*: structuring information, picking components for
each kind of data, and knowing when a value should be a read-only display
versus an interactive control.

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

The class names and DOM structure of every component reproduce the original
thelcars.com template so its theme stylesheets apply unmodified. **Treat the
markup as load-bearing** — never rename a class or restructure a component's
DOM to restyle it; override via the documented props and theme CSS variables
instead.

---

## 2. The information hierarchy

Design a page in two tiers: the **page layout** (tier 1) and the **sections**
inside its body (tier 2).

```
LCARSPage                         ← TIER 1: the page layout (frame + chrome)
├─ banner                         ← identity — site / section name
├─ topFrame                       ← quick controls for THIS page's system + headline telemetry
├─ sidebar (SidePanels)           ← page-to-page navigation (sibling pages)
├─ column1/column2 (ultra)        ← at-a-glance summary of THIS page
├─ main (children)                ← TIER 2: a sequence of sections, each delimited by…
│  ├─ LCARSBar "Group A" + content    …a heading bar (the default), OR
│  ├─ Bezel { Group B content }       …an elbow frame (to compartmentalize / emphasize)
│  └─ LCARSBar "Group C" + content
└─ footer                         ← status line (data link, attribution)
```

### Tier 1 — the page layout (`LCARSPage`)

One per page; it never nests. The `banner` answers "where am I"; the `footer`
carries the live connection status and the required attribution. Two parts of
the frame are easy to confuse, so be deliberate:

- **`sidebar` (SidePanels) is navigation** — it links to *sibling pages* of the
  dashboard. This is the only place page-to-page navigation belongs.
- **`topFrame` is for the current page, not navigation.** The header button row
  holds **quick controls for *this* page's system** — the one or two actions
  you would reach for first on this screen (arm the alarm, all-lights-off, run
  a speed test, restart the server). They act on the page's own subsystem and
  call `ha.callService`; they are *not* links to other pages.

### Tier 2 — sections within `main`

The body is a flat sequence of sections. A section is introduced **either** by
a heading bar **or** by a bezel frame — these are **peers**, two
interchangeable ways to delimit a section, not nested levels. Choose per
section, and mix them freely down the page.

1. **Heading bar — `LCARSBar`.** The default, lightweight section divider — the
   equivalent of an `<h2>`/`<h3>` (use the `heading` prop for depth). This is
   the normal rhythm of the page: a title, then its `DataCascade` / gauges /
   controls below it.

2. **Bezel frame — `TopBezel` / `MiddleBezel` / `BottomBezel`.** The
   alternative: instead of (or occasionally in addition to) a heading, box the
   section's content in an elbow frame. Use a bezel to **compartmentalize or
   emphasize a subset of entities** — to make a cluster read as one contained,
   visually distinct unit (a status callout, a grouped control cluster, a
   framed summary) rather than just another bar-headed block. A bezeled
   section sits at the same tier as a bar-headed one; it is not subordinate to
   the bar above it.

### Using bezels to express grouping

`TopBezel`, `MiddleBezel`, and `BottomBezel` are the elbow frame pieces. Each
is a coloured outer block with a black content area inset by padding, so the
colour shows as a top/bottom bar plus one or both side rails.

- **Cap a group** with a single `TopBezel` (bar on top, content hangs below)
  or `BottomBezel` (content sits on a base bar) when you want one decorated
  edge.
- **Box a group fully** by stacking `TopBezel` → `MiddleBezel` (straight side
  rails, no bar) → `BottomBezel`. They stack flush; the elbow opens around the
  shared content column.
- **Pick the side** with `sides="left" | "right" | "both"`. Convention: the
  rail points toward the navigation/structure it belongs to. Left-railed
  frames are the default and read as "attached to the page's left spine".
- **Colour carries meaning.** Set `color` to the category colour for the group
  (e.g. `var(--mars)` for a warnings block, `var(--ice)` for network). The
  inner area defaults to black; keep it black unless you have a reason.

Defaults are tuned so the elbow looks right at small sizes: 15px top/bottom
bar, 60px side rails, 30px inner radius, and an **outer radius of
`innerRadius + bar`** so the coloured band keeps a uniform thickness through
the bend. Override `bar`, `side`, `innerRadius`, or `outerRadius` per instance
if you need a chunkier or tighter frame, but prefer the defaults for
consistency across a site.

> Bars and bezels are peers — don't box *every* section, or the emphasis a
> bezel is supposed to carry disappears. Make `LCARSBar` headings the normal
> rhythm of the page and switch to a bezel for the few sections that genuinely
> deserve to be compartmentalized or emphasized.

### The "ultra" columns — page-level summary

`LCARSPage` supports an optional far-left column via the `column1` (and
narrower `column2`) snippet/prop. Passing either switches the page into the
three-column ultra layout. Use `column1` for an **at-a-glance summary of the
current page**: typically an `LCARSFrame` holding a `BarChart`/`LevelMeter` of
the page's headline metric, plus a `Pillbox` of key entity states. It is
hidden on narrow viewports, so never put anything essential *only* there.

---

## 3. Site anatomy (the three files)

Every site is three files in `sites/<name>/` (demos) or `private/<name>/`
(real dashboards):

- **`index.html`** — boilerplate; a `<div id="app">` and a module script
  loading `main.js`.
- **`main.js`** — mounts `App.svelte`, imports exactly one theme CSS, imports
  the MDI icon font if you use icons, and calls `initHA(...)` once.
- **`App.svelte`** — the page: one `<LCARSPage>` with its snippets filled in.

The theme is chosen in **two coordinated places** and they must agree:

```js
// main.js
import '$shared/assets/classic.css';        // 1. the stylesheet
```
```svelte
<!-- App.svelte -->
<LCARSPage theme="classic">                 <!-- 2. the matching prop -->
```

Valid pairs: `classic.css`+`theme="classic"`, `nemesis-blue.css`+`"nemesis"`,
`lower-decks.css`+`"lower-decks"`, `lower-decks-padd.css`+`"padd"`. The prop
controls bar-segment counts, the padd divider, and sidebar defaults; the CSS
controls the actual styling.

`initHA` connects the Home Assistant store; see §6.

---

## 4. Component catalogue by purpose

### Structure & navigation
| Component | Use for |
|---|---|
| `LCARSPage` | The page frame. One per page. Snippets: `topFrame`, `sidebar`, `column1`/`column2`, `children`, `footer`. |
| `LCARSBar` | Titled divider introducing a content group (`heading="h2"…"h4"`), or a plain rule with no `title`. |
| `TopBezel` / `MiddleBezel` / `BottomBezel` | Elbow frames to compartmentalize/emphasize a section (a peer to an `LCARSBar` heading, not a sub-level). |
| `SidePanels` | The left sidebar panel stack — page-to-page navigation. Each panel can take an `href`. |
| `NavButtons` | The header button row (placed in `topFrame`) — quick controls for the current page's system, not page navigation. |
| `Pillbox` | A 2-column grid of pill buttons; good for a cluster of quick links or status pills. |

### Read-only data displays
| Component | Best for |
|---|---|
| `DataCascade` | Dense, multi-entity telemetry — the animated number grid. The workhorse for "show me lots of sensor values at once." Carries per-cell formatting and threshold-colouring. |
| `PillGauge` (house-specific, in `private/live/`) | A single value against a range, with a gradient fill whose leading edge colours by level. Great for %, temperature, speed. |
| `BarChart` | A handful of comparable values as horizontal bars (each with its own `max`, `unit`, `color`). Good inside an `LCARSFrame` / `column1`. |
| `LevelMeter` | A single value as a segmented vertical gauge (VU-meter style). Good for a tall, glanceable indicator. |
| `LCARSFrame` | A decorative bordered display cell; drop a chart/meter into its centre, or let it animate its default "audio level" bars as pure decoration. |
| `ImageFrame` / `Gallery` | A camera snapshot or image, framed; or a set of images. |

### Actionable controls
| Component | Use for |
|---|---|
| `LCARSButton` | A single action or link. `onclick` (e.g. an `ha.callService`) runs before the beep/navigation. Variants: default pill, `sidebar`, `panel`. |
| `ButtonRow` | Layout wrapper to group several `LCARSButton`s with a `justify` mode. |
| `Pillbox` pills | Quick toggles/links in a compact grid (give each pill an `onclick` or `href`). |
| `Accordion` | Collapsible section to tuck away secondary controls/detail. |

---

## 5. Choosing a component from the Home Assistant entity

The single most important design decision is **read-only vs. actionable**, and
it follows almost entirely from the entity's domain (the part before the dot in
`switch.foo`).

### Rule of thumb

> A value is **read-only** if Home Assistant exposes no service that changes
> it. It is **actionable** if there is a service that sets/toggles/invokes it
> *and* the user should be able to trigger that from this screen.

Read-only entities must never be rendered as buttons — a button that can't do
anything is a lie. Conversely, an actionable entity rendered as plain text
wastes the control.

### Domain → treatment

**Read-only (display with cascade / gauge / meter / text):**

| Domain | Typical meaning | Suggested display |
|---|---|---|
| `sensor.*` | Numeric/text telemetry (temp, power, speed, counts) | `DataCascade` cell; `PillGauge`/`BarChart`/`LevelMeter` for a single highlighted value |
| `binary_sensor.*` | on/off state you can't set (motion, door, presence) | `DataCascade` cell with threshold class, or a status pill (non-button) / coloured text |
| `weather.*` | Forecast attributes | read attributes via `ha.entity(id).attributes`, show as text/cascade |
| `device_tracker.*`, `person.*` | home/away | status text or pill |
| `sun.*`, `update.*`, diagnostic sensors | state | cascade / text |

**Actionable (render a control that calls a service):**

| Domain | Service to call | Suggested control |
|---|---|---|
| `switch.*` | `switch.turn_on` / `turn_off` (toggle on current state) | `LCARSButton` or pill; colour reflects on/off |
| `light.*` | `light.toggle` / `turn_on` (+`brightness`) | button; optionally a brightness control |
| `fan.*` | `fan.turn_on`/`turn_off`, `fan.set_percentage` | button + level control |
| `lock.*` | `lock.lock` / `lock.unlock` | button (see note below) |
| `cover.*` | `cover.open_cover` / `close_cover` / `stop_cover` | up/down/stop buttons |
| `climate.*` | `climate.set_hvac_mode`, `climate.set_temperature` | mode-cycle buttons + setpoint control; state = current hvac mode |
| `button.*` | `button.press` | single momentary button |
| `script.*` | `script.turn_on` (or `script.<name>`) | button |
| `scene.*` | `scene.turn_on` | button |
| `input_boolean.*` | `input_boolean.toggle` | toggle button |
| `input_number.*` / `number.*` | `*.set_value` | stepper/level control |
| `input_select.*` / `select.*` | `*.select_option` | option-cycle button |
| `media_player.*` | `media_player.media_play_pause`, `volume_set`, … | transport buttons |

Call services through the store:

```js
ha.callService('switch', 'turn_on', {}, { entity_id: 'switch.adguard' });
ha.callService('button', 'press', {}, { entity_id: 'button.ha_restart' });
ha.callService('climate', 'set_hvac_mode', { hvac_mode: 'cool' },
               { entity_id: 'climate.kitchen' });
```

### Notes & patterns

- **Always send the command unconditionally for state-ambiguous domains.** For
  a lock whose state may be `unknown` (never reported), call `lock.lock`
  regardless of the read-back state rather than guarding on
  `state === 'unlocked'` — the guard can wedge the control. Toggle domains
  (switch/light) are the exception: read current state to decide on vs off.
- **Reflect state in the control.** A toggle button should look different when
  its entity is on vs off (colour, icon, a `.pill-off` opacity class). The
  control both *acts* and *reads back*.
- **Mixed cluster.** A group often has both: a `DataCascade` of the
  read-only sensors plus a `Pillbox`/`ButtonRow` of the actionable ones for
  the same subsystem. Group them under one `LCARSBar` (or in one bezel).
- **High-level controls go up top.** Put this page's system's primary actions
  (the one or two things you'd reach for first) in the `topFrame` header button
  row; put the long tail of detail in `main`. The header is for acting on the
  current page, not for jumping to other pages — that's the sidebar's job.

---

## 6. Home Assistant data

The reactive store lives in `shared/ha.svelte.js`. In `main.js` call `initHA`
once; in components read from `ha`:

```js
import { ha } from '$shared/ha.svelte.js';

ha.state('sensor.temp')    // reactive state string, or undefined until known
ha.entity('sensor.temp')   // reactive full object: { state, attributes, ... }
ha.connected               // reactive boolean — drive the footer status line
ha.callService(domain, service, data, target)   // returns a Promise
```

Reactivity is per-entity (a `SvelteMap`), so reading `ha.state(id)` inside a
template or `$derived` re-runs only when *that* entity changes.

- **Mock mode.** Without `.env.local` credentials (or with `?mock` in the
  URL), `initHA({ mock: true })` fabricates any entity a component reads and
  drifts it randomly. This means **all layout work happens offline** — every
  `ha.state(...)` you reference just appears as plausible data. `binary_sensor.*`
  mocks as on/off; everything else as a drifting number.
- **Numbers arrive as strings.** `ha.state` returns the raw string. Parse
  before maths: `const n = parseFloat(ha.state(id))`. Guard with
  `Number.isFinite`.
- **Formatting & thresholds** are first-class in `DataCascade`: each value can
  be `{ entity, format, cls }` where `format(state)` returns the display text
  and `cls(state)` returns extra classes (e.g. `'font-red blink'` above a
  limit). `{ text }` makes a static label cell. See `private/*/App.svelte`
  for worked examples.
- **Unavailable states.** `DataCascade` renders `'unavailable'`/`'unknown'` as
  `----`; mirror that convention in custom displays (the gauges show `----`
  for non-finite values).

---

## 7. Svelte 5 idioms & gotchas

The codebase is Svelte 5 throughout; match it.

- **Runes, not legacy.** `$props()`, `$state()`, `$derived()`,
  `$derived.by()`, `$effect()`. Snippets (`{#snippet}` / `{@render}`), never
  slots.
- **`$derived` vs `$derived.by`.** `$derived(expr)` takes an expression;
  `$derived(() => {...})` makes the *value* a function (a common bug). For a
  multi-line computation use `$derived.by(() => { ...; return x })`.
- **`{@const}` placement.** It must be an immediate child of a block
  (`{#snippet}`, `{#each}`, `{#if}`, `{:else}`) — *not* inside a plain
  `<div>`. Hoist `{@const}` declarations to the top of the snippet/block.
- **Theme override specificity.** Component CSS bundles *before* the theme
  stylesheet, so an equal-specificity override loses. In a site's `<style>`,
  add an extra ancestor to win (e.g. `main .lcars-text-bar { ... }`), and use
  `:global(...)` to reach into component DOM.
- **Avoid the class name `active`** — the theme styles it globally (accordion
  behaviour).
- **Keep the footer attribution.** The template's terms require the
  TheLCARS.com credit line; `LCARSPage` renders it in the footer — don't
  remove it.
- **Sounds.** `playBeep(n)` (n=1–4) for a beep; `playBeepAndGo(n, url)` to
  beep then navigate. Most navigation components beep already.

---

## 8. A minimal page skeleton

```svelte
<script>
  import LCARSPage from '$shared/components/LCARSPage.svelte';
  import LCARSBar from '$shared/components/LCARSBar.svelte';
  import DataCascade from '$shared/components/DataCascade.svelte';
  import ButtonRow from '$shared/components/ButtonRow.svelte';
  import LCARSButton from '$shared/components/LCARSButton.svelte';
  import SidePanels from '$shared/components/SidePanels.svelte';
  import { ha } from '$shared/ha.svelte.js';

  const num = id => parseFloat(ha.state(id));
  const hot = limit => v => parseFloat(v) >= limit ? 'font-red blink' : '';
</script>

<LCARSPage theme="classic" banner="HIGHBROOK • CLIMATE" topPanel={{ label: 'OPS', href: '../live/' }}>
  {#snippet sidebar()}
    <SidePanels panels={[{ label: '03', hop: '-OPS', href: '../live/' }]} />
  {/snippet}

  <!-- main content -->
  <LCARSBar title="Temperatures" heading="h3" />
  <DataCascade values={[
    { text: 'LIVING' },
    { entity: 'sensor.living_room_temp', format: v => parseFloat(v).toFixed(1) + '°', cls: hot(78) },
    { text: 'KITCHEN' },
    { entity: 'sensor.kitchen_temp', format: v => parseFloat(v).toFixed(1) + '°', cls: hot(78) }
  ]} filler={2} pattern={[1, 2]} />

  <LCARSBar title="Controls" heading="h3" />
  <ButtonRow>
    <LCARSButton label="Cool" color="ice"
      onclick={() => ha.callService('climate', 'set_hvac_mode', { hvac_mode: 'cool' }, { entity_id: 'climate.house' })} />
    <LCARSButton label="Off" color="mars"
      onclick={() => ha.callService('climate', 'set_hvac_mode', { hvac_mode: 'off' }, { entity_id: 'climate.house' })} />
  </ButtonRow>

  {#snippet footer()}
    Highbrook Climate &#149;
    <span class={ha.connected ? 'font-gold' : 'font-red blink'}>{ha.connected ? 'ONLINE' : 'OFFLINE'}</span> &#149;
  {/snippet}
</LCARSPage>
```

`sites/test/App.svelte` is the living reference — every component is
demonstrated there in the section that documents its props. When you add a
component or prop to the engine, update that page too.
