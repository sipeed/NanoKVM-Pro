import { w3cwebsocket } from 'websocket';

import { http } from '@/lib/http.ts';
import { getBaseUrl } from '@/lib/service.ts';

// set video stream mode
export function setMode(mode: string) {
  const data = { mode };
  return http.post('/api/stream/mode', data);
}

// set stream quality
export function setQuality(quality: number) {
  const data = { quality };
  return http.post('/api/stream/quality', data);
}

// set stream GOP
export function setGop(gop: number) {
  const data = { gop };
  return http.post('/api/stream/gop', data);
}

// WebRTC H.264 connection
export function webrtcH264() {
  return getWebSocket('/api/stream/h264/webrtc');
}

// direct H.264 connection
export function directH264() {
  return getWebSocket('/api/stream/h264/direct');
}

function getWebSocket(url: string) {
  const baseUrl = getBaseUrl('ws');

  return new w3cwebsocket(`${baseUrl}${url}`);
}
