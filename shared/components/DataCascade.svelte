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
   *              string or { entity, format, cls } where format(stateString)
   *              returns the display text and cls(stateString) returns
   *              extra classes for the cell (e.g. 'font-red blink' above a
   *              threshold). { text } makes a static cell instead (e.g. a
   *              column label; '' for an empty cell). Entity cells show
   *              random digits until the state arrives, and '----' when HA
   *              reports it unavailable/unknown. Values fill column-major,
   *              so with pattern [1,2,3,4] every 4 items form one column.
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

  /** @type {{ values?: Array<string | {entity?: string, format?: (state: string) => string, cls?: (state: string) => string, text?: string}>, columns?: number, filler?: number, theme?: 'classic'|'nemesis'|'lower-decks'|'padd', pattern?: number[], frozen?: boolean, wrapperClass?: string }} */
  let {
    values = null,
    columns = 24,
    filler,
    theme = 'classic',
    pattern,
    frozen = false,
    wrapperClass
  } = $props();

  const isPadd = theme === 'padd';
  const isLowerDecks = theme === 'lower-decks' || isPadd;

  pattern ??= isLowerDecks ? [1, 2, 3, 4] : [1, 1, 2, 3, 3, 4, 5, 6, 7];
  wrapperClass ??= isLowerDecks ? 'data-wrapper' : 'data-cascade-wrapper';

  const rows = pattern.length;

  // padd cells sprinkle these font colors across the cascade
  const paddFonts = ['font-arctic-ice', 'font-night-rain', 'font-alpha-blue'];

  function digits(len) {
    let s = '';
    for (let i = 0; i < len; i++) s += Math.floor(Math.random() * 10);
    return s;
  }

  function randomCell() {
    if (isPadd) {
      // padd mixes short cells, 7-digit, and 11-digit columns, with
      // random font colors on roughly half the cells
      const lens = [2, 2, 2, 3, 3, 4, 7, 11];
      return {
        seed: digits(lens[Math.floor(Math.random() * lens.length)]),
        font: Math.random() < 0.5
          ? paddFonts[Math.floor(Math.random() * paddFonts.length)]
          : undefined
      };
    }
    // lower-decks cells are short with extra hidden digits (.hide-data,
    // dropped below 1100px); classic/nemesis cells are 1-7 plain digits
    return isLowerDecks
      ? { seed: digits(2 + Math.floor(Math.random() * 2)), extra: digits(8) }
      : { seed: digits(1 + Math.floor(Math.random() * 7)) };
  }

  // lower-decks/padd interleave dark filler columns ('0'/'000' dark-on-dark
  // at the top, dim digits below)
  function darkColumn() {
    return pattern.map((_, i) => {
      const edge = isPadd ? i === 0 : i === 0 || i === pattern.length - 1;
      return {
        seed: edge ? (isPadd ? '000' : '0') : digits(2),
        dark: true,
        darkFont: edge,
        font: !edge && isPadd
          ? paddFonts[Math.floor(Math.random() * paddFonts.length)]
          : undefined
      };
    });
  }

  // Column model: cell = { entity?, format?, seed, extra?, dark?, darkFont? }.
  // Built once per values/filler change; live updates flow through display().
  const cols = $derived.by(() => {
    const cells = (values ?? []).map(v => {
      if (typeof v === 'string') return { entity: v, ...randomCell() };
      if (v.text != null) return { text: v.text };
      return { entity: v.entity, format: v.format, cls: v.cls, ...randomCell() };
    });

    const valueCols = Math.ceil(cells.length / rows);
    while (cells.length < valueCols * rows) cells.push(randomCell());

    const result = [];
    for (let c = 0; c < valueCols; c++) {
      result.push(cells.slice(c * rows, (c + 1) * rows));
    }

    const fillerCols = filler ?? Math.max(0, columns - valueCols);
    for (let c = 0; c < fillerCols; c++) {
      // lower-decks alternates dim "darkspace" columns into the filler
      result.push(
        isLowerDecks && c % 2 === 1
          ? darkColumn()
          : pattern.map(() => randomCell())
      );
    }
    return result;
  });

  function display(cell) {
    if (cell.text != null) return cell.text;
    if (!cell.entity) return cell.seed;
    const raw = ha.state(cell.entity);
    if (raw == null) return cell.seed;
    if (raw === 'unavailable' || raw === 'unknown') return '----';
    return cell.format ? cell.format(raw) : raw;
  }

  function extraClass(cell) {
    if (!cell.cls || !cell.entity) return '';
    const raw = ha.state(cell.entity);
    return raw == null ? '' : (cell.cls(raw) ?? '');
  }
</script>

<div class={wrapperClass} id={frozen ? 'frozen' : 'default'}>
  {#each cols as column}
    <div class="data-column">
      {#each column as cell, i}
        <div
          class="dc-row-{pattern[i]} {cell.font ?? ''} {extraClass(cell)}"
          class:darkspace={cell.dark}
          class:darkfont={cell.darkFont}
        >{display(cell)}{#if cell.extra}<span class="hide-data">{cell.extra}</span>{/if}</div>
      {/each}
    </div>
  {/each}
</div>
