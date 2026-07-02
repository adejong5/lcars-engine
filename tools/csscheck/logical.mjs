// logical-properties check: the compat build rewrites margin/padding-block/
// -inline, border-block, and inset to physical longhands. Compare computed
// physical values, compat vs originals, on elements whose margins the
// flex-gap fallbacks don't touch (padding is safe everywhere).
import { chromium } from 'playwright';
import { BASE, routeOriginals, finish } from './lib.mjs';

const cases = [
  ['h2', ['marginTop', 'marginBottom']],
  ['h3', ['marginTop', 'marginBottom']],
  ['.lcars-bar', ['marginTop', 'marginBottom']],
  ['.lcars-text-bar', ['marginTop', 'marginBottom']],
  ['.lcars-accordion', ['marginTop', 'marginBottom']],
  ['.demo', ['marginTop', 'marginBottom']], // two-value margin-block
  ['.pill-gauge', ['paddingLeft', 'paddingRight']], // padding-inline
  ['.pill-gauge .fill', ['top', 'right', 'bottom', 'left']], // inset: 0
  ['.accordion', ['paddingTop', 'paddingBottom', 'paddingLeft', 'paddingRight']],
];

const b = await chromium.launch();
async function read(swap, w) {
  const p = await b.newPage({ viewport: { width: w, height: 900 } });
  if (swap) await routeOriginals(p);
  await p.goto(BASE, { waitUntil: 'load' });
  await p.evaluate(() => document.fonts.ready);
  const o = await p.evaluate((cases) => {
    const out = {};
    for (const [s, props] of cases) {
      const el = document.querySelector(s);
      if (!el) continue;
      const cs = getComputedStyle(el);
      out[s] = Object.fromEntries(props.map((pr) => [pr, cs[pr]]));
    }
    return out;
  }, cases);
  await p.close();
  return o;
}
let fails = 0, checks = 0;
for (const w of [1600, 1000, 600]) {
  const c = await read(false, w), u = await read(true, w);
  for (const [s, props] of cases) {
    if (!c[s] || !u[s]) continue;
    for (const pr of props) {
      checks++;
      if (c[s][pr] !== u[s][pr]) {
        fails++;
        console.log(`${w} ${s.padEnd(18)} ${pr.padEnd(14)} compat=${c[s][pr]} upstream=${u[s][pr]} ***`);
      }
    }
  }
}
await b.close();
finish(fails, `logical: ALL ${checks} physical values match the originals`);
