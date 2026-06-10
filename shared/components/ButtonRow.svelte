<script>
  /**
   * Layout wrapper for a group of LCARSButtons (.buttons + justify-* helper).
   *
   * Props:
   *   justify — 'space-between' | 'center' | 'flex-end' | 'space-around' |
   *             'space-evenly'; omit for default flex-start
   */

  /** @type {{ justify?: string, children?: import('svelte').Snippet }} */
  let { justify, children } = $props();
</script>

<div class="buttons {justify ? `justify-${justify}` : ''}">
  {@render children?.()}
</div>
