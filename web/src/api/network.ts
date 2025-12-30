import { http } from '@/lib/http.ts';

// wake on lan
export function wol(mac: string) {
  const data = {
    mac
  };
  return http.post('/api/network/wol', data);
}

// get wake-on-lan macs history
export function getWolMacs() {
  return http.get('/api/network/wol/mac');
}

export function deleteWolMac(mac: string) {
  const data = {
    mac
  };
  return http.delete('/api/network/wol/mac', data);
}

// set Mac name
export function setWolMacName(mac: string, name: string) {
  return http.post('/api/network/wol/mac/name', { mac, name });
}

// get static IP
export function getStaticIp() {
  return http.get('/api/network/static-ip');
}

// set static IP
export function setStaticIP(enabled: boolean, ip: string) {
  const data = {
    enabled: enabled,
    ip: ip
  };
  return http.post('/api/network/static-ip', data);
}

// scan wifi
export function scanWiFi() {
  return http.get('/api/network/wifi/scan');
}

// connect wifi (only available in wifi configuration mode)
export function connectWifiNoAuth(ssid: string, password: string) {
  const data = {
    ssid,
    password
  };
  return http.post('/api/network/wifi', data);
}

// connect wifi
export function connectWifi(ssid: string, password: string) {
  const data = {
    ssid,
    password
  };
  return http.post('/api/network/wifi/connect', data);
}

// disconnect wifi
export function disconnectWifi() {
  return http.post('/api/network/wifi/disconnect');
}

// get connected wifi information
export function getWiFi() {
  return http.get('/api/network/wifi');
}
