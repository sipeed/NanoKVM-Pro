package webrtc

import (
	"NanoKVM-Server/common"
	"NanoKVM-Server/service/stream"
	"github.com/pion/webrtc/v4/pkg/media"
	"strconv"
)

func (c *Client) stopVideoStream() {
	c.videoStopOnce.Do(func() { close(c.videoDone) })
}

func (c *Client) startVideoStream() {
	c.videoOnce.Do(func() { go c.sendVideoStream() })
}

func (c *Client) sendVideoStream() {
	select {
	case <-c.videoDone:
		return
	default:
	}
	hub := stream.GetVideoHub()
	subscription := hub.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
	stop := subscription.CloseOn(c.videoDone)
	defer stop()
	status := 0
	for {
		select {
		case <-c.videoDone:
			return
		case nextStatus := <-subscription.Updates:
			if nextStatus != status {
				if err := c.WriteMessage("video-status", strconv.Itoa(nextStatus)); err != nil {
					return
				}
				status = nextStatus
			}

		case frame, ok := <-subscription.Frames:
			if !ok {
				return
			}

			if frame.Result < 0 || !hub.IsCurrent(frame) || !hub.IsMode(common.STREAM_TYPE_H264_WEBRTC) {
				continue
			}
			c.track.writeVideo(media.Sample{Data: frame.Data, Duration: frame.Duration})
		}
	}
}
