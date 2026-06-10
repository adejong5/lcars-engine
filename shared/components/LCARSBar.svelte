<script>
  /**
   * Horizontal LCARS divider bar for the content area.
   *
   * Without a title: the plain .lcars-bar (rounded end caps + center bar).
   * With a title: the .lcars-text-bar carrying a heading.
   *
   * Colors come from the theme vars --lcars-bar-color,
   * --lcars-bar-start-color, --lcars-bar-end-color, --lcars-bar-text-color;
   * override them via the style prop for one-off colors.
   *
   * Props:
   *   title   — heading text; omit for a plain bar
   *   heading — heading tag for the title: 'h2' (default) | 'h3' | 'h4' | 'span'
   *   end     — true to right-align the title (.the-end)
   *   style   — inline style passthrough (e.g. '--lcars-bar-color: var(--orange)')
   */

  /** @type {{ title?: string, heading?: 'h2'|'h3'|'h4'|'span', end?: boolean, style?: string }} */
  let { title, heading = 'h2', end = false, style } = $props();
</script>

{#if title}
  <div class="lcars-text-bar {end ? 'the-end' : ''}" {style}>
    <svelte:element this={heading}>{title}</svelte:element>
  </div>
{:else}
  <div class="lcars-bar" {style}></div>
{/if}
