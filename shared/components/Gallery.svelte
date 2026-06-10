<script>
  /**
   * LCARS image gallery (ul.gallery — flagged experimental in the theme CSS).
   *
   * Props:
   *   images — [{ src, alt?, caption?, href? }]; href wraps the image in a
   *            link (theme highlights the border on hover)
   *   thumbs — true to size images as 340px thumbnails (.thumbs)
   */

  /** @type {{ images?: Array<{src: string, alt?: string, caption?: string, href?: string}>, thumbs?: boolean }} */
  let { images = [], thumbs = false } = $props();
</script>

<ul class="gallery {thumbs ? 'thumbs' : ''}">
  {#each images as img}
    <li>
      {#if img.href}
        <a href={img.href}><img src={img.src} alt={img.alt ?? ''} /></a>
      {:else}
        <img src={img.src} alt={img.alt ?? ''} />
      {/if}
      {#if img.caption}
        <div>{img.caption}</div>
      {/if}
    </li>
  {/each}
</ul>
