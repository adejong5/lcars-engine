/**
 * Reactive Home Assistant store (Svelte 5 runes).
 *
 * Site setup (main.js):
 *   import { initHA } from '$shared/ha.svelte.js';
 *   initHA({ host: '192.168.1.100', token: '...' });   // real
 *   initHA({ mock: true });                            // fake data for offline dev
 *
 * Component usage:
 *   import { ha } from '$shared/ha.svelte.js';
 *   ha.state('sensor.temp')   // reactive state string (undefined until known)
 *   ha.entity('sensor.temp')  // reactive full state object (attributes etc.)
 *   ha.connected              // reactive connection status
 *
 * Reactivity is fine-grained: updates are written per-key into a SvelteMap,
 * so a state change touches only the template expressions reading that entity.
 *
 * Mock mode: any entity a component reads is lazily seeded with a plausible
 * value (on/off for binary_sensor.*, a drifting number otherwise) and a
 * random subset of entities updates every `mockInterval` ms, mimicking HA's
 * push behavior.
 */

import { SvelteMap } from 'svelte/reactivity';
import { HAClient } from './ha-websocket.js';

const entities = new SvelteMap(); // entity_id → HA state object
let connected = $state(false);

let client = null;
let mockMode = false;

export const ha = {
  get connected() {
    return connected;
  },

  /** Latest state string for an entity, or undefined if unknown. */
  state(id) {
    return this.entity(id)?.state;
  },

  /** Full state object ({ state, attributes, ... }), or undefined. */
  entity(id) {
    // Lazily fabricate any entity a component asks for in mock mode.
    // Deferred so we never write state during a template read.
    if (mockMode && !entities.has(id)) queueMicrotask(() => mockSeed(id));
    return entities.get(id);
  },

  callService(domain, service, data = {}, target = null) {
    if (mockMode) return Promise.resolve();
    if (!client) return Promise.reject(new Error('initHA() has not been called'));
    return client.callService(domain, service, data, target);
  }
};

/**
 * Connect the store. Call once per site, before or after mount.
 * @param {{ host?: string, port?: number, token?: string, ssl?: boolean,
 *           mock?: boolean, mockInterval?: number }} config
 */
export async function initHA(config = {}) {
  if (config.mock) {
    mockMode = true;
    connected = true;
    setInterval(mockTick, config.mockInterval ?? 2000);
    return;
  }

  client = new HAClient({
    host: config.host,
    port: config.port,
    token: config.token,
    ssl: config.ssl,
    onConnectionChange: up => { connected = up; }
  });

  await client.connect();
  connected = true;

  for (const s of await client.getStates()) {
    entities.set(s.entity_id, s);
  }
  await client.onStateChange('*', (newState, entityId) => {
    if (newState) entities.set(entityId, newState);
  });
}

// ── Mock data ─────────────────────────────────────────────────

function mockSeed(id) {
  if (!entities.has(id)) {
    entities.set(id, { entity_id: id, state: mockValue(id), attributes: {} });
  }
}

function mockValue(id) {
  const prev = entities.get(id);

  if (id.startsWith('binary_sensor.')) {
    if (!prev) return Math.random() > 0.5 ? 'on' : 'off';
    // flip occasionally
    return Math.random() > 0.8 ? (prev.state === 'on' ? 'off' : 'on') : prev.state;
  }

  if (!prev) return String(Math.floor(2 + Math.random() * 9998));
  // drift ±5% so values look like live telemetry
  const base = parseFloat(prev.state);
  return String(Math.round(base * (1 + (Math.random() * 0.1 - 0.05)) * 10) / 10);
}

function mockTick() {
  for (const [id, e] of entities) {
    if (Math.random() < 0.35) {
      entities.set(id, { ...e, state: mockValue(id) });
    }
  }
}
