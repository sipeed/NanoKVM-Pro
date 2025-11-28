package ws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"NanoKVM-Server/service/ui"
	"NanoKVM-Server/service/vm/jiggler"
)

const (
	Heartbeats = iota
	KeyboardEvent
	MouseEvent
)

func (c *WsClient) Start() {
	defer c.Clean()

	ui.Online.Store(true)
	go c.hid.Keyboard(c.keyboard)
	go c.hid.Mouse(c.mouse)

	_ = c.Read()
}

func (c *WsClient) Read() error {
	var zeroTime time.Time
	_ = c.conn.SetReadDeadline(zeroTime)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}

		log.Debugf("received message: %s", message)

		var event []int
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Debugf("received invalid message: %s", message)
			continue
		}

		if event[0] == Heartbeats {
			c.updateHeartbeat()
			continue
		}

		switch event[0] {
		case KeyboardEvent:
			c.keyboard <- event[1:]
		case MouseEvent:
			c.mouse <- event[1:]
		default:
			log.Debugf("received invalid event: %d", event[0])
		}

		jiggler.GetJiggler().Update()
	}
}

func (c *WsClient) Write(event string, data string) error {
	message := &Message{
		Type: event,
		Data: data,
	}

	messageByte, err := json.Marshal(message)
	if err != nil {
		log.Errorf("failed to marshal message: %s", err)
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.conn.WriteMessage(websocket.TextMessage, messageByte)
}

func (c *WsClient) Clean() {
	RemoveClient(c)

	_ = c.conn.Close()

	ui.Online.Store(false)
	go clearQueue(c.keyboard)
	close(c.keyboard)

	go clearQueue(c.mouse)
	close(c.mouse)

	log.Debug("websocket disconnected")
}

func clearQueue(queue chan []int) {
	for range queue {
	}
}

func (c *WsClient) updateHeartbeat() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastHeartbeat = time.Now()
}

func (c *WsClient) getLastHeartbeat() time.Time {
	return c.lastHeartbeat
}

func (c *WsClient) keepAlive() {
	const timeout = 60 * time.Second
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(c.getLastHeartbeat()) > timeout {
			ui.Online.Store(false)
			log.Debug("heartbeat timeout, set offline")
			return
		} else {
			ui.Online.Store(true)
		}
	}
}
