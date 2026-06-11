<script>
  /**
   * A single LCARS button/link for the content area.
   *
   * Variants (matching the theme CSS):
   *   default   — the big rounded pill (.lcars-button / .buttons style)
   *   'sidebar' — flat block styled like a sidebar panel (.sidebar-button)
   *   'panel'   — .panel-button; combine with `pan` to mimic an established
   *               panel's height/color (.pan-0/4/5/6/7/8/10)
   *
   * Props:
   *   label   — button text (or use children)
   *   href    — renders an <a>; omit for a <button> that just beeps
   *   color   — theme color name, applied as button-{color}
   *             (e.g. 'orange', 'african-violet', 'butterscotch')
   *   variant — see above
   *   pan     — panel number for variant 'panel'
   *   sounds  — false to disable the beep
   *   onclick — click handler (e.g. an ha.callService call); runs before
   *             the beep/navigation
   */

  import { playBeepAndGo } from '../sounds.js';

  /** @type {{ label?: string, href?: string, color?: string, variant?: 'sidebar'|'panel', pan?: number, sounds?: boolean, onclick?: (e: MouseEvent) => void, children?: import('svelte').Snippet }} */
  let { label, href, color, variant, pan, sounds = true, onclick, children } = $props();

  const cls = [
    variant === 'sidebar' ? 'sidebar-button' : variant === 'panel' ? 'panel-button' : 'lcars-button',
    variant === 'panel' && pan != null ? `pan-${pan}` : '',
    color ? `button-${color}` : ''
  ].filter(Boolean).join(' ');

  function click(e) {
    onclick?.(e);
    if (!sounds) return;
    e.preventDefault();
    playBeepAndGo(2, href ?? '#');
  }
</script>

{#if href}
  <a class={cls} {href} onclick={click}>
    {#if children}{@render children()}{:else}{label}{/if}
  </a>
{:else}
  <button class={cls} onclick={click}>
    {#if children}{@render children()}{:else}{label}{/if}
  </button>
{/if}
