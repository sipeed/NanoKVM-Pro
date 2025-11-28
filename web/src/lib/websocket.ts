import { IMessageEvent, w3cwebsocket as W3cWebSocket } from 'websocket';

import { getBaseUrl } from '@/lib/service.ts';

type Event = (message: IMessageEvent) => void;

const eventMap: Map<string, Event> = new Map<string, Event>();
const wildcardHandlers = new Set<Event>();

export class WsClient {
  private readonly url: string;
  private instance: W3cWebSocket;

  public connected: (() => void) | null = null;
  public disconnected: (() => void) | null = null;

  constructor(url = `${getBaseUrl('ws')}/api/ws`) {
    this.url = url;
    this.instance = new W3cWebSocket(this.url);
    this.setEvents();
  }

  public connect() {
    this.close();

    if (this.connected) {
      this.connected();
    }

    this.instance = new W3cWebSocket(this.url);
    this.setEvents();
  }

  public send(data: number[]) {
    if (this.instance.readyState !== W3cWebSocket.OPEN) {
      return;
    }

    const message = JSON.stringify(data);
    this.instance.send(message);
  }

  public close() {
    if (this.instance.readyState === W3cWebSocket.OPEN) {
      if (this.disconnected) {
        this.disconnected();
      }

      this.instance.close();
    }
  }

  public register(type: string, fn: (message: IMessageEvent) => void) {
    if (type === '*') {
      wildcardHandlers.add(fn);
    } else {
      eventMap.set(type, fn);
    }

    this.setEvents();
  }

  public unregister(type: string) {
    eventMap.delete(type);

    this.setEvents();
  }

  private setEvents() {
    this.instance.onmessage = (message) => {
      const data = JSON.parse(message.data as string);

      const fn = eventMap.get(data.type);
      if (fn) {
        fn(message);
      }

      wildcardHandlers.forEach((fn) => fn(message));
    };
  }
}

export const client = new WsClient();
