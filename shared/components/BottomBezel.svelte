<script>
  /**
   * LCARS "elbow" bezel mirrored to the bottom — a coloured frame piece with
   * a bottom bar and one or both side bars wrapping a black content area.
   * Pairs with TopBezel (and optional MiddleBezel) to close a full frame.
   *
   * Same real-div construction as TopBezel; the rounded corners are on the
   * bottom edge instead of the top.
   *
   * Props:
   *   sides        — 'left' | 'right' | 'both' (default 'left'): which side
   *                  bar(s) accompany the bottom bar
   *   color        — frame colour (default var(--bar-6-color))
   *   contentColor — inner area colour (default black)
   *   bar          — bottom bar thickness (default 15px)
   *   side         — side bar thickness (default 60px, independent of bar)
   *   innerRadius  — interior elbow radius (default 30px)
   *   outerRadius  — exterior corner radius; defaults to innerRadius + bar so
   *                  the coloured band keeps a uniform thickness through the
   *                  bend (the real LCARS-elbow relationship — additive, not a
   *                  fixed ratio, so it holds at any scale)
   *   children     — content rendered inside the black area
   */

  /** @type {{ sides?: 'left'|'right'|'both', color?: string, contentColor?: string, bar?: string, side?: string, outerRadius?: string, innerRadius?: string, children?: import('svelte').Snippet }} */
  let {
    sides = 'left',
    color = 'var(--bar-6-color, #c66)',
    contentColor = 'black',
    bar = '15px',
    side = '60px',
    innerRadius = '30px',
    outerRadius,
    children
  } = $props();

  // Outer radius defaults to innerRadius + bar, keeping the coloured band a
  // uniform thickness around the elbow — scale-correct at any bar size.
  const outer = $derived(outerRadius ?? `calc(${innerRadius} + ${bar})`);
  const hasLeft = $derived(sides === 'left' || sides === 'both');
  const hasRight = $derived(sides === 'right' || sides === 'both');
</script>

<div
  class="bottom-bezel"
  style:background-color={color}
  style:padding-bottom={bar}
  style:padding-left={hasLeft ? side : null}
  style:padding-right={hasRight ? side : null}
  style:border-bottom-left-radius={hasLeft ? outer : null}
  style:border-bottom-right-radius={hasRight ? outer : null}
>
  <div
    class="bezel-content"
    style:background-color={contentColor}
    style:border-bottom-left-radius={hasLeft ? innerRadius : null}
    style:border-bottom-right-radius={hasRight ? innerRadius : null}
  >
    {@render children?.()}
  </div>
</div>

<style>
  .bottom-bezel {
    box-sizing: border-box;
    width: 100%;
    margin-top: 5px;
  }
  .bezel-content {
    width: 100%;
    height: 100%;
  }
</style>
