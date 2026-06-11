<script>
  /**
   * Two-column grid of pill buttons (.pillbox / .pillbox-2 from the ultra
   * layout). Colors are assigned by position via the theme's nth-child
   * rules (6 designed slots per variant).
   *
   * Props:
   *   pills   — a number (that many auto-labeled pills) or an array of
   *             { label?, href?, spacer? }; label auto-generates like
   *             'J-001', omit href for a decorative beep-only pill,
   *             spacer: true renders an empty non-button slot
   *   variant — 2 for the second color set (.pillbox-2)
   *   sounds  — false to disable the beep
   */

  import { playBeepAndGo } from '../sounds.js';

  /** @type {{ pills?: number | Array<{label?: string, href?: string, spacer?: boolean}>, variant?: 1|2, sounds?: boolean }} */
  let { pills = 6, variant = 1, sounds = true } = $props();

  const letters = 'JRICAFGHKL';

  const items = (typeof pills === 'number'
    ? Array.from({ length: pills }, () => ({}))
    : pills
  ).map((p, i) => p.spacer ? p : {
    label: p.label ?? `${letters[i % letters.length]}-${String(i + 1).padStart(3, '0')}`,
    href: p.href
  });

  const pillClass = variant === 2 ? 'pill-2' : 'pill';

  function click(p) {
    if (sounds) playBeepAndGo(1, p.href ?? '#');
    else if (p.href) window.location.href = p.href;
  }
</script>

<div class={variant === 2 ? 'pillbox-2' : 'pillbox'}>
  {#each items as p}
    {#if p.spacer}
      <div class={pillClass}>&nbsp;</div>
    {:else}
      <button class={pillClass} onclick={() => click(p)}>{p.label}</button>
    {/if}
  {/each}
</div>
