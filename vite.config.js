import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';

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
        // Add a new entry here for each site you create.
        index: path.resolve('./index.html'),
        test: path.resolve('./sites/test/index.html'),
        padd: path.resolve('./sites/padd/index.html')
      }
    }
  }
});
