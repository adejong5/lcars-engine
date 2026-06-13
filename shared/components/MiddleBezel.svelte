<script>
  /**
   * LCARS straight side rails — the middle section between TopBezel and
   * BottomBezel. Draws a coloured bar down one or both sides of a black
   * content area, with no top/bottom bar and no rounded corners, so it
   * stacks flush against the bezels above and below.
   *
   * Same real-div construction as TopBezel: the coloured outer div shows
   * through the side padding of the inner content div.
   *
   * Props:
   *   sides        — 'left' | 'right' | 'both' (default 'left')
   *   color        — frame colour (default var(--bar-6-color))
   *   contentColor — inner area colour (default black)
   *   side         — side bar thickness (default 60px)
   *   children     — content rendered inside the black area
   */

  /** @type {{ sides?: 'left'|'right'|'both', color?: string, contentColor?: string, side?: string, children?: import('svelte').Snippet }} */
  let {
    sides = 'left',
    color = 'var(--bar-6-color, #c66)',
    contentColor = 'black',
    side = '60px',
    children
  } = $props();

  const hasLeft = $derived(sides === 'left' || sides === 'both');
  const hasRight = $derived(sides === 'right' || sides === 'both');
</script>

<div
  class="middle-bezel"
  style:background-color={color}
  style:padding-left={hasLeft ? side : null}
  style:padding-right={hasRight ? side : null}
>
  <div class="bezel-content" style:background-color={contentColor}>
    {@render children?.()}
  </div>
</div>

<style>
  .middle-bezel {
    box-sizing: border-box;
    width: 100%;
  }
  .bezel-content {
    width: 100%;
    height: 100%;
  }
</style>
