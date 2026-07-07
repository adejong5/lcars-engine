// panel-2 check: components.compat.css precomputes min-height per breakpoint in place
// of max(outer-radius, bar-height + inner-radius). Reconstruct the original
// max() rule, serve it in place of the table, and assert the computed
// min-height is identical at widths straddling every breakpoint.
import { chromium } from 'playwright';
import { BASE, srcCSS, finish } from './lib.mjs';

const table = /\.left-frame-top \.panel-2\s*\{ min-height: 160px; \}[\s\S]*?min-height: 44px; \} \}/;
const maxRule = `:root                      { --elbow-outer-radius: 160px; --elbow-inner-radius: 60px; }
@media (max-width: 1500px) { :root { --elbow-outer-radius: 130px; --elbow-inner-radius: 60px; } }
@media (max-width: 1300px) { :root { --elbow-outer-radius: 100px; --elbow-inner-radius: 40px; } }
@media (max-width: 750px)  { :root { --elbow-outer-radius: 80px;  --elbow-inner-radius: 34px; } }
@media (max-width: 525px)  { :root { --elbow-outer-radius: 40px;  --elbow-inner-radius: 34px; } }
.left-frame-top .panel-2 {
  min-height: max(var(--elbow-outer-radius), calc(var(--bar-height) + var(--elbow-inner-radius)));
}`;
const orig = srcCSS('components.compat.css').replace(table, maxRule);
if (!orig.includes('max(var(--elbow-outer-radius)')) {
  console.error('panel2: could not reconstruct the max() rule — components.compat.css table changed?');
  process.exit(1);
}

const b = await chromium.launch();
async function mh(swap, w) {
  const p = await b.newPage({ viewport: { width: w, height: 900 } });
  if (swap)
    await p.route('**/components.compat.css', (r) =>
      r.fulfill({ contentType: 'text/css', body: orig }));
  await p.goto(BASE, { waitUntil: 'load' });
  const v = await p.evaluate(() => {
    const el = document.querySelector('.left-frame-top .panel-2');
    return el ? getComputedStyle(el).minHeight : null;
  });
  await p.close();
  return v;
}
let fails = 0;
for (const w of [1920, 1501, 1500, 1301, 1300, 900, 751, 750, 526, 525, 400]) {
  const c = await mh(false, w), u = await mh(true, w);
  if (c !== u) {
    fails++;
    console.log(`${String(w).padEnd(5)} compat=${c}  max()-rule=${u} ***`);
  }
}
await b.close();
finish(fails, 'panel-2: precomputed table == live max() at all widths');
