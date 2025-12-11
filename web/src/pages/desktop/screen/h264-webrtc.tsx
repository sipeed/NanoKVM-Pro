import { useEffect, useRef, useState } from 'react';
import { Spin } from 'antd';
import clsx from 'clsx';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';

import * as api from '@/api/stream.ts';
import { VideoStatus } from '@/types';
import { mouseStyleAtom } from '@/jotai/mouse.ts';
import { videoParametersAtom, videoStatusAtom, videoVolumeAtom } from '@/jotai/screen.ts';

export const H264Webrtc = () => {
  const videoParameters = useAtomValue(videoParametersAtom);
  const mouseStyle = useAtomValue(mouseStyleAtom);
  const setVideoStatus = useSetAtom(videoStatusAtom);
  const [volume, setVolume] = useAtom(videoVolumeAtom);

  const [isPlaying, setIsPlaying] = useState(true);
  const [isLoading, setIsLoading] = useState(true);

  const videoRef = useRef<HTMLVideoElement | null>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  const videoOfferSent = useRef(false);
  const audioOfferSent = useRef(false);

  const videoIceCandidates = useRef<RTCIceCandidate[]>([]);
  const audioIceCandidates = useRef<RTCIceCandidate[]>([]);

  useEffect(() => {
    const ws = api.webrtcH264();

    const iceServers = [{ urls: ['stun:stun.l.google.com:19302'] }];
    const video = new RTCPeerConnection({ iceServers });
    const audio = new RTCPeerConnection({ iceServers });

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
      audio.addTransceiver('audio', { direction: 'recvonly' });
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data as string);
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
            handleVideoStatus(data);
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

    const handleVideoAnswer = (data: any) => {
      if (video.signalingState !== 'have-local-offer') {
        videoOfferSent.current = false;
        console.warn(`Video signaling state incorrect for answer: ${video.signalingState}`);
        return;
      }

      video
        .setRemoteDescription(new RTCSessionDescription(data))
        .then(() => {
          videoOfferSent.current = false;
          videoIceCandidates.current.forEach((candidate) => {
            video
              .addIceCandidate(candidate)
              .catch((e) => console.error('Video candidate failed to add:', e.message));
          });
          videoIceCandidates.current = [];
        })
        .catch((error) => {
          console.error('Video answer set failed:', error);
          videoOfferSent.current = false;
        });
    };

    const handleVideoCandidate = (data: any) => {
      if (!data.candidate) {
        return;
      }

      const candidate = new RTCIceCandidate(data);
      if (video.remoteDescription) {
        video
          .addIceCandidate(candidate)
          .catch((e) => console.error('Video candidate failed to add:', e.message));
      } else {
        videoIceCandidates.current.push(candidate);
      }
    };

    const handleAudioAnswer = (data: any) => {
      if (audio.signalingState !== 'have-local-offer') {
        audioOfferSent.current = false;
        console.warn(`Audio signaling state incorrect for answer: ${audio.signalingState}`);
        return;
      }

      audio
        .setRemoteDescription(new RTCSessionDescription(data))
        .then(() => {
          audioOfferSent.current = false;
          audioIceCandidates.current.forEach((candidate) => {
            audio
              .addIceCandidate(candidate)
              .catch((e) => console.error('Audio candidate failed to add:', e.message));
          });
          audioIceCandidates.current = [];
        })
        .catch((error) => {
          console.error('Audio answer set failed:', error);
          audioOfferSent.current = false;
        });
    };

    const handleAudioCandidate = (data: any) => {
      if (!data.candidate) {
        return;
      }

      const candidate = new RTCIceCandidate(data);
      if (audio.remoteDescription) {
        audio
          .addIceCandidate(candidate)
          .catch((e) => console.error('Audio candidate failed to add:', e.message));
      } else {
        audioIceCandidates.current.push(candidate);
      }
    };

    const handleVideoStatus = (data: any) => {
      switch (Number(data)) {
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
          console.log('Unhandled event:', data);
          break;
      }
    };

    const sendMsg = (event: string, data: string) => {
      if (ws.readyState !== WebSocket.OPEN) {
        return;
      }

      try {
        ws.send(JSON.stringify({ event, data }));
      } catch (err) {
        console.error('Error sending event: ', err);
      }
    };

    const heartbeatTimer = setInterval(() => {
      sendMsg('heartbeat', '');
    }, 60 * 1000);

    setTimeout(() => {
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

      if (heartbeatTimer) {
        clearInterval(heartbeatTimer);
      }
    };
  }, []);

  useEffect(() => {
    if (!audioRef.current) return;

    if (volume > 0) {
      if (audioRef.current.paused) {
        audioRef.current.play().catch(console.error);
      }
      if (audioRef.current.muted) {
        audioRef.current.muted = false;
      }
    }

    audioRef.current.volume = volume / 100;
  }, [volume]);

  return (
    <Spin size="large" tip="Loading" spinning={isLoading}>
      <div className="flex h-screen w-screen items-start justify-center xl:items-center">
        <video
          id="screen"
          ref={videoRef}
          className={clsx(
            'block max-h-full min-h-[480px] min-w-[640px] max-w-full select-none object-scale-down',
            isPlaying ? 'opacity-100' : 'opacity-0',
            mouseStyle
          )}
          style={{ transform: `scale(${videoParameters.scale})` }}
          muted
          autoPlay
          playsInline
          controls={false}
          onClick={(e) => {
            e.stopPropagation();
            e.preventDefault();
          }}
          onPlaying={() => setIsLoading(false)}
          onError={(e) => console.error('Video element error:', e)}
        />

        <audio
          ref={audioRef}
          muted
          autoPlay
          playsInline
          onError={(e) => console.error('Audio element error:', e)}
        />
      </div>
    </Spin>
  );
};
