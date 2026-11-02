package webrtc

import (
	"NanoKVM-Server/common"
	"NanoKVM-Server/service/stream"
	"github.com/pion/webrtc/v4/pkg/media"
)

func (t *Track) stopVideoStream() {
	t.videoStopOnce.Do(func() { close(t.videoDone) })
}

func (t *Track) startVideoStream() {
	t.videoOnce.Do(func() { go t.sendVideoStream() })
}

func (t *Track) sendVideoStream() {
	select {
	case <-t.videoDone:
		return
	default:
	}
	hub := stream.GetVideoHub()
	subscription := hub.Subscribe(common.STREAM_TYPE_H265_WEBRTC)
	stop := subscription.CloseOn(t.videoDone)
	defer stop()

	for {
		select {
		case <-t.videoDone:
			return
		case frame, ok := <-subscription.Frames:
			if !ok {
				return
			}

			if frame.Result < 0 || !hub.IsCurrent(frame) || !hub.IsMode(common.STREAM_TYPE_H265_WEBRTC) {
				continue
			}
			t.writeVideo(media.Sample{Data: frame.Data, Duration: frame.Duration})
		}
	}
}
