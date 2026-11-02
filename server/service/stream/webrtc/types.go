package webrtc

import (
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

type WebRTCManager struct {
	clients        map[*websocket.Conn]*Client
	clientSnapshot atomic.Pointer[[]*Client]
	audioSending   int32
	mutex          sync.Mutex
}

type Client struct {
	videoOnce     sync.Once
	videoStopOnce sync.Once
	videoDone     chan struct{}
	ws            *websocket.Conn
	video         *webrtc.PeerConnection
	audio         *webrtc.PeerConnection
	track         *Track
	mutex         sync.Mutex
}

type SignalingHandler struct {
	client *Client
}

type Track struct {
	playoutDelayExtensionID   atomic.Uint32
	playoutDelayExtensionData []byte
	videoPacketizer           rtp.Packetizer
	video                     *webrtc.TrackLocalStaticRTP
	audio                     *webrtc.TrackLocalStaticSample
}

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
