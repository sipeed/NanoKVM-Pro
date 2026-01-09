import { MutableRefObject, useEffect, useRef, useState } from 'react';
import { Spin } from 'antd';
import clsx from 'clsx';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';

import * as api from '@/api/stream.ts';
import { VideoStatus } from '@/types';
import { microphoneEnabledAtom } from '@/jotai/audio.ts';
import { mouseStyleAtom } from '@/jotai/mouse.ts';
import { videoParametersAtom, videoStatusAtom, videoVolumeAtom } from '@/jotai/screen.ts';

interface WebRTCMessage {
  event: string;
  data: string;
}

const createAnswerHandler = (
  connection: RTCPeerConnection,
  offerSentRef: MutableRefObject<boolean>,
  candidatesRef: MutableRefObject<RTCIceCandidate[]>,
  type: 'video' | 'audio'
) => {
  return (data: any) => {
    if (connection.signalingState !== 'have-local-offer') {
      offerSentRef.current = false;
      console.warn(`${type} signaling state incorrect for answer: ${connection.signalingState}`);
      return;
    }

    connection
      .setRemoteDescription(new RTCSessionDescription(data))
      .then(() => {
        offerSentRef.current = false;
        candidatesRef.current.forEach((candidate) => {
          connection
            .addIceCandidate(candidate)
            .catch((e) => console.error(`${type} candidate failed to add:`, e.message));
        });
        candidatesRef.current = [];
      })
      .catch((error) => {
        console.error(`${type} answer set failed:`, error);
        offerSentRef.current = false;
      });
  };
};

const createCandidateHandler = (
  connection: RTCPeerConnection,
  candidatesRef: MutableRefObject<RTCIceCandidate[]>,
  type: 'video' | 'audio'
) => {
  return (data: any) => {
    if (!data.candidate) {
      return;
    }

    const candidate = new RTCIceCandidate(data);
    if (connection.remoteDescription) {
      connection
        .addIceCandidate(candidate)
        .catch((e) => console.error(`${type} candidate failed to add:`, e.message));
    } else {
      candidatesRef.current.push(candidate);
    }
  };
};

export const H264Webrtc = () => {
  const videoParameters = useAtomValue(videoParametersAtom);
  const mouseStyle = useAtomValue(mouseStyleAtom);
  const setVideoStatus = useSetAtom(videoStatusAtom);
  const [volume, setVolume] = useAtom(videoVolumeAtom);
  const [micEnabled, setMicEnabled] = useAtom(microphoneEnabledAtom);

  const [isPlaying, setIsPlaying] = useState(true);
  const [isLoading, setIsLoading] = useState(true);

  const videoRef = useRef<HTMLVideoElement | null>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  const videoOfferSent = useRef(false);
  const audioOfferSent = useRef(false);

  const videoIceCandidates = useRef<RTCIceCandidate[]>([]);
  const audioIceCandidates = useRef<RTCIceCandidate[]>([]);

  // References for microphone management
  const audioConnectionRef = useRef<RTCPeerConnection | null>(null);
  const micSenderRef = useRef<RTCRtpSender | null>(null);
  const micStreamRef = useRef<MediaStream | null>(null);
  const micTrackRef = useRef<MediaStreamTrack | null>(null);

  useEffect(() => {
    const ws = api.webrtcH264();

    const iceServers = [{ urls: ['stun:stun.l.google.com:19302'] }];
    const video = new RTCPeerConnection({ iceServers });
    const audio = new RTCPeerConnection({ iceServers });

    const handleVideoAnswer = createAnswerHandler(
      video,
      videoOfferSent,
      videoIceCandidates,
      'video'
    );
    const handleAudioAnswer = createAnswerHandler(
      audio,
      audioOfferSent,
      audioIceCandidates,
      'audio'
    );
    const handleVideoCandidate = createCandidateHandler(video, videoIceCandidates, 'video');
    const handleAudioCandidate = createCandidateHandler(audio, audioIceCandidates, 'audio');

    const handleVideoStatus = (data: number) => {
      switch (data) {
        case 1:
          setVideoStatus(VideoStatus.Normal);
          setIsPlaying(true);
          break;
        case -1:
          setVideoStatus(VideoStatus.NoImage);
          setIsPlaying(false);
          break;
        case -4:
          setVideoStatus(VideoStatus.InconsistentVideoMode);
          setIsPlaying(false);
          break;
        default:
          console.log('Unhandled video status:', data);
          break;
      }
    };

    const sendMsg = (event: string, data: string) => {
      if (ws.readyState !== WebSocket.OPEN) {
        return;
      }

      try {
        const message: WebRTCMessage = { event, data };
        ws.send(JSON.stringify(message));
      } catch (err) {
        console.error('Error sending event:', err);
      }
    };

    // --- Init Video ---
    video.onnegotiationneeded = async () => {
      if (videoOfferSent.current || video.signalingState !== 'stable') {
        console.log('Skipping video negotiation - Waiting for answer or state unstable');
        return;
      }

      try {
        videoOfferSent.current = true;
        const offer = await video.createOffer({
          offerToReceiveVideo: true,
          offerToReceiveAudio: false
        });

        await video.setLocalDescription(offer);
        sendMsg('video-offer', JSON.stringify(video.localDescription));
      } catch (error) {
        videoOfferSent.current = false;
        console.error('Video negotiation failed:', error);
      }
    };

    video.onconnectionstatechange = () => {
      if (video.iceConnectionState === 'connected') {
        setIsLoading(false);
      }
    };

    video.ontrack = (event) => {
      if (videoRef.current && event.track.kind === 'video') {
        videoRef.current.srcObject = new MediaStream([event.track]);
      }
    };

    // --- Init Audio ---
    audio.onnegotiationneeded = async () => {
      if (audioOfferSent.current || audio.signalingState !== 'stable') {
        console.log('Skipping audio negotiation - Waiting for answer or state unstable');
        return;
      }

      try {
        audioOfferSent.current = true;
        const offer = await audio.createOffer({
          offerToReceiveVideo: false,
          offerToReceiveAudio: true
        });

        await audio.setLocalDescription(offer);
        sendMsg('audio-offer', JSON.stringify(audio.localDescription));
      } catch (error) {
        audioOfferSent.current = false;
        console.error('Audio negotiation failed:', error);
      }
    };

    audio.ontrack = (event) => {
      if (audioRef.current && event.track.kind === 'audio') {
        audioRef.current.srcObject = new MediaStream([event.track]);
      }
    };

    // --- WebSocket Message Handling ---
    ws.onopen = () => {
      videoOfferSent.current = false;
      audioOfferSent.current = false;

      video.onicecandidate = (event) => {
        if (event.candidate) {
          sendMsg('video-candidate', JSON.stringify(event.candidate));
        }
      };

      audio.onicecandidate = (event) => {
        if (event.candidate) {
          sendMsg('audio-candidate', JSON.stringify(event.candidate));
        }
      };

      video.addTransceiver('video', { direction: 'recvonly' });
      audio.addTransceiver('audio', { direction: 'sendrecv' });

      audioConnectionRef.current = audio;
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data as string) as WebRTCMessage;
        if (!msg?.data) return;

        const data = JSON.parse(msg.data);
        if (!data) return;

        switch (msg.event) {
          case 'video-answer':
            handleVideoAnswer(data);
            break;
          case 'video-candidate':
            handleVideoCandidate(data);
            break;
          case 'audio-answer':
            handleAudioAnswer(data);
            break;
          case 'audio-candidate':
            handleAudioCandidate(data);
            break;
          case 'video-status':
            handleVideoStatus(Number(data));
            break;
          case 'heartbeat':
            break;
          default:
            console.log('Unhandled event:', msg.event);
        }
      } catch (err) {
        console.error('Message processing error:', err);
      }
    };

    const heartbeatTimer = setInterval(() => {
      sendMsg('heartbeat', '');
    }, 60 * 1000);

    const loadingTimer = setTimeout(() => {
      setIsLoading(false);
    }, 10 * 1000);

    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }

      video.close();
      audio.close();

      videoOfferSent.current = false;
      audioOfferSent.current = false;

      setVolume(0);

      clearInterval(heartbeatTimer);
      clearTimeout(loadingTimer);

      if (micStreamRef.current) {
        micStreamRef.current.getTracks().forEach((track) => track.stop());
        micStreamRef.current = null;
      }
      micSenderRef.current = null;
      micTrackRef.current = null;
    };
  }, []);

  useEffect(() => {
    const audioElement = audioRef.current;
    if (!audioElement) return;

    if (volume > 0) {
      if (audioElement.paused) {
        audioElement.play().catch(console.error);
      }
      if (audioElement.muted) {
        audioElement.muted = false;
      }
    }

    audioElement.volume = volume / 100;
  }, [volume]);

  // Handle microphone enable/disable
  useEffect(() => {
    const audio = audioConnectionRef.current;
    if (!audio) return;

    const enableMicrophone = async () => {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({
          audio: {
            sampleRate: 48000,
            channelCount: 2,
            echoCancellation: true,
            noiseSuppression: true,
            autoGainControl: true
          }
        });
        micStreamRef.current = stream;

        const track = stream.getAudioTracks()[0];
        micTrackRef.current = track;

        micSenderRef.current = audio.addTrack(track, stream);
      } catch (err) {
        console.error('Failed to get microphone:', err);
        setMicEnabled(false);
      }
    };

    const disableMicrophone = () => {
      if (micSenderRef.current) {
        try {
          audio.removeTrack(micSenderRef.current);
        } catch (e) {
          console.error('Failed to remove track:', e);
        }
        micSenderRef.current = null;
      }

      micTrackRef.current = null;
      if (micStreamRef.current) {
        micStreamRef.current.getTracks().forEach((track) => track.stop());
        micStreamRef.current = null;
      }
    };

    if (micEnabled) {
      enableMicrophone();
    } else {
      disableMicrophone();
    }
  }, [micEnabled, setMicEnabled]);

  return (
    <Spin size="large" tip="Loading" spinning={isLoading}>
      <div className="flex h-screen w-screen items-start justify-center xl:items-center">
        <video
          id="screen"
          ref={videoRef}
          className={clsx(
            'block max-h-full min-h-[50vh] min-w-[50vw] max-w-full select-none object-scale-down',
            isPlaying ? 'opacity-100' : 'opacity-0',
            mouseStyle
          )}
          style={{ transform: `scale(${videoParameters.scale})` }}
          muted
          autoPlay
          playsInline
          controls={false}
          onPlaying={() => setIsLoading(false)}
          onClick={(e) => {
            e.stopPropagation();
            e.preventDefault();
          }}
        />

        <audio ref={audioRef} muted autoPlay playsInline />
      </div>
    </Spin>
  );
};
