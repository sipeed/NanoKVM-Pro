package direct

import (
	"encoding/binary"
	"time"

	"NanoKVM-Server/common"
	"NanoKVM-Server/service/stream"
	"github.com/gorilla/websocket"
)

func writeStream(ws *websocket.Conn, done <-chan struct{}) {
	hub := stream.GetVideoHub()
	hub.SetMode(common.STREAM_TYPE_H265_DIRECT)
	subscription := hub.Subscribe(common.STREAM_TYPE_H265_DIRECT)
	defer subscription.Close()
	for {
		select {
		case <-done:
			return
		case frame, ok := <-subscription.Frames:
			if !ok {
				return
			}
			if frame.Result < 0 || !hub.IsCurrent(frame) || !hub.IsMode(common.STREAM_TYPE_H265_DIRECT) {
				continue
			}
			// Keep the existing Direct wire format: key flag, LE microsecond
			// timestamp, then the complete Annex B access unit.
			message := make([]byte, 9+len(frame.Data))
			if frame.KeyFrame {
				message[0] = 1
			}
			binary.LittleEndian.PutUint64(message[1:9], uint64(frame.Timestamp))
			copy(message[9:], frame.Data)
			_ = ws.SetWriteDeadline(time.Now().Add(time.Second))
			if err := ws.WriteMessage(websocket.BinaryMessage, message); err != nil {
				_ = ws.Close() // also unblock the request's read loop
				return
			}
		}
	}
}
