import { http } from '@/lib/http.ts';

// UpdateSource mirrors the server's shared application/firmware source config.
export type UpdateSource = {
  url: string;
  isOfficial: boolean;
  enabled: boolean;
};

// ManualUpdateState is the set of persisted states exposed to the UI.
export type ManualUpdateState =
  | 'staged'
  | 'installing'
  | 'succeeded'
  | 'failed'
  | 'reboot_scheduled';

// ManualUpdateInspection is returned after an upload passes server inspection.
export type ManualUpdateInspection = {
  id: string;
  version: string;
  size?: number;
  state: ManualUpdateState;
};

// get application version
export function getVersion() {
  return http.get('/api/application/version');
}

// update application to latest version
export function update() {
  return http.request({
    method: 'post',
    url: '/api/application/update',
    timeout: 15 * 60 * 1000 // 15 minutes
  });
}

// enable/disable preview updates
export function setPreviewUpdates(enable: boolean) {
  const data = {
    enable
  };
  return http.post('/api/application/preview', data);
}

// get preview updates state
export function getPreviewUpdates() {
  return http.get('/api/application/preview');
}

// get the single update-server root used for application and firmware updates
export function getUpdateSource() {
  return http.get('/api/application/update-source');
}

// save and validate a custom update-server root
export function setUpdateSource(url: string) {
  return http.post('/api/application/update-source', { url });
}

// switch future updates back to the official update-server root
export function resetUpdateSource() {
  return http.post('/api/application/update-source/reset');
}

// stage and inspect an uploaded package without installing it
export function inspectManualUpdate(data: FormData) {
  return http.request({
    method: 'post',
    url: '/api/application/update/manual/inspect',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data,
    timeout: 15 * 60 * 1000
  });
}

// install a previously inspected and staged package
export function confirmManualUpdate(id: string) {
  return http.request({
    method: 'post',
    url: '/api/application/update/manual/confirm',
    data: { id },
    timeout: 15 * 60 * 1000
  });
}

// get the state of a manual update after it has been confirmed
export function getManualUpdate(id: string) {
  return http.get(`/api/application/update/manual/${encodeURIComponent(id)}`);
}

// discard a staged package that will not be installed
export function cancelManualUpdate(id: string) {
  return http.delete(`/api/application/update/manual/${encodeURIComponent(id)}`);
}
