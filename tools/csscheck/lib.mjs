// Shared bits for the compat-CSS verification harness. Each check renders the
// running app twice — once as served (*.compat.css) and once with the original
// modern stylesheets route-swapped in — and asserts the compat rendering
// matches the true modern behaviour. Point CSSCHECK_URL at the app (default
// http://localhost:8088/) and run `npm run check` with the server up.
import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

export const BASE = process.env.CSSCHECK_URL || 'http://localhost:8088/';

const root = join(dirname(fileURLToPath(import.meta.url)), '..', '..');
export const srcCSS = (name) =>
  readFileSync(join(root, 'internal/render/static', name), 'utf8');

// originals maps each served compat filename to the pristine modern source, for
// route-swapping so a modern browser renders the upstream behaviour.
export const originals = () => ({
  'classic.compat.css': srcCSS('classic.css'),
  'components.compat.css': srcCSS('components.css'),
});

export async function routeOriginals(page, files = originals()) {
  for (const [n, css] of Object.entries(files))
    await page.route('**/' + n, (r) =>
      r.fulfill({ contentType: 'text/css', body: css }));
}

export function finish(fails, okMsg) {
  if (fails) {
    console.log(`\n${fails} mismatch(es)`);
    process.exitCode = 1;
  } else {
    console.log('\n' + okMsg);
  }
}
