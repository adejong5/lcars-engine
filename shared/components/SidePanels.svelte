<script>
  /**
   * The colored panel stack in the left sidebar, placed via LCARSPage's
   * sidebar snippet (LCARSPage renders a default <SidePanels /> when the
   * snippet is omitted).
   *
   * The theme CSS gives each panel-N class a designed height (panel-9 is
   * the fluid one that absorbs leftover page height, panel-10 pins to the
   * bottom via the left-frame's space-between). Panels cycle through
   * classes panel-3 … panel-9; the bottom panel is always panel-10.
   *
   * Props:
   *   panels — a number (that many auto-labeled filler panels, default 7
   *            to match the original template) or an array of
   *            { label?, hop?, href? }:
   *              label — big text; auto-numbered ('03'…) when omitted
   *              hop   — small suffix; auto-generated digits when omitted,
   *                      '' to suppress
   *              href  — renders the panel as a link (with beep)
   *   bottomPanel — same shape for the pinned panel-10 block; true for
   *                 auto filler (default), false to omit
   *   sounds — false to disable the beep on link panels
   */

  import { playBeepAndGo } from '../sounds.js';

  /** @type {{ panels?: number | Array<{label?: string, hop?: string, href?: string}>, bottomPanel?: boolean | {label?: string, hop?: string, href?: string}, sounds?: boolean }} */
  let { panels = 7, bottomPanel = true, sounds = true } = $props();

  function randomHop() {
    let s = '-';
    for (let i = 0; i < 6; i++) s += Math.floor(Math.random() * 10);
    return s;
  }

  function fill(p, i) {
    return {
      label: p.label ?? String(i + 3).padStart(2, '0'),
      hop: p.hop ?? randomHop(),
      href: p.href,
      cls: `panel-${3 + (i % 7)}`
    };
  }

  const topPanels = (typeof panels === 'number'
    ? Array.from({ length: panels }, () => ({}))
    : panels
  ).map(fill);

  const bottom = bottomPanel
    ? {
        label: bottomPanel.label ?? '10',
        hop: bottomPanel.hop ?? randomHop(),
        href: bottomPanel.href,
        cls: 'panel-10'
      }
    : null;

  function click(e, href) {
    if (!sounds) return;
    e.preventDefault();
    playBeepAndGo(2, href);
  }
</script>

{#snippet panel(p)}
  {#if p.href}
    <a class={p.cls} href={p.href} onclick={e => click(e, p.href)}>
      {p.label}{#if p.hop}<span class="hop">{p.hop}</span>{/if}
    </a>
  {:else}
    <div class={p.cls}>
      {p.label}{#if p.hop}<span class="hop">{p.hop}</span>{/if}
    </div>
  {/if}
{/snippet}

<div>
  {#each topPanels as p}
    {@render panel(p)}
  {/each}
</div>
<div>
  {#if bottom}
    {@render panel(bottom)}
  {/if}
</div>

<style>
  /* The theme styles panel-N divs; links need block layout for the
     class heights/padding to apply the same way. */
  a {
    display: block;
  }
</style>
