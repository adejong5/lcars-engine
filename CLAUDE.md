# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- `npm run dev` — start the Vite dev server (opens `/sites/test/`, listens on all interfaces)
- `npm run build` — build all sites to `dist/`
- `npm run preview` — serve the production build

There are no tests or linters configured.

## What this is

A Svelte 5 + Vite re-implementation of the thelcars.com LCARS website template, used to build multiple small dashboard sites (Star Trek LCARS styling) with live data from Home Assistant.

## Architecture

### Multi-site MPA

Each site lives in `sites/<name>/` with three files: `index.html`, `main.js` (mounts the app, picks a theme CSS, calls `initHA`), and `App.svelte`. **Creating a new site requires registering its `index.html` in `build.rollupOptions.input` in `vite.config.js`**, and optionally adding a link in the root `index.html` (the site directory page).

`$shared` is a Vite alias for `./shared` (also mapped in `jsconfig.json`).

### Theme system — class names are load-bearing

Components in `shared/components/` reproduce the exact HTML structure and class names of the original thelcars.com template so its theme stylesheets apply unmodified. The CSS files in `shared/assets/` (`classic.css`, `nemesis-blue.css`, `lower-decks.css`, `lower-decks-padd.css`) are taken from that template — **do not rename classes in components or restructure their markup without checking the theme CSS**.

A site picks a theme in two coordinated places:
1. `main.js` imports the theme CSS (e.g. `import '$shared/assets/classic.css'`)
2. `App.svelte` passes the matching `theme` prop to `<LCARSPage>` (`'classic' | 'nemesis' | 'lower-decks' | 'padd'`; padd uses `lower-decks-padd.css`)

`LCARSPage.svelte` is the page frame (banner, side panels, bar segments, footer). Content goes in its `topFrame`, `sidebar`, `children`, and `footer` snippets. It also takes `layout="ultra"` (classic/nemesis themes only) for the original three-column ultra layout, with `column1`/`column2` snippets overriding the decorative side columns. The footer attribution to TheLCARS.com is required by the template's terms — keep it.

`sites/test/App.svelte` is a living documentation page: every component is demonstrated in the section that documents its props. Keep it updated when adding components or props.

### Home Assistant integration

`shared/ha.svelte.js` is a reactive store built on Svelte 5 runes and `SvelteMap` (fine-grained: per-entity reactivity). Sites call `initHA(...)` once in `main.js`; components read `ha.state(id)` / `ha.entity(id)` / `ha.connected` and call `ha.callService(...)`.

- Real connection config comes from `.env.local` (`VITE_HA_HOST`, `VITE_HA_TOKEN`), which is gitignored.
- Without those env vars, sites fall back to **mock mode**: any entity a component reads is lazily fabricated and randomly drifts, so all UI development works offline.
- `shared/ha-websocket.js` is the underlying plain-JS WebSocket client (auth, subscriptions, reconnect).

### Other shared modules

- `shared/sounds.js` — LCARS beep effects (`playBeep`, `playBeepAndGo` for beep-then-navigate).
- `online-template/` — the original downloaded thelcars.com HTML templates, kept locally as reference only (gitignored, not part of the build).
- `template/` — an older hand-rolled CSS prototype with a different class scheme; not used by the Svelte sites.

## Conventions

- Svelte 5 idioms throughout: `$props()`, `$state()`, snippets (`{#snippet}` / `{@render}`) — not slots or legacy reactive statements.
- Components are JSDoc-typed plain JS (no TypeScript).
