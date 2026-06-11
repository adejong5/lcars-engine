import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';
import fs from 'fs';

// Every <dir>/<name>/index.html is a site — no manual registration needed.
// `sites/` holds the public demos; `private/` is the gitignored clone of the
// lcars-engine-deployments repo and may be absent (e.g. on CI).
function siteEntries(dir) {
  if (!fs.existsSync(dir)) return {};
  return Object.fromEntries(
    fs.readdirSync(dir)
      .filter(name => fs.existsSync(path.resolve(dir, name, 'index.html')))
      .map(name => [`${dir}-${name}`, path.resolve(dir, name, 'index.html')])
  );
}

export default defineConfig({
  plugins: [svelte()],

  // Relative asset URLs so the build works at any mount path
  // (e.g. GitHub Pages serves it under /<repo-name>/).
  base: './',

  resolve: {
    alias: {
      // Import shared assets as $shared/... instead of ../../shared/...
      '$shared': path.resolve('./shared')
    }
  },

  server: {
    host: true,
    allowedHosts: ['pandaserver'],
    open: '/sites/test/'
  },

  build: {
    rollupOptions: {
      input: {
        index: path.resolve('./index.html'),
        ...siteEntries('sites'),
        ...siteEntries('private')
      }
    }
  }
});
