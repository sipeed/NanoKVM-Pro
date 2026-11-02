// Receives Server H.264/H.265 RTP, saves Annex B and exercises video PLI/FIR.
package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	"os"
	"sync"
	"time"
)

type message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	codec := flag.String("codec", "h264", "h264 or h265")
	address := flag.String("url", "ws://127.0.0.1:18080/api/stream/h264/webrtc", "signaling URL")
	output := flag.String("output", "received.h264", "Annex B output")
	count := flag.Int("frames", 90, "access units to receive")
	flag.Parse()
	if (*codec != "h264" && *codec != "h265") || *count < 1 || *count > 100000 {
		panic("unsupported codec or frame count")
	}
	started := time.Now()
	engine := &webrtc.MediaEngine{}
	mime := webrtc.MimeTypeH265
	if *codec == "h264" {
		mime = webrtc.MimeTypeH264
	}
	must(engine.RegisterCodec(webrtc.RTPCodecParameters{RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: mime, ClockRate: 90000, RTCPFeedback: []webrtc.RTCPFeedback{{Type: "nack", Parameter: "pli"}, {Type: "ccm", Parameter: "fir"}}}, PayloadType: 96}, webrtc.RTPCodecTypeVideo))
	pc, err := webrtc.NewAPI(webrtc.WithMediaEngine(engine)).NewPeerConnection(webrtc.Configuration{})
	must(err)
	defer pc.Close()
	_, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly})
	must(err)
	ws, _, err := websocket.DefaultDialer.Dial(*address, nil)
	must(err)
	defer ws.Close()
	var wsMutex sync.Mutex
	write := func(event string, value any) {
		data, err := json.Marshal(value)
		must(err)
		wsMutex.Lock()
		defer wsMutex.Unlock()
		must(ws.WriteJSON(message{event, string(data)}))
	}
	done := make(chan map[string]any, 1)
	pc.OnTrack(func(track *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		file, err := os.Create(*output)
		must(err)
		defer file.Close()
		frames, packets, keys, gaps := 0, 0, 0, 0
		var lastSequence uint16
		var first time.Time
		var keyFrames []int
		prefix := []byte{0, 0, 0, 1}
		emit := func(nal []byte) { _, err := file.Write(prefix); must(err); _, err = file.Write(nal); must(err) }
		for frames < *count {
			packet, _, err := track.ReadRTP()
			if err != nil {
				return
			}
			if packets > 0 && packet.SequenceNumber != lastSequence+1 {
				gaps++
			}
			lastSequence = packet.SequenceNumber
			packets++
			p := packet.Payload
			if len(p) < 3 {
				continue
			}
			nalType := (p[0] >> 1) & 63
			if *codec == "h265" {
				switch nalType {
				case 48:
					for offset := 2; offset+2 <= len(p); {
						size := int(binary.BigEndian.Uint16(p[offset:]))
						offset += 2
						if size < 2 || offset+size > len(p) {
							panic("invalid AP")
						}
						emit(p[offset : offset+size])
						offset += size
					}
				case 49:
					if p[2]&128 != 0 {
						_, err := file.Write(prefix)
						must(err)
						_, err = file.Write([]byte{(p[0] & 129) | ((p[2] & 63) << 1), p[1]})
						must(err)
						if p[2]&63 >= 16 && p[2]&63 <= 21 {
							keys++
							keyFrames = append(keyFrames, frames)
						}
					}
					_, err = file.Write(p[3:])
					must(err)
				default:
					emit(p)
					if nalType >= 16 && nalType <= 21 {
						keys++
						keyFrames = append(keyFrames, frames)
					}
				}
			} else {
				nalType = p[0] & 31
				switch nalType {
				case 24:
					for offset := 1; offset+2 <= len(p); {
						size := int(binary.BigEndian.Uint16(p[offset:]))
						offset += 2
						if size < 1 || offset+size > len(p) {
							panic("invalid STAP")
						}
						emit(p[offset : offset+size])
						offset += size
					}
				case 28:
					if p[1]&128 != 0 {
						_, err := file.Write(prefix)
						must(err)
						_, err = file.Write([]byte{p[0]&224 | p[1]&31})
						must(err)
						if p[1]&31 == 5 {
							keys++
							keyFrames = append(keyFrames, frames)
						}
					}
					_, err = file.Write(p[2:])
					must(err)
				default:
					emit(p)
					if nalType == 5 {
						keys++
						keyFrames = append(keyFrames, frames)
					}
				}
			}
			if packet.Marker {
				frames++
				if first.IsZero() {
					first = time.Now()
				}
				if frames == 15 {
					for range 20 {
						must(pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}}))
					}
				}
				if frames == 30 {
					must(pc.WriteRTCP([]rtcp.Packet{&rtcp.FullIntraRequest{MediaSSRC: uint32(track.SSRC()), FIR: []rtcp.FIREntry{{SSRC: uint32(track.SSRC()), SequenceNumber: 1}}}}))
				}
			}
		}
		done <- map[string]any{"codec": *codec, "frames": frames, "packets": packets, "keys": keys, "key_frame_indices": keyFrames, "rtp_gaps": gaps, "first_frame_ms": first.Sub(started).Milliseconds(), "elapsed_ms": time.Since(started).Milliseconds(), "pli_burst": 20, "fir": 1}
	})
	candidates := []webrtc.ICECandidateInit{}
	go func() {
		for {
			var m message
			if ws.ReadJSON(&m) != nil {
				return
			}
			switch m.Event {
			case "video-answer":
				var answer webrtc.SessionDescription
				must(json.Unmarshal([]byte(m.Data), &answer))
				must(pc.SetRemoteDescription(answer))
				for _, candidate := range candidates {
					must(pc.AddICECandidate(candidate))
				}
				candidates = nil
			case "video-candidate":
				var candidate webrtc.ICECandidateInit
				must(json.Unmarshal([]byte(m.Data), &candidate))
				if pc.RemoteDescription() == nil {
					candidates = append(candidates, candidate)
				} else {
					must(pc.AddICECandidate(candidate))
				}
			}
		}
	}()
	offer, err := pc.CreateOffer(nil)
	must(err)
	gathered := webrtc.GatheringCompletePromise(pc)
	must(pc.SetLocalDescription(offer))
	<-gathered
	write("video-offer", pc.LocalDescription())
	select {
	case result := <-done:
		encoded, _ := json.Marshal(result)
		fmt.Println(string(encoded))
	case <-time.After(15 * time.Second):
		panic("WebRTC frame timeout")
	}
}
