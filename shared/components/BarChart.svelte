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
  .bar-chart {
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 0.4rem;
    padding: 0.5rem 0.6rem;
    width: 100%;
    height: 100%;
    box-sizing: border-box;
  }

  .bar-row {
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .bar-label {
    font-size: 0.62rem;
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
    height: 11px;
    background: rgba(255 255 255 / 0.1);
    border-radius: 6px;
    overflow: hidden;
  }

  .bar-fill {
    height: 100%;
    border-radius: 6px;
    transition: width 0.7s ease;
  }

  .bar-value {
    font-size: 0.68rem;
    font-weight: bold;
    color: var(--space-white, #f5f6fa);
    min-width: 5rem;
    font-variant-numeric: tabular-nums;
    white-space: nowrap;
  }
</style>
