import { IMessageEvent, w3cwebsocket as W3cWebSocket } from 'websocket';

import { getBaseUrl } from '@/lib/service.ts';

type MessageHandler = (message: IMessageEvent) => void;
type ConnectionHandler = () => void;
type SendData = number[] | ArrayBuffer | Uint8Array;

export enum MessageEvent {
  Heartbeat = 0,
  Keyboard = 1,
  Mouse = 2
}

interface WsClientOptions {
  url?: string;
  heartbeatInterval?: number;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

const DEFAULT_OPTIONS: Required<WsClientOptions> = {
  url: `${getBaseUrl('ws')}/api/ws`,
  heartbeatInterval: 10 * 1000,
  reconnectInterval: 3 * 1000,
  maxReconnectAttempts: 1
};

export class WsClient {
  private readonly options: Required<WsClientOptions>;
  private instance: W3cWebSocket | null = null;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private reconnectAttempts = 0;
  private shouldReconnect = true;

  private readonly eventHandlers = new Map<string, Set<MessageHandler>>();
  // Connection hooks are separate from JSON message handlers and also run for explicit close().
  private readonly openHandlers = new Set<ConnectionHandler>();
  private readonly closeHandlers = new Set<ConnectionHandler>();

  constructor(options: WsClientOptions = {}) {
    this.options = { ...DEFAULT_OPTIONS, ...options };
  }

  public connect(): void {
    this.shouldReconnect = true;
    this.reconnectAttempts = 0;
    this.createConnection();
  }

  public close(): void {
    this.shouldReconnect = false;
    this.cleanup();

    const instance = this.instance;
    if (!instance) {
      return;
    }

    // Notify while the socket is still writable so input owners can send their neutral reports.
    this.notifyConnectionHandlers(this.closeHandlers, 'close');
    this.instance = null;
    this.closeConnection(instance);
  }

  public on(type: string, handler: MessageHandler): () => void {
    if (!this.eventHandlers.has(type)) {
      this.eventHandlers.set(type, new Set());
    }

    this.eventHandlers.get(type)!.add(handler);

    return () => {
      const handlers = this.eventHandlers.get(type);
      if (handlers) {
        handlers.delete(handler);
        if (handlers.size === 0) {
          this.eventHandlers.delete(type);
        }
      }
    };
  }

  public off(type: string, handler?: MessageHandler): void {
    if (handler) {
      const handlers = this.eventHandlers.get(type);
      if (handlers) {
        handlers.delete(handler);
        if (handlers.size === 0) {
          this.eventHandlers.delete(type);
        }
      }
    } else {
      this.eventHandlers.delete(type);
    }
  }

  public onOpen(handler: ConnectionHandler): () => void {
    this.openHandlers.add(handler);
    return () => this.openHandlers.delete(handler);
  }

  public onClose(handler: ConnectionHandler): () => void {
    this.closeHandlers.add(handler);
    return () => this.closeHandlers.delete(handler);
  }

  public send(data: SendData): boolean {
    if (!this.instance || !this.isConnected) {
      return false;
    }

    try {
      if (data instanceof ArrayBuffer || (data as unknown) instanceof Uint8Array) {
        this.instance.send(data);
      } else {
        this.instance.send(JSON.stringify(data));
      }
    } catch (error) {
      console.error('[WebSocket] Send error:', error);
      return false;
    }

    return true;
  }

  public get isConnected(): boolean {
    return this.instance?.readyState === W3cWebSocket.OPEN;
  }

  private createConnection(): void {
    this.cleanup();

    // Retire a replaced socket; identity checks below ignore any events it delivers later.
    const previousInstance = this.instance;
    if (previousInstance) {
      this.instance = null;
      this.closeConnection(previousInstance);
    }

    const instance = new W3cWebSocket(this.options.url);
    instance.binaryType = 'arraybuffer';
    this.instance = instance;

    instance.onopen = () => {
      // A CONNECTING socket may finish opening after close()/connect() has replaced it.
      if (this.instance !== instance) {
        this.closeConnection(instance);
        return;
      }
      this.handleOpen();
    };
    instance.onclose = () => {
      // Only the current socket may stop heartbeat, release input, or schedule reconnect.
      if (this.instance !== instance) {
        return;
      }
      this.instance = null;
      this.handleClose();
    };
    instance.onerror = (error) => {
      if (this.instance === instance) {
        this.handleError(error);
      }
    };
    instance.onmessage = (message) => {
      if (this.instance === instance) {
        this.handleMessage(message);
      }
    };
  }

  private handleOpen(): void {
    this.reconnectAttempts = 0;
    this.notifyConnectionHandlers(this.openHandlers, 'open');
    this.startHeartbeat();
  }

  private handleClose(): void {
    this.stopHeartbeat();
    this.notifyConnectionHandlers(this.closeHandlers, 'close');
    this.scheduleReconnect();
  }

  private handleError(error: Error): void {
    console.error('[WebSocket] Error:', error);
  }

  private handleMessage(message: IMessageEvent): void {
    try {
      const data = JSON.parse(message.data as string);
      const handlers = this.eventHandlers.get(data.type);

      if (handlers) {
        handlers.forEach((handler) => handler(message));
      }
    } catch (err) {
      console.log(err);
    }
  }

  private startHeartbeat(): void {
    this.stopHeartbeat();
    this.heartbeatTimer = setInterval(() => {
      this.send(new Uint8Array([MessageEvent.Heartbeat]));
    }, this.options.heartbeatInterval);
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  private scheduleReconnect(): void {
    if (!this.shouldReconnect) {
      return;
    }

    if (this.reconnectAttempts >= this.options.maxReconnectAttempts) {
      console.error('[WebSocket] Max reconnect attempts reached');
      return;
    }

    this.reconnectAttempts++;
    console.log(`[WebSocket] Reconnecting... (attempt ${this.reconnectAttempts})`);

    this.reconnectTimer = setTimeout(() => {
      this.createConnection();
    }, this.options.reconnectInterval);
  }

  private cleanup(): void {
    this.stopHeartbeat();

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }

  private closeConnection(instance: W3cWebSocket): void {
    // close() can race with CONNECTING; failure is harmless because stale callbacks are ignored.
    if (
      instance.readyState !== W3cWebSocket.CONNECTING &&
      instance.readyState !== W3cWebSocket.OPEN
    ) {
      return;
    }

    try {
      instance.close();
    } catch (error) {
      console.error('[WebSocket] Close error:', error);
    }
  }

  private notifyConnectionHandlers(
    handlers: Set<ConnectionHandler>,
    event: 'close' | 'open'
  ): void {
    // One lifecycle consumer must not prevent heartbeat/reconnect or other cleanup handlers.
    handlers.forEach((handler) => {
      try {
        handler();
      } catch (error) {
        console.error(`[WebSocket] ${event} handler error:`, error);
      }
    });
  }
}

export const client = new WsClient();
