/**
 * LCARS beep sound effects (from the thelcars.com template assets).
 *
 * Usage:
 *   import { playBeep, playBeepAndGo } from '$shared/sounds.js';
 *   playBeep(2);                       // fire-and-forget
 *   playBeepAndGo(2, '/sites/other/'); // beep, then navigate
 */

import beep1 from './assets/beep1.mp3';
import beep2 from './assets/beep2.mp3';
import beep3 from './assets/beep3.mp3';
import beep4 from './assets/beep4.mp3';

const urls = { 1: beep1, 2: beep2, 3: beep3, 4: beep4 };
const cache = {};

function get(n) {
  if (!cache[n]) cache[n] = new Audio(urls[n]);
  return cache[n];
}

/** Play beep n (1–4). */
export function playBeep(n = 2) {
  const audio = get(n);
  audio.currentTime = 0;
  audio.play().catch(() => {}); // autoplay restrictions: ignore
}

/** Play beep n, then navigate to url when the sound finishes. */
export function playBeepAndGo(n, url) {
  if (!url || url === '#') { playBeep(n); return; }
  const audio = get(n);
  audio.currentTime = 0;
  audio.onended = () => { window.location.href = url; };
  audio.play().catch(() => { window.location.href = url; });
}
