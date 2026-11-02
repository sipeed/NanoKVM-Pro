package webrtc

import (
	"io"
	"os/exec"
	"sync"
	"sync/atomic"

	"NanoKVM-Server/service/stream/opus"
	"github.com/gorilla/websocket"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

type WebRTCManager struct {
	tracks        map[*websocket.Conn]*Track
	trackSnapshot atomic.Pointer[[]*Track]
	audioSending  int32
	mutex         sync.Mutex
}

type Client struct {
	ws    *websocket.Conn
	video *webrtc.PeerConnection
	audio *webrtc.PeerConnection
	mutex sync.Mutex
}

type SignalingHandler struct {
	client *Client
	track  *Track
}

type Track struct {
	videoOnce                 sync.Once
	videoStopOnce             sync.Once
	videoDone                 chan struct{}
	playoutDelayExtensionID   atomic.Uint32
	playoutDelayExtensionData []byte
	videoPacketizer           rtp.Packetizer
	video                     *webrtc.TrackLocalStaticRTP
	audio                     *webrtc.TrackLocalStaticSample
}

type AudioInputPlayer struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	mutex   sync.Mutex
	active  bool
	decoder *opus.Decoder
}

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
