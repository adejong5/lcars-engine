<script>
  /**
   * LCARS "elbow" bezel — a coloured frame piece with a top bar and one or
   * both side bars wrapping a black content area. Stack with MiddleBezel and
   * BottomBezel to build a full LCARS frame around arbitrary content.
   *
   * The coloured frame is a real element (the outer div), not a CSS border:
   * the inner content div is inset by padding so the colour shows through as
   * the top bar and the side bar(s). The top corner on each side that carries
   * a side bar is rounded — outerRadius on the coloured exterior, innerRadius
   * on the content interior — which forms the elbow curve.
   *
   * Props:
   *   sides        — 'left' | 'right' | 'both' (default 'left'): which side
   *                  bar(s) accompany the top bar
   *   color        — frame colour (default var(--bar-6-color))
   *   contentColor — inner area colour (default black)
   *   bar          — top bar thickness (default 15px)
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
  class="top-bezel"
  style:background-color={color}
  style:padding-top={bar}
  style:padding-left={hasLeft ? side : null}
  style:padding-right={hasRight ? side : null}
  style:border-top-left-radius={hasLeft ? outer : null}
  style:border-top-right-radius={hasRight ? outer : null}
>
  <div
    class="bezel-content"
    style:background-color={contentColor}
    style:border-top-left-radius={hasLeft ? innerRadius : null}
    style:border-top-right-radius={hasRight ? innerRadius : null}
  >
    {@render children?.()}
  </div>
</div>

<style>
  .top-bezel {
    box-sizing: border-box;
    width: 100%;
    margin-bottom: 5px;
  }
  .bezel-content {
    width: 100%;
    height: 100%;
  }
</style>
