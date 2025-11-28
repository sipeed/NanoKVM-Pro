package webrtc

import (
	"encoding/json"
	"errors"

	"github.com/pion/webrtc/v4"
	log "github.com/sirupsen/logrus"
)

func NewSignalingHandler(client *Client) *SignalingHandler {
	return &SignalingHandler{
		client: client,
	}
}

// RegisterCallbacks Register callback functions
func (s *SignalingHandler) RegisterCallbacks() {
	// video ICE candidate
	s.client.video.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		candidateByte, err := json.Marshal(candidate.ToJSON())
		if err != nil {
			log.Errorf("failed to marshal video candidate: %s", err)
			return
		}

		if err := s.client.WriteMessage("video-candidate", string(candidateByte)); err != nil {
			log.Errorf("failed to send video candidate: %s", err)
		}
	})

	// audio ICE candidate
	s.client.audio.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		candidateByte, err := json.Marshal(candidate.ToJSON())
		if err != nil {
			log.Errorf("failed to marshal audio candidate: %s", err)
			return
		}

		if err := s.client.WriteMessage("audio-candidate", string(candidateByte)); err != nil {
			log.Errorf("failed to send audio candidate or: %s", err)
			return
		}
	})

	manager := getManager()

	// video connection state change
	s.client.video.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		if state == webrtc.ICEConnectionStateConnected {
			manager.StartVideoStream()
		}

		log.Debugf("video connection state changed to %s", state.String())
	})

	// audio connection state change
	s.client.audio.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		if state == webrtc.ICEConnectionStateConnected {
			manager.StartAudioStream()
		}

		log.Debugf("audio connection state changed to %s", state.String())
	})
}

// HandleMessage handle the received message
func (s *SignalingHandler) HandleMessage(message *Message) error {
	switch message.Event {
	case "video-offer":
		return s.handleVideoOffer(message.Data)
	case "video-candidate":
		return s.handleVideoCandidate(message.Data)
	case "audio-offer":
		return s.handleAudioOffer(message.Data)
	case "audio-candidate":
		return s.handleAudioCandidate(message.Data)
	case "heartbeat":
		return s.handleHeartbeat()
	default:
		log.Debugf("Unhandled message event: %s", message.Event)
		return nil
	}
}

func (s *SignalingHandler) handleVideoOffer(data string) error {
	if s.client.video.SignalingState() != webrtc.SignalingStateStable {
		err := errors.New("video signaling is not stable")
		log.Error(err)
		return err
	}

	offer := webrtc.SessionDescription{}
	if err := json.Unmarshal([]byte(data), &offer); err != nil {
		log.Errorf("failed to unmarshal video offer: %s", err)
		return err
	}

	if err := s.client.video.SetRemoteDescription(offer); err != nil {
		log.Errorf("failed to set remote description: %s", err)
		return err
	}

	answer, err := s.client.video.CreateAnswer(nil)
	if err != nil {
		log.Errorf("failed to create answer: %s", err)
		return err
	}

	if err := s.client.video.SetLocalDescription(answer); err != nil {
		log.Errorf("failed to set local description: %s", err)
		return err
	}

	if err := s.updateHeaderExtensionID(); err != nil {
		log.Errorf("could not update header extension ID: %v", err)
		return err
	}

	answerByte, err := json.Marshal(answer)
	if err != nil {
		log.Errorf("failed to marshal answer: %s", err)
		return err
	}

	return s.client.WriteMessage("video-answer", string(answerByte))
}

// set extension ID
func (s *SignalingHandler) updateHeaderExtensionID() error {
	receivers := s.client.video.GetReceivers()
	if len(receivers) == 0 {
		return errors.New("no RTP receiver found for video")
	}

	params := receivers[0].GetParameters()
	if len(params.HeaderExtensions) == 0 {
		return errors.New("no header extensions found in negotiated parameters")
	}

	for _, ext := range params.HeaderExtensions {
		if ext.URI == "http://www.webrtc.org/experiments/rtp-hdrext/playout-delay" {
			s.client.track.playoutDelayExtensionID = uint8(ext.ID)
			log.Debugf("found and set playout delay extension ID to: %d", ext.ID)
			return nil
		}
	}

	log.Warnf("no track extension found in negotiated parameters, use default value 5")
	return nil
}

// handle video offer
func (s *SignalingHandler) handleAudioOffer(data string) error {
	if s.client.audio.SignalingState() != webrtc.SignalingStateStable {
		err := errors.New("audio signaling is not stable")
		log.Error(err)
		return err
	}

	offer := webrtc.SessionDescription{}
	if err := json.Unmarshal([]byte(data), &offer); err != nil {
		log.Errorf("failed to unmarshal audio offer: %s", err)
		return err
	}

	if err := s.client.audio.SetRemoteDescription(offer); err != nil {
		log.Errorf("failed to set remote description: %s", err)
		return err
	}

	answer, err := s.client.audio.CreateAnswer(nil)
	if err != nil {
		log.Errorf("failed to create answer: %s", err)
		return err
	}

	if err := s.client.audio.SetLocalDescription(answer); err != nil {
		log.Errorf("failed to set local description: %s", err)
		return err
	}

	answerByte, err := json.Marshal(answer)
	if err != nil {
		log.Errorf("failed to marshal answer: %s", err)
		return err
	}

	return s.client.WriteMessage("audio-answer", string(answerByte))
}

// handle video candidate
func (s *SignalingHandler) handleVideoCandidate(data string) error {
	candidate := webrtc.ICECandidateInit{}
	if err := json.Unmarshal([]byte(data), &candidate); err != nil {
		log.Errorf("failed to unmarshal candidate: %s", err)
		return err
	}

	if err := s.client.video.AddICECandidate(candidate); err != nil {
		log.Errorf("failed to add ICECandidate: %s", err)
		return err
	}

	return nil
}

// handl audio candidate
func (s *SignalingHandler) handleAudioCandidate(data string) error {
	candidate := webrtc.ICECandidateInit{}
	if err := json.Unmarshal([]byte(data), &candidate); err != nil {
		log.Errorf("failed to unmarshal candidate: %s", err)
		return err
	}

	if err := s.client.audio.AddICECandidate(candidate); err != nil {
		log.Errorf("failed to add ICECandidate: %s", err)
		return err
	}

	return nil
}

// handle heartbeat
func (s *SignalingHandler) handleHeartbeat() error {
	return s.client.WriteMessage("heartbeat", "")
}
