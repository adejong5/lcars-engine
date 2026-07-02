// text-box check: the compat build drops text-box and hoists the theme's
// @-moz-document block (the author's own no-text-box rendering — Firefox never
// shipped text-box). The reference is upstream-minus-text-box in Firefox,
// where the moz block applies natively: compat in Chrome AND Firefox must
// match it. Bar headings are excluded — the compat build deliberately replaces
// the author's top:-.7vh fix with measured transform centering.
import { chromium, firefox } from 'playwright';
import { BASE, srcCSS, finish } from './lib.mjs';

const strip = (s) => s.replace(/^[ \t]*text-box:[^;]*;[ \t]*\r?\n?/gm, '');
const up = {
  'classic.compat.css': strip(srcCSS('classic.css')),
  'components.compat.css': strip(srcCSS('components.css')),
};

const sels = ['h4', 'p', '.gallery li div', '.dc-row-2', '.lcars-accordion summary'];
async function metrics(engine, swap) {
  const b = await engine.launch();
  const p = await b.newPage({ viewport: { width: 1400, height: 950 } });
  if (swap)
    for (const [n, css] of Object.entries(up))
      await p.route('**/' + n, (r) => r.fulfill({ contentType: 'text/css', body: css }));
  await p.goto(BASE, { waitUntil: 'load' });
  await p.evaluate(() => document.fonts.ready);
  const out = await p.evaluate((sels) => {
    const o = {};
    for (const s of sels) {
      const e = document.querySelector(s);
      if (!e) continue;
      const cs = getComputedStyle(e), r = e.getBoundingClientRect();
      o[s] = { h: +r.height.toFixed(1), mt: cs.marginTop, mb: cs.marginBottom };
    }
    return o;
  }, sels);
  await b.close();
  return out;
}
const ref = await metrics(firefox, true); // old-Firefox simulation: author intent
const cChrome = await metrics(chromium, false);
const cFx = await metrics(firefox, false);
let fails = 0;
for (const s of sels) {
  if (!ref[s]) continue;
  for (const m of ['h', 'mt', 'mb']) {
    const near = (x, y) => Math.abs(parseFloat(x) - parseFloat(y)) <= 1.0;
    if (!near(ref[s][m], cChrome[s]?.[m]) || !near(ref[s][m], cFx[s]?.[m])) {
      fails++;
      console.log(`${s.padEnd(24)} ${m.padEnd(3)} ref=${ref[s][m]} chrome=${cChrome[s]?.[m]} fx=${cFx[s]?.[m]} ***`);
    }
  }
}
finish(fails, "text-box: compat == the author's no-text-box design in both engines");
