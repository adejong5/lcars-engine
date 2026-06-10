/**
 * Home Assistant WebSocket client
 *
 * Usage:
 *   import { HAClient } from '../../shared/ha-websocket.js';
 *
 *   const ha = new HAClient({ host: '192.168.1.100', token: 'your-long-lived-token' });
 *   await ha.connect();
 *
 *   const state = await ha.getState('light.living_room');
 *
 *   const unsub = await ha.onStateChange('sensor.temperature', (newState) => {
 *     document.getElementById('temp').textContent = newState.state;
 *   });
 *
 *   // Listen to all state changes:
 *   await ha.onStateChange('*', (newState, entityId) => { ... });
 *
 *   ha.callService('light', 'turn_on', { brightness: 200 }, { entity_id: 'light.living_room' });
 */

export class HAClient {
  #host;
  #token;
  #ssl;
  #port;
  #ws = null;
  #msgId = 1;
  #pending = new Map();      // id → { resolve, reject }
  #subs = new Map();         // subId → callback (raw event callback)
  #stateListeners = new Map(); // entityId | '*' → Set<callback>
  #stateSubId = null;
  #intentionalClose = false;
  #connectResolve = null;
  #connectReject = null;
  #reconnectAttempts = 0;
  #onConnectionChange = null;

  constructor({ host, port = 8123, token, ssl = false, onConnectionChange = null }) {
    this.#host  = host;
    this.#port  = port;
    this.#token = token;
    this.#ssl   = ssl;
    this.#onConnectionChange = onConnectionChange;
  }

  get url() {
    const proto = this.#ssl ? 'wss' : 'ws';
    return `${proto}://${this.#host}:${this.#port}/api/websocket`;
  }

  // ── Public API ──────────────────────────────────────────────

  /** Connect and authenticate. Returns a Promise that resolves on auth_ok. */
  connect() {
    this.#intentionalClose = false;
    return new Promise((resolve, reject) => {
      this.#connectResolve = resolve;
      this.#connectReject  = reject;
      this.#open();
    });
  }

  /** Get the current state object for a single entity. */
  async getState(entityId) {
    const states = await this.getStates();
    return states.find(s => s.entity_id === entityId) ?? null;
  }

  /** Get all entity states as an array. */
  getStates() {
    return this.#send({ type: 'get_states' });
  }

  /**
   * Subscribe to state changes for an entity (or '*' for all).
   * Returns a Promise<unsubscribeFn>.
   *
   * Callback signature:
   *   (newState, oldState) for a specific entity
   *   (newState, entityId) when listening to '*'
   */
  async onStateChange(entityId, callback) {
    if (!this.#stateListeners.has(entityId)) {
      this.#stateListeners.set(entityId, new Set());
    }
    this.#stateListeners.get(entityId).add(callback);

    if (!this.#stateSubId) {
      this.#stateSubId = await this.#subscribeEvent('state_changed', () => {});
    }

    return () => {
      this.#stateListeners.get(entityId)?.delete(callback);
    };
  }

  /**
   * Call a Home Assistant service.
   *
   * @param {string} domain       e.g. 'light'
   * @param {string} service      e.g. 'turn_on'
   * @param {object} serviceData  e.g. { brightness: 200 }
   * @param {object|null} target  e.g. { entity_id: 'light.living_room' }
   */
  callService(domain, service, serviceData = {}, target = null) {
    const msg = { type: 'call_service', domain, service, service_data: serviceData };
    if (target) msg.target = target;
    return this.#send(msg);
  }

  /** Close the connection without auto-reconnecting. */
  disconnect() {
    this.#intentionalClose = true;
    this.#ws?.close();
  }

  // ── Internal ────────────────────────────────────────────────

  #open() {
    this.#ws = new WebSocket(this.url);

    this.#ws.addEventListener('message', e => {
      try { this.#onMessage(JSON.parse(e.data)); }
      catch { /* malformed message */ }
    });

    this.#ws.addEventListener('close', () => this.#onClose());

    this.#ws.addEventListener('error', () => {
      if (this.#connectReject) {
        this.#connectReject(new Error(`WebSocket error connecting to ${this.url}`));
        this.#connectReject = null;
      }
    });
  }

  #onMessage(msg) {
    switch (msg.type) {
      case 'auth_required':
        this.#ws.send(JSON.stringify({ type: 'auth', access_token: this.#token }));
        break;

      case 'auth_ok':
        this.#reconnectAttempts = 0;
        this.#onConnectionChange?.(true);
        if (this.#connectResolve) {
          this.#connectResolve(this);
          this.#connectResolve = null;
        } else {
          // Re-authenticated after a reconnect — restore subscriptions.
          this.#restoreSubscriptions();
        }
        break;

      case 'auth_invalid':
        const authErr = new Error(`HA auth failed: ${msg.message}`);
        if (this.#connectReject) { this.#connectReject(authErr); this.#connectReject = null; }
        this.#intentionalClose = true; // don't reconnect on bad token
        break;

      case 'result': {
        const p = this.#pending.get(msg.id);
        if (!p) break;
        this.#pending.delete(msg.id);
        if (msg.success) p.resolve(msg.result);
        else p.reject(new Error(msg.error?.message ?? 'Unknown HA error'));
        break;
      }

      case 'event': {
        const cb = this.#subs.get(msg.id);
        if (cb) cb(msg.event);

        // Route state_changed events to per-entity listeners.
        if (msg.id === this.#stateSubId && msg.event?.event_type === 'state_changed') {
          const { entity_id, new_state, old_state } = msg.event.data;
          this.#stateListeners.get(entity_id)?.forEach(fn => fn(new_state, old_state));
          this.#stateListeners.get('*')?.forEach(fn => fn(new_state, entity_id));
        }
        break;
      }
    }
  }

  #onClose() {
    this.#ws = null;
    this.#onConnectionChange?.(false);
    if (this.#intentionalClose) return;

    // Exponential back-off: 3s, 6s, 12s, … capped at 60s
    const delay = Math.min(3000 * 2 ** this.#reconnectAttempts, 60000);
    this.#reconnectAttempts++;
    console.warn(`[HAClient] Disconnected. Reconnecting in ${delay / 1000}s…`);
    setTimeout(() => this.#open(), delay);
  }

  async #restoreSubscriptions() {
    // Re-subscribe to state_changed if there are active listeners.
    if (this.#stateListeners.size > 0) {
      this.#stateSubId = null;
      this.#stateSubId = await this.#subscribeEvent('state_changed', () => {});
    }
  }

  #send(msg) {
    return new Promise((resolve, reject) => {
      const id = this.#msgId++;
      msg.id = id;
      this.#pending.set(id, { resolve, reject });
      this.#ws.send(JSON.stringify(msg));
    });
  }

  #subscribeEvent(eventType, callback) {
    return new Promise((resolve, reject) => {
      const id = this.#msgId++;
      const msg = { id, type: 'subscribe_events', event_type: eventType };
      this.#pending.set(id, {
        resolve: () => {
          this.#subs.set(id, callback);
          resolve(id);
        },
        reject
      });
      this.#ws.send(JSON.stringify(msg));
    });
  }
}

// ── Stardate utility ──────────────────────────────────────────

/**
 * Calculate a TNG-era stardate.
 * Year 2364 = 41000; each Earth year ≈ 1000 stardate units.
 */
export function getStardate(date = new Date()) {
  const year      = date.getFullYear();
  const start     = new Date(year, 0, 1);
  const end       = new Date(year + 1, 0, 1);
  const fraction  = (date - start) / (end - start);
  const stardate  = (year - 2323) * 1000 + fraction * 1000;
  return stardate.toFixed(1);
}
