<script>
  /**
   * The animated cascade of numbers, placed in the top frame via LCARSPage's
   * topFrame snippet. Purely random by default; pass `values` to carry live
   * Home Assistant telemetry (requires initHA() in the site's main.js).
   *
   * Values fill column-major (down each column, then the next) so every
   * entity keeps a stable row position — and therefore a stable color/blink
   * group from the theme's dc-row-N classes. Real values come first in DOM
   * order; filler columns follow, so the theme's overflow clipping eats the
   * filler on narrow screens, not your data.
   *
   * Props:
   *   values   — list of entities to display; each item is an entity_id
   *              string or { entity, format } where format(stateString)
   *              returns the display text. Cells show random digits until
   *              the entity's state arrives, and '----' when HA reports it
   *              unavailable/unknown.
   *   columns  — total column count to pad to with random filler
   *              (default 24, the original template's density)
   *   filler   — exact number of random filler columns, overriding the
   *              pad-to-`columns` default; 0 = values only
   *   theme    — 'classic' | 'nemesis' | 'lower-decks'; sets the row
   *              pattern and wrapper class the theme CSS expects
   *   pattern  — override the theme's dc-row-N class index per row
   *   frozen   — id="frozen" disables the color-cycling animation
   *   wrapperClass — override the theme's wrapper class
   */

  import { ha } from '../ha.svelte.js';

  /** @type {{ values?: Array<string | {entity: string, format?: (state: string) => string}>, columns?: number, filler?: number, theme?: 'classic'|'nemesis'|'lower-decks', pattern?: number[], frozen?: boolean, wrapperClass?: string }} */
  let {
    values = null,
    columns = 24,
    filler,
    theme = 'classic',
    pattern,
    frozen = false,
    wrapperClass
  } = $props();

  pattern ??= theme === 'lower-decks' ? [1, 2, 3, 4] : [1, 1, 2, 3, 3, 4, 5, 6, 7];
  wrapperClass ??= theme === 'lower-decks' ? 'data-wrapper' : 'data-cascade-wrapper';

  const rows = pattern.length;

  function randomCell() {
    const len = 1 + Math.floor(Math.random() * 7);
    let s = '';
    for (let i = 0; i < len; i++) s += Math.floor(Math.random() * 10);
    return s;
  }

  // Column model: cell = { entity?, format?, seed }. Built once per
  // values/filler change; live updates flow through display() instead.
  const cols = $derived.by(() => {
    const cells = (values ?? []).map(v =>
      typeof v === 'string'
        ? { entity: v, seed: randomCell() }
        : { entity: v.entity, format: v.format, seed: randomCell() }
    );

    const valueCols = Math.ceil(cells.length / rows);
    while (cells.length < valueCols * rows) cells.push({ seed: randomCell() });

    const fillerCols = filler ?? Math.max(0, columns - valueCols);
    for (let i = 0; i < fillerCols * rows; i++) cells.push({ seed: randomCell() });

    const result = [];
    for (let c = 0; c < valueCols + fillerCols; c++) {
      result.push(cells.slice(c * rows, (c + 1) * rows));
    }
    return result;
  });

  function display(cell) {
    if (!cell.entity) return cell.seed;
    const raw = ha.state(cell.entity);
    if (raw == null) return cell.seed;
    if (raw === 'unavailable' || raw === 'unknown') return '----';
    return cell.format ? cell.format(raw) : raw;
  }
</script>

<div class={wrapperClass} id={frozen ? 'frozen' : 'default'}>
  {#each cols as column}
    <div class="data-column">
      {#each column as cell, i}
        <div class="dc-row-{pattern[i]}">{display(cell)}</div>
      {/each}
    </div>
  {/each}
</div>
