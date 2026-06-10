<script>
  /**
   * The main navigation button group from the top frame, placed via
   * LCARSPage's topFrame snippet (the theme CSS styles `nav` inside
   * .data-cascade-button-group).
   *
   * Props:
   *   buttons — a number (that many auto-labeled buttons, default 4 to
   *             match the original template) or an array of
   *             { label?, href? }:
   *               label — auto-numbered ('01'…) when omitted
   *               href  — navigation target; omit for a decorative
   *                       button that just beeps
   *   sounds  — false to disable the beep on click
   */

  import { playBeepAndGo } from '../sounds.js';

  /** @type {{ buttons?: number | Array<{label?: string, href?: string}>, sounds?: boolean }} */
  let { buttons = 4, sounds = true } = $props();

  const items = (typeof buttons === 'number'
    ? Array.from({ length: buttons }, () => ({}))
    : buttons
  ).map((b, i) => ({
    label: b.label ?? String(i + 1).padStart(2, '0'),
    href: b.href
  }));

  function nav(href) {
    if (sounds) playBeepAndGo(2, href ?? '#');
    else if (href) window.location.href = href;
  }
</script>

<nav>
  {#each items as b}
    <button onclick={() => nav(b.href)}>{b.label}</button>
  {/each}
</nav>
