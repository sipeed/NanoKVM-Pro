import { http } from '@/lib/http.ts';

// get NanoKVM information
export function getInfo() {
  return http.get('/api/vm/info');
}

// get hardware information
export function getHardware() {
  return http.get('/api/vm/hardware');
}

// set gpio value
export function setGpio(type: string, duration: number) {
  const data = {
    type,
    duration
  };
  return http.post('/api/vm/gpio', data);
}

// get gpio value
export function getGpio() {
  return http.get('/api/vm/gpio');
}

// get OLED configuration
export function getOLED() {
  return http.get('/api/vm/oled');
}

// set OLED configuration
export function setOLED(sleep: number) {
  return http.post('/api/vm/oled', { sleep });
}

// get LCD time format
export function getLCDTimeFormat() {
  return http.get('/api/vm/lcd/time/format');
}

// set LCD time format
export function setLCDTimeFormat(format: string) {
  return http.post('/api/vm/lcd/time/format', {
    format
  });
}

// get HDMI capture status
export function getHdmiCapture() {
  return http.get('/api/vm/hdmi/capture');
}

// set HDMI capture status
export function setHdmiCapture(enabled: boolean) {
  return http.post('/api/vm/hdmi/capture', { enabled });
}

// get HDMI passthrough status
export function getHdmiPassthrough() {
  return http.get('/api/vm/hdmi/passthrough');
}

// set HDMI passthrough status
export function setHdmiPassthrough(enabled: boolean) {
  return http.post('/api/vm/hdmi/passthrough', { enabled });
}

// get SSH state
export function getSSHState() {
  return http.get('/api/vm/ssh');
}

// enable SSH
export function enableSSH() {
  return http.post('/api/vm/ssh/enable');
}

// disable SSH
export function disableSSH() {
  return http.post('/api/vm/ssh/disable');
}

// get mouse jiggler
export function getMouseJiggler() {
  return http.get('/api/vm/mouse-jiggler');
}

// set mouse jiggler
export function setMouseJiggler(enabled: boolean, mode: string) {
  return http.post('/api/vm/mouse-jiggler', { enabled, mode });
}

// get Hostname
export function getHostname() {
  return http.get('/api/vm/hostname');
}

// set Hostname
export function setHostname(hostname: string) {
  return http.post('/api/vm/hostname', { hostname });
}

// get WebTitle
export function getWebTitle() {
  return http.get('/api/vm/web-title');
}

// set WebTitle
export function setWebTitle(title: string) {
  return http.post('/api/vm/web-title', { title });
}

// get mDNS state
export function getMdnsState() {
  return http.get('/api/vm/mdns');
}

// enable mDNS
export function enableMdns() {
  return http.post('/api/vm/mdns/enable');
}

// disable mDNS
export function disableMdns() {
  return http.post('/api/vm/mdns/disable');
}

// get low power mode
export function getLowPowerMode() {
  return http.get('/api/vm/low-power');
}

// set low power mode
export function setLowPowerMode(enable: boolean) {
  return http.post('/api/vm/low-power', { enable });
}

// switch system to PiKVM
export function switchPiKVM() {
  return http.post('/api/vm/system/pikvm');
}

// reboot
export function reboot() {
  return http.post('/api/vm/system/reboot');
}

export function getLedConfig() {
  return http.get('/api/vm/ledstrip/get');
}

export function setLedConfig(config: {
  brightness: number;
  hor?: number;
  on?: boolean;
  ver?: number;
}) {
  return http.post('/api/vm/ledstrip/set', config);
}

export function updateEdid(edid: string) {
  return http.post('/api/vm/edid', {
    edid
  });
}

export function uploadEdid(formData: FormData) {
  return http.request({
    method: 'post',
    url: '/api/vm/edid/upload',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
}

export function getCustomEdidList() {
  return http.get('/api/vm/edid/custom');
}

export function deleteEdid(edid: string) {
  return http.post('/api/vm/edid/delete', { edid });
}

export function getEdid() {
  return http.get('/api/vm/edid');
}

export function getTimezone() {
  return http.get('/api/vm/timezone');
}

export function setTimezone(timezone: string) {
  return http.post('/api/vm/timezone', {
    timezone
  });
}

export function getTimeStatus() {
  return http.get('/api/vm/time/status');
}

export function syncTime() {
  return http.post('/api/vm/time/sync');
}

// get menu bar config
export function getMenuBarConfig() {
  return http.get('/api/vm/menubar');
}

// set menu bar config
export function setMenuBarConfig(disabledItems: string[]) {
  return http.post('/api/vm/menubar', { disabledItems });
}
