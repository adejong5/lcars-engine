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
   *   theme       — 'classic' | 'nemesis' | 'lower-decks'; controls the
   *                 wrapper class
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
   *   children — main content
   *   footer   — replaces the default footer text (attribution kept)
   */

  import SidePanels from './SidePanels.svelte';
  import { playBeep, playBeepAndGo } from '../sounds.js';

  /** @type {{ theme?: 'classic'|'nemesis'|'lower-decks', banner?: string, topPanel?: {label: string, href?: string}, panel2?: {label: string, hop?: string}, sounds?: boolean, floorText?: string, topFrame?: import('svelte').Snippet, sidebar?: import('svelte').Snippet, children?: import('svelte').Snippet, footer?: import('svelte').Snippet }} */
  let {
    theme = 'classic',
    banner = 'LCARS • 47988',
    topPanel = {},
    panel2 = { label: '02', hop: '-262000' },
    sounds = true,
    floorText,
    topFrame,
    sidebar,
    children,
    footer
  } = $props();

  const isLowerDecks = theme === 'lower-decks';

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

<svelte:element this={isLowerDecks ? 'div' : 'section'}
  class={isLowerDecks ? 'wrap-all' : 'wrap-standard'}
  id={isLowerDecks ? undefined : 'column-3'}>

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
        <div class="bar-1"></div>
        <div class="bar-2"></div>
        <div class="bar-3"></div>
        <div class="bar-4"></div>
        <div class="bar-5"></div>
      </div>
    </div>
  </div>

  <div class="wrap" id="gap">
    <div class="left-frame">
      <button onclick={toTop} id="topBtn"
        style:display={scrollY > 200 ? 'block' : 'none'}>
        <span class="hop">screen</span> top
      </button>
      {#if sidebar}
        {@render sidebar()}
      {:else}
        <SidePanels {sounds} />
      {/if}
    </div>
    <div class="right-frame">
      <div class="bar-panel">
        <div class="bar-6"></div>
        <div class="bar-7"></div>
        <div class="bar-8"></div>
        <div class="bar-9"></div>
        <div class="bar-10"></div>
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

</svelte:element>

<div class="headtrim"></div>
<div class="baseboard"></div>

<style>
  /* The theme CSS lays out <body> with display:flex and expects the page
     wrapper as a direct flex child; make the mount node transparent. */
  :global(#app) {
    display: contents;
  }
</style>
