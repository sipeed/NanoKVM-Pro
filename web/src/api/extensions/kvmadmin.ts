import { http } from '@/lib/http.ts';

// install kvmadmin
export function install() {
  return http.post('/api/extensions/kvmadmin/install');
}

// uninstall kvmadmin
export function uninstall() {
  return http.post('/api/extensions/kvmadmin/uninstall');
}

// get kvmadmin status
export function getStatus() {
  return http.get('/api/extensions/kvmadmin/status');
}

// start kvmadmin
export function start() {
  return http.post('/api/extensions/kvmadmin/start');
}

// stop kvmadmin
export function stop() {
  return http.post('/api/extensions/kvmadmin/stop');
}
