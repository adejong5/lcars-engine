// flex-gap check: the compat build replaces flex gap with margin-based
// spacing. Measure the pixel distance between the first two children of each
// flex container, compat vs originals (real flex-gap), at several widths.
import { chromium } from 'playwright';
import { BASE, routeOriginals, finish } from './lib.mjs';

const sels = ['.bar-chart', '.demo-render', '.gallery', '.data-cascade-wrapper', '.pill-gauge'];

const b = await chromium.launch();
async function measure(swap, w) {
  const p = await b.newPage({ viewport: { width: w, height: 1000 } });
  if (swap) await routeOriginals(p);
  await p.goto(BASE, { waitUntil: 'load' });
  await p.evaluate(() => document.fonts.ready);
  const res = await p.evaluate((sels) => {
    const out = {};
    for (const s of sels) {
      const el = document.querySelector(s);
      if (!el) continue;
      const kids = [...el.children].filter((c) => {
        const cs = getComputedStyle(c);
        return cs.position !== 'absolute' && cs.display !== 'none';
      });
      if (kids.length < 2) { out[s] = null; continue; }
      const a = kids[0].getBoundingClientRect(), c = kids[1].getBoundingClientRect();
      const horiz = Math.abs(a.top - c.top) < a.height / 2;
      out[s] = horiz ? +(c.left - a.right).toFixed(2) : +(c.top - a.bottom).toFixed(2);
    }
    return out;
  }, sels);
  await p.close();
  return res;
}
let fails = 0;
for (const w of [1400, 900, 600]) {
  const compat = await measure(false, w), upstream = await measure(true, w);
  for (const s of sels) {
    if (compat[s] == null || upstream[s] == null) continue;
    if (Math.abs(compat[s] - upstream[s]) > 1.0) {
      fails++;
      console.log(`${w} ${s.padEnd(22)} compat=${compat[s]}  upstream=${upstream[s]} ***`);
    }
  }
}
await b.close();
finish(fails, 'flex-gap: spacing matches upstream flex-gap within 1px');
