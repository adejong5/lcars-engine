<script>
  /**
   * Segmented vertical level meter — pill-shaped segments filling from the
   * bottom up, each segment showing a slice of the gradient at its position
   * (so the colour reflects the level, not just the fill amount).
   *
   * Props:
   *   value    — current level
   *   min/max  — range (default 0–100)
   *   segments — number of pill segments (default 16)
   *   gap      — gap between segments in SVG units (default 3)
   *   stops    — gradient colours bottom→top, as CSS colour strings
   *              (CSS custom properties work: 'var(--gold)')
   */

  /** @type {{ value?: number, min?: number, max?: number, segments?: number, gap?: number, stops?: string[] }} */
  let {
    value    = 0,
    min      = 0,
    max      = 100,
    segments = 16,
    gap      = 3,
    stops    = ['var(--butterscotch,#f96)', 'var(--african-violet,#c9f)', 'var(--ice,#9cf)']
  } = $props();

  // Stable per-instance gradient ID — not reactive, computed once.
  const gid = `lm${(Math.random() * 0xfffff | 0).toString(16)}`;

  // Fixed SVG coordinate space; preserveAspectRatio="none" scales it to
  // whatever CSS size the parent gives the element.
  const VW = 40;
  const VH = 280;

  const frac      = $derived(Math.min(1, Math.max(0, (value - min) / (max - min))));
  const litCount  = $derived(Math.round(frac * segments));

  // Geometry (in SVG units)
  const segH = $derived((VH - (segments - 1) * gap) / segments);
  const r    = $derived(segH / 2); // pill: radius = half height
</script>

<svg
  class="level-meter"
  viewBox="0 0 {VW} {VH}"
  preserveAspectRatio="none"
  role="meter"
  aria-valuemin={min}
  aria-valuemax={max}
  aria-valuenow={value}
>
  <defs>
    <!--
      Gradient runs in user-space from y=VH (bottom, first stop)
      to y=0 (top, last stop), so each rect inherits the gradient
      colour at its own y position — bottom segments get the first
      colour, top segments get the last.
    -->
    <linearGradient id={gid} x1="0" y1={VH} x2="0" y2="0" gradientUnits="userSpaceOnUse">
      {#each stops as color, i}
        <stop offset="{(i / (stops.length - 1)) * 100}%" style="stop-color: {color}" />
      {/each}
    </linearGradient>
  </defs>

  {#each Array.from({ length: segments }) as _, i}
    {@const fromBottom = segments - 1 - i}
    {@const y = i * (segH + gap)}
    <rect
      x="2"
      y={y}
      width={VW - 4}
      height={segH}
      rx={r}
      ry={r}
      fill={fromBottom < litCount ? `url(#${gid})` : 'rgba(255,255,255,0.06)'}
    />
  {/each}
</svg>

<style>
  .level-meter {
    display: block;
    width: 100%;
    height: 100%;
  }
</style>
