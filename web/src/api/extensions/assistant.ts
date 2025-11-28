import { http } from '@/lib/http.ts';

// install assistant dependencies
export function install() {
  return http.post('/api/extensions/assistant/install');
}

// start assistant dependencies
export function start() {
  return http.post('/api/extensions/assistant/start');
}
