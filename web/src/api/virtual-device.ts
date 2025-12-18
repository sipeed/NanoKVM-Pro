import { http } from '@/lib/http.ts';

// get virtual devices status
export function getVirtualDevice() {
  return http.get('/api/vm/device/virtual');
}

// mount/unmount virtual device
export function updateVirtualDevice(device: string, type: string) {
  const data = {
    device,
    type
  };

  return http.post('/api/vm/device/virtual', data);
}

// refresh virtual device
export function refreshVirtualDevice(device: string) {
  const data = {
    device
  };
  return http.post('/api/vm/device/virtual/refresh', data);
}
