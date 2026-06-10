<script>
  /**
   * Decorative framed display with animated "audio level" bars
   * (.lcars-frame from the ultra layout — works anywhere wide enough,
   * 280px tall, ~320px wide).
   *
   * Props:
   *   horizontal — rotate the bar display 90° (.display-horizontal)
   *   lines      — number of animated bars (default 16, the original; the
   *                theme's nth-child animations cover 16, extras animate
   *                like the outermost bars)
   *   children   — replaces the animated bars with arbitrary content in
   *                the frame's center cell
   */

  /** @type {{ horizontal?: boolean, lines?: number, children?: import('svelte').Snippet }} */
  let { horizontal = false, lines = 16, children } = $props();
</script>

<div class="lcars-frame">
  <div class="frame-col-1">
    <div class="frame-col-1-cell-a"></div>
    <div class="frame-col-1-cell-b"></div>
    <div class="frame-col-1-cell-c"></div>
  </div>
  <div class="frame-col-2"></div>
  <div class="frame-col-3 {children ? '' : horizontal ? 'display-horizontal' : 'display-vertical'}">
    {#if children}
      {@render children()}
    {:else}
      {#each { length: lines } as _}
        <div class="line"></div>
      {/each}
    {/if}
  </div>
  <div class="frame-col-4"></div>
  <div class="frame-col-5">
    <div class="frame-col-5-cell-a"></div>
    <div class="frame-col-5-cell-b"></div>
    <div class="frame-col-5-cell-c"></div>
  </div>
</div>
