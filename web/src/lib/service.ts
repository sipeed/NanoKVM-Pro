export function getHostname(): string {
  return window.location.hostname;
}

export function getPort(): number {
  const port = window.location.port;
  if (port) {
    return parseInt(port, 10);
  }

  return window.location.protocol === 'https:' ? 443 : 80;
}

export function getBaseUrl(type: 'http' | 'ws'): string {
  let protocol = window.location.protocol;
  if (type === 'ws') {
    protocol = protocol === 'https:' ? 'wss:' : 'ws:';
  }

  const hostname = getHostname();
  const port = getPort();

  const isDefaultPort =
    ((protocol === 'https:' || protocol === 'wss:') && port === 443) ||
    ((protocol === 'http:' || protocol === 'ws:') && port === 80);

  return isDefaultPort ? `${protocol}//${hostname}` : `${protocol}//${hostname}:${port}`;
}
