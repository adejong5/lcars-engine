<script>
  /**
   * Horizontal bar chart — pass reactive items from the parent using $derived().
   *
   * Props:
   *   items — Array<{ label: string, value: number, max: number,
   *                   unit?: string, color?: string, decimals?: number }>
   */

  /** @type {{ items?: { label: string, value: number, max: number, unit?: string, color?: string, decimals?: number }[] }} */
  let { items = [] } = $props();

  function pct(value, max) {
    if (!Number.isFinite(value) || max <= 0) return 0;
    return Math.min(100, Math.max(0, (value / max) * 100));
  }

  function fmt(item) {
    if (!Number.isFinite(item.value)) return '—';
    return item.value.toFixed(item.decimals ?? 0) + (item.unit ?? '');
  }
</script>

<div class="bar-chart">
  {#each items as item}
    <div class="bar-row">
      <span class="bar-label">{item.label}</span>
      <div class="bar-track">
        <div
          class="bar-fill"
          style:width="{pct(item.value, item.max)}%"
          style:background={item.color ?? 'var(--butterscotch, #f96)'}
        ></div>
      </div>
      <span class="bar-value">{fmt(item)}</span>
    </div>
  {/each}
</div>

<style>
  /* font-size here is the em anchor — padding, gaps, and track height
     are all in em so they scale if the anchor changes */
  .bar-chart {
    display: flex;
    flex-direction: column;
    justify-content: center;
    font-size: 0.6rem;
    gap: 0.5em;
    padding: 0.9em 1em;
    width: 100%;
    height: 100%;
    box-sizing: border-box;
  }

  .bar-row {
    display: flex;
    align-items: center;
    gap: 0.65em;
  }

  .bar-label {
    font-size: 0.56rem; /* 10% smaller than original 0.62rem */
    line-height: 1;
    font-weight: bold;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--space-white, #f5f6fa);
    min-width: 5.5rem;
    text-align: right;
    white-space: nowrap;
    opacity: 0.85;
  }

  .bar-track {
    flex: 1;
    height: 1.2em; /* 0.6rem × 1.2 ≈ 11.5px; text is ~9px so ~1.25px above & below */
    background: rgba(255 255 255 / 0.1);
    border-radius: 0.5em;
    overflow: hidden;
  }

  .bar-fill {
    height: 100%;
    border-radius: 0.5em;
    transition: width 0.7s ease;
  }

  .bar-value {
    font-size: 0.61rem; /* 10% smaller than original 0.68rem */
    line-height: 1;
    font-weight: bold;
    color: var(--space-white, #f5f6fa);
    min-width: 5rem;
    font-variant-numeric: tabular-nums;
    white-space: nowrap;
  }
</style>
