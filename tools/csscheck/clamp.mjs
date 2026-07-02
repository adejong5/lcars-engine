// clamp check: the compat build replaces clamp(MIN, PREF, MAX) with PREF plus
// @media blocks pinning MIN/MAX at the crossover viewports. For each original
// clamp expression, inject a hidden probe element carrying the real clamp()
// next to the styled element and assert their computed values agree — at
// viewports spanning every breakpoint and the theme's rem-regime boundary.
import { chromium } from 'playwright';
import { BASE, finish } from './lib.mjs';

const cases = [
  ['.banner', 'font-size', 'clamp(1.25rem, 0.75rem + 4vw, 4rem)'],
  ['h1', 'font-size', 'clamp(1.5rem, 1.25rem + 3.5vw, 4rem)'],
  ['h2', 'font-size', 'clamp(1.4rem, 1.1rem + 2.25vw, 2.3rem)'],
  ['.left-frame', 'font-size', 'clamp(.875rem, 2vw, 1rem)'],
  ['.lcars-text-bar', 'height', 'clamp(16px, 4vh, 41px)'],
  ['.lcars-bar', 'height', 'clamp(15px, 2vh, 25px)'],
  ['.header-content', 'padding-left', 'clamp(20px, 3vw, 50px)'],
  ['.first-bar-panel', 'margin-top', 'clamp(15px, 2vw, 30px)'],
  ['.postmeta', 'font-size', 'clamp(1.2rem, 0.88rem + 1.28vw, 1.6rem)'],
  ['.go-big', 'font-size', 'clamp(1.2rem, 1rem + 1vw, 1.275rem)'],
];
const V = [[240, 300], [360, 640], [600, 760], [768, 900], [840, 700], [960, 700],
  [1000, 750], [1290, 1250], [1300, 1251], [1310, 700], [1600, 900], [1729, 900],
  [1788, 1000], [1920, 1080], [2200, 1300], [2560, 1440], [500, 250]];

const b = await chromium.launch();
let fails = 0, checks = 0;
for (const [w, h] of V) {
  const p = await b.newPage({ viewport: { width: w, height: h } });
  await p.goto(BASE, { waitUntil: 'load' });
  await p.evaluate(() => document.fonts.ready);
  const res = await p.evaluate((cases) => {
    const out = [];
    for (const [sel, prop, expr] of cases) {
      const el = document.querySelector(sel);
      if (!el) continue;
      const probe = document.createElement('div');
      probe.style.setProperty(prop, expr);
      probe.style.position = 'absolute';
      probe.style.visibility = 'hidden';
      el.parentNode.appendChild(probe); // sibling: same vw/vh and rem context
      const cam = prop.replace(/-([a-z])/g, (_, c) => c.toUpperCase());
      const want = parseFloat(getComputedStyle(probe)[cam]);
      const got = parseFloat(getComputedStyle(el)[cam]);
      probe.remove();
      out.push({ sel, prop, got, want });
    }
    return out;
  }, cases);
  for (const r of res) {
    checks++;
    if (Math.abs(r.got - r.want) > 1.0) {
      fails++;
      console.log(`${(w + 'x' + h).padEnd(10)} ${r.sel.padEnd(16)} ${r.prop.padEnd(13)} got=${r.got.toFixed(1)} clamp=${r.want.toFixed(1)} ***`);
    }
  }
  await p.close();
}
await b.close();
finish(fails, `clamp: ALL ${checks} match true clamp() within 1px`);
