// Shared bits for the compat-CSS verification harness. Both served
// stylesheets are maintained compat sources (no modern originals to compare
// against — see decisions/0008): the remaining checks are self-contained.
// clamp.mjs probes served rules against real clamp(); panel2.mjs reconstructs
// its max() rule from the served table. Point CSSCHECK_URL at the app
// (default http://localhost:8088/) and run `npm run check` with the server up.
import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

export const BASE = process.env.CSSCHECK_URL || 'http://localhost:8088/';

const root = join(dirname(fileURLToPath(import.meta.url)), '..', '..');
export const srcCSS = (name) =>
  readFileSync(join(root, 'internal/render/static', name), 'utf8');

export function finish(fails, okMsg) {
  if (fails) {
    console.log(`\n${fails} mismatch(es)`);
    process.exitCode = 1;
  } else {
    console.log('\n' + okMsg);
  }
}
