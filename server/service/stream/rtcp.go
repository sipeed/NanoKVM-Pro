package stream

import (
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
)

// Drain audio RTCP too, but only video PLI/FIR may request a keyframe. A compound
// packet produces at most one request; VideoHub merges all peers within 100 ms.
func ReadRTCP(sender *webrtc.RTPSender, requestKeyFrame func()) {
	for {
		packets, _, err := sender.ReadRTCP()
		if err != nil {
			return
		}
		if requestKeyFrame != nil && needsKeyFrame(packets) {
			requestKeyFrame()
		}
	}
}

func needsKeyFrame(packets []rtcp.Packet) bool {
	for _, packet := range packets {
		switch packet.(type) {
		case *rtcp.PictureLossIndication, *rtcp.FullIntraRequest:
			return true
		}
	}
	return false
}
