package webrtc

import (
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4/pkg/media"
	log "github.com/sirupsen/logrus"
)

func (t *Track) updateExtension() {
	if t.playoutDelayExtensionData == nil || len(t.playoutDelayExtensionData) == 0 {
		playoutDelay := &rtp.PlayoutDelayExtension{
			MinDelay: 0,
			MaxDelay: 0,
		}
		playoutDelayExtensionData, err := playoutDelay.Marshal()
		if err == nil {
			t.playoutDelayExtensionData = playoutDelayExtensionData
		}
	}
}

func (t *Track) writeVideoSample(sample media.Sample) error {
	samples := uint32(sample.Duration.Seconds() * 90000)
	packets := t.videoPacketizer.Packetize(sample.Data, samples)

	for _, p := range packets {
		if id := t.playoutDelayExtensionID.Load(); id != 0 {
			if err := p.Header.SetExtension(uint8(id), t.playoutDelayExtensionData); err != nil {
				return err
			}
		}

		if err := t.video.WriteRTP(p); err != nil {
			log.Errorf("failed to write RTP: %v", err)
			return err
		}
	}

	return nil
}

func (t *Track) writeVideo(sample media.Sample) {
	err := t.writeVideoSample(sample)
	if err != nil {
		log.Errorf("failed to write h264 video: %s", err)
	}
}

func (t *Track) writeAudio(sample media.Sample) {
	err := t.audio.WriteSample(sample)
	if err != nil {
		log.Errorf("failed to send h264 audio: %s", err)
	}
}
