<script>
  /**
   * The decorated LCARS image frame (.image-frame): title bar, bordered
   * image body, and the colored base blocks.
   *
   * Frame colors come from theme vars --primary-color, --secondary-color,
   * --accent-color, --title-color, --image-border-color; override via the
   * style prop for one-off colors.
   *
   * Props:
   *   title — text in the frame's title bar
   *   src   — image source (or pass arbitrary content as children)
   *   alt   — image alt text
   *   style — inline style passthrough for var overrides
   */

  /** @type {{ title?: string, src?: string, alt?: string, style?: string, children?: import('svelte').Snippet }} */
  let { title = '', src, alt = '', style, children } = $props();
</script>

<div class="image-frame" {style}>
  <div class="imgf-title">
    <div class="h4-wrapper">
      <h4>{title}</h4>
    </div>
  </div>
  <div class="imgf-image-body">
    <div class="image-holder">
      {#if children}
        {@render children()}
      {:else if src}
        <img {src} {alt} />
      {/if}
    </div>
  </div>
  <div class="imgf-base">
    <div class="imgf-block-1"></div>
    <div class="imgf-block-2"></div>
    <div class="imgf-block-3"></div>
    <div class="imgf-block-4"></div>
    <div class="imgf-block-5"></div>
  </div>
</div>
