import { client, MessageEvent } from '@/lib/websocket.ts';

/**
 * Coordinates mouse state across mode components and WebSocket reconnects.
 * Reports queued here outlive an Absolute/Relative component unmount.
 */
type ReleaseHandler = () => void;
export type MouseInterface = 'absolute' | 'relative';

let releaseHandler: ReleaseHandler | null = null;
// Each HID interface needs its own neutral report because report length selects hidg1/hidg2.
const pendingReleases = new Map<MouseInterface, Uint8Array>();

function buildMouseMessage(report: Uint8Array): Uint8Array {
  return new Uint8Array([MessageEvent.Mouse, ...report]);
}

function flushPendingMouseReleases(): boolean {
  // A stale release must reach the host before any input from a new mode/connection.
  for (const [mouseInterface, data] of pendingReleases) {
    if (client.send(data)) {
      pendingReleases.delete(mouseInterface);
    }
  }

  return pendingReleases.size === 0;
}

export function initializeMouseInputLifecycle(): () => void {
  // Unexpected close clears local button state; open retries the matching HID releases first.
  const removeOpenHandler = client.onOpen(flushPendingMouseReleases);
  const removeCloseHandler = client.onClose(releaseMouseInput);

  return () => {
    removeOpenHandler();
    removeCloseHandler();
  };
}

export function sendMouseReport(mouseInterface: MouseInterface, report: Uint8Array): boolean {
  // Do not let a new press overtake a neutral report left behind by a broken connection.
  if (!flushPendingMouseReleases()) {
    return false;
  }

  const sent = client.send(buildMouseMessage(report));
  if (sent && report[0] === 0) {
    pendingReleases.delete(mouseInterface);
  }

  return sent;
}

export function sendMouseRelease(mouseInterface: MouseInterface, report: Uint8Array): void {
  const data = buildMouseMessage(report);
  // Record first, then remove only after send accepted the complete neutral report.
  pendingReleases.set(mouseInterface, data);

  if (client.send(data)) {
    pendingReleases.delete(mouseInterface);
  }
}

export function registerMouseReleaseHandler(handler: ReleaseHandler): () => void {
  // Only the mounted mouse mode owns the current local button/touch state.
  releaseHandler = handler;

  return () => {
    if (releaseHandler === handler) {
      releaseHandler = null;
    }
  };
}

export function releaseMouseInput(): void {
  releaseHandler?.();
}
