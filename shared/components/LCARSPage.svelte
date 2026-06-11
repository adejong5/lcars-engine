<script>
  /**
   * LCARS page frame, structured after the thelcars.com standard layout
   * (classic / nemesis-blue / lower-decks themes — import the matching
   * theme CSS in your site's main.js).
   *
   * The HTML structure and class names match the original template so the
   * theme stylesheets apply unmodified:
   *
   *   wrap (top)    — left-frame-top (panel-1, panel-2)
   *                   right-frame-top (banner, topFrame snippet, bar-panel 1-5)
   *   wrap#gap      — left-frame  (topBtn, panels 3-9, panel-10)
   *                   right-frame (bar-panel 6-10, main, footer)
   *
   * Props:
   *   theme       — 'classic' | 'nemesis' | 'lower-decks' | 'padd';
   *                 controls wrapper class, bar-panel segment count, the
   *                 padd divider block, and sidebar defaults. Import the
   *                 matching CSS in main.js (padd = lower-decks-padd.css).
   *   layout      — 'standard' (default) | 'ultra'. Ultra wraps the page
   *                 in the three-column wrap-everything layout from the
   *                 original ultra templates, adding decorative column-1
   *                 and column-2 to the left (classic and nemesis only —
   *                 the lower-decks CSS has no ultra classes).
   *   column1,    — pick ultra columns individually: true renders the
   *   column2       default decorative content, a snippet renders your
   *                 own. Setting either (or both) enables the ultra
   *                 frame without needing layout="ultra"; layout="ultra"
   *                 is shorthand for both columns with defaults.
   *   banner      — top banner text
   *   topPanel    — { label?, href? } for the big top-left panel button;
   *                 label defaults to 'LCARS', omit href for a decorative
   *                 button that just beeps
   *   panel2      — { label, hop } small panel under it
   *   sounds      — false to disable beeps
   *   floorText   — extra footer line (nemesis theme styles this)
   *
   * Snippets:
   *   topFrame — content for the top frame area, typically <DataCascade>
   *              and/or <NavButtons>
   *   sidebar  — content for the left sidebar, typically <SidePanels>;
   *              defaults to <SidePanels /> (7 auto-filler panels + bottom)
   *   column1  — far-left ultra column; defaults to the template's
   *              <LCARSFrame> + <Pillbox> + status list + <Pillbox
   *              variant={2}>
   *   column2  — narrow second ultra column; defaults to the template's
   *              panels 11–15 with three sidebar buttons
   *   children — main content
   *   footer   — replaces the default footer text (attribution kept)
   */

  import SidePanels from './SidePanels.svelte';
  import LCARSFrame from './LCARSFrame.svelte';
  import Pillbox from './Pillbox.svelte';
  import { playBeep, playBeepAndGo } from '../sounds.js';

  /** @type {{ theme?: 'classic'|'nemesis'|'lower-decks'|'padd', layout?: 'standard'|'ultra', banner?: string, topPanel?: {label: string, href?: string}, panel2?: {label: string, hop?: string}, sounds?: boolean, floorText?: string, topFrame?: import('svelte').Snippet, sidebar?: import('svelte').Snippet, column1?: boolean | import('svelte').Snippet, column2?: boolean | import('svelte').Snippet, children?: import('svelte').Snippet, footer?: import('svelte').Snippet }} */
  let {
    theme = 'classic',
    layout = 'standard',
    banner = 'LCARS • 47988',
    topPanel = {},
    panel2 = { label: '02', hop: '-262000' },
    sounds = true,
    floorText,
    topFrame,
    sidebar,
    column1,
    column2,
    children,
    footer
  } = $props();

  // reactive so a site can switch layouts at runtime; either column prop
  // (true = defaults, snippet = custom) enables the ultra frame on its own
  const ultra = $derived(layout === 'ultra' || !!column1 || !!column2);
  const showCol1 = $derived(layout === 'ultra' || !!column1);
  const showCol2 = $derived(layout === 'ultra' || !!column2);
  // default ultra column-2 buttons; labels/colors from the original templates
  const sideButtons = theme === 'nemesis'
    ? [['JS2B-01', 'button-evening'], ['JS2B-02', 'button-moonbeam'], ['MS2B-03', 'button-evening']]
    : [['JS2B-01', 'button-almond-creme'], ['JS2B-02', 'button-butterscotch'], ['MS2B-03', 'button-african-violet']];

  const isPadd = theme === 'padd';
  // lower-decks and padd use the wrap-all wrapper
  const wrapAll = theme === 'lower-decks' || isPadd;
  // lower-decks bar panels have 4 segments; classic/nemesis/padd have 5
  const topBars = theme === 'lower-decks' ? [1, 2, 3, 4] : [1, 2, 3, 4, 5];
  const bottomBars = theme === 'lower-decks' ? [6, 7, 8, 9] : [6, 7, 8, 9, 10];

  let scrollY = $state(0);

  function nav(href) {
    if (sounds) playBeepAndGo(2, href);
    else if (href && href !== '#') window.location.href = href;
  }

  function toTop() {
    if (sounds) playBeep(4);
    document.documentElement.scrollTop = 0;
    document.body.scrollTop = 0;
  }
</script>

<svelte:window bind:scrollY />

{#snippet pageFrame()}
  <div class="wrap">
    <div class="left-frame-top">
      <button onclick={() => nav(topPanel.href)} class="panel-1-button">{topPanel.label ?? 'LCARS'}</button>
      <div class="panel-2">{panel2.label}<span class="hop">{panel2.hop}</span></div>
    </div>
    <div class="right-frame-top">
      <div class="banner">{banner}</div>
      <div class="data-cascade-button-group">
        {@render topFrame?.()}
      </div>
      <div class="bar-panel first-bar-panel">
        {#each topBars as n}
          <div class="bar-{n}"></div>
        {/each}
      </div>
    </div>
  </div>

  {#if isPadd}
    <div class="divider">
      <div class="block-left"></div>
      <div class="block-right">
        <div class="block-row">
          <div class="bar-11"></div>
          <div class="bar-12"></div>
          <div class="bar-13"></div>
          <div class="bar-14">
            <div class="blockhead"></div>
          </div>
        </div>
      </div>
    </div>
  {/if}

  <div class="wrap" id="gap">
    <div class="left-frame">
      <button onclick={toTop} id="topBtn"
        style:display={scrollY > 200 ? 'block' : 'none'}>
        <span class="hop">screen</span> top
      </button>
      {#if sidebar}
        {@render sidebar()}
      {:else}
        <SidePanels {theme} {sounds} />
      {/if}
    </div>
    <div class="right-frame">
      <div class="bar-panel">
        {#each bottomBars as n}
          <div class="bar-{n}"></div>
        {/each}
      </div>
      <main>
        {@render children?.()}
      </main>
      {#if floorText}
        <div class="floor-text">{floorText}</div>
      {/if}
      <footer>
        {#if footer}
          {@render footer()}
        {/if}
        <!-- Required attribution for the theme CSS: -->
        LCARS Inspired Website Template by <a href="https://www.thelcars.com">www.TheLCARS.com</a>.
      </footer>
    </div>
  </div>
{/snippet}

{#if ultra}
  <div class="wrap-everything">
    {#if showCol1}
    <section id="column-1">
      {#if typeof column1 === 'function'}
        {@render column1()}
      {:else}
        <LCARSFrame horizontal={theme === 'nemesis'} />
        <Pillbox {sounds} />
        <div class="lcars-list-2 uppercase">
          <ul>
            <li>Subspace Link: Established</li>
            <li>Starfleet Database: Connected</li>
            <li>Quantum Memory Field: stable</li>
            <li class={theme === 'nemesis'
              ? 'bullet-moonbeam font-moonbeam'
              : 'bullet-almond-creme font-almond-creme'}>Optical Data Network: rerouting</li>
          </ul>
        </div>
        <Pillbox variant={2} {sounds} pills={[
          { label: 'F12-22' }, { label: 'G24-22' }, { spacer: true },
          { label: 'H-07AM' }, { label: 'I50-72' }, { label: 'J5369' }
        ]} />
      {/if}
    </section>
    {/if}
    {#if showCol2}
    <section id="column-2">
      {#if typeof column2 === 'function'}
        {@render column2()}
      {:else}
        <div class="panel-11">11-1524</div>
        {#each sideButtons as [label, color]}
          <button onclick={() => nav()} class="sidebar-button {color}">{label}</button>
        {/each}
        <div class="panel-12">12-0730</div>
        <div class="panel-13">13-318</div>
        <div class="panel-14">14-DL44</div>
        <div class="panel-15">15-3504</div>
      {/if}
    </section>
    {/if}
    <section id="column-3">
      {@render pageFrame()}
    </section>
  </div>
{:else}
  <svelte:element this={wrapAll ? 'div' : 'section'}
    class={wrapAll ? 'wrap-all' : 'wrap-standard'}
    id={wrapAll ? undefined : 'column-3'}>
    {@render pageFrame()}
  </svelte:element>
{/if}

<div class="headtrim"></div>
<div class="baseboard"></div>

<style>
  /* The theme CSS lays out <body> with display:flex and expects the page
     wrapper as a direct flex child; make the mount node transparent. */
  :global(#app) {
    display: contents;
  }
</style>
