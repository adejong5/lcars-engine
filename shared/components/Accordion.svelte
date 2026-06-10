<script>
  /**
   * LCARS accordion (collapsible section). Ports the toggle behavior from
   * the template's lcars.js: .active class on the header, animated
   * max-height on the content.
   *
   * Props:
   *   title      — header text
   *   open       — initial state (default closed)
   *   limitWidth — cap the width at 1240px (.limit-width)
   */

  /** @type {{ title: string, open?: boolean, limitWidth?: boolean, children?: import('svelte').Snippet }} */
  let { title, open = false, limitWidth = false, children } = $props();

  let isOpen = $state(open);
  let contentEl = $state(null);
</script>

<div class="accordion-wrapper {limitWidth ? 'limit-width' : ''}">
  <button class="accordion" class:active={isOpen} onclick={() => (isOpen = !isOpen)}>
    {title}
  </button>
  <div
    class="accordionContent"
    bind:this={contentEl}
    style:max-height={isOpen && contentEl ? contentEl.scrollHeight + 'px' : null}
  >
    {@render children?.()}
  </div>
</div>
