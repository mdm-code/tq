importScripts('wasm_exec.js');
importScripts('sw.js');

const wasm = 'api.wasm';

addEventListener('install', (event) => {
  event.waitUntil(caches.open('api').then((cache) => cache.add(wasm)));
});

addEventListener('activate', (event) => {
  event.waitUntil(clients.claim());
});

registerWasmHTTPListener(wasm, { base: 'api' });
