import { mount } from 'svelte';
import App from './App.svelte';
import { initHA } from '$shared/ha.svelte.js';
import '$shared/assets/classic.css';

// Real connection if configured in .env.local (VITE_HA_HOST / VITE_HA_TOKEN),
// otherwise mock data for offline development.
initHA(
  import.meta.env.VITE_HA_HOST
    ? { host: import.meta.env.VITE_HA_HOST, token: import.meta.env.VITE_HA_TOKEN }
    : { mock: true }
);

mount(App, { target: document.getElementById('app') });
