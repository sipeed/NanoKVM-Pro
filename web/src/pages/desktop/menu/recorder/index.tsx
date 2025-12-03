import { useEffect, useRef, useState } from 'react';
import { Tooltip } from 'antd';
import { useAtomValue } from 'jotai';
import { Video } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { videoModeAtom } from '@/jotai/screen.ts';

// Recording relies on captureStream() of the <video> element, which only exists
// in the WebRTC modes (mjpeg renders an <img>, direct modes an offscreen canvas).
const RECORDABLE_MODES = [null, 'h264-webrtc', 'h265-webrtc'];

const MIME_CANDIDATES = [
  { mimeType: 'video/mp4;codecs=avc1.640034', ext: '.mp4', accept: { 'video/mp4': ['.mp4'] } },
  { mimeType: 'video/webm;codecs=vp9', ext: '.webm', accept: { 'video/webm': ['.webm'] } },
  { mimeType: 'video/webm', ext: '.webm', accept: { 'video/webm': ['.webm'] } }
];

const pickMimeType = () => {
  return MIME_CANDIDATES.find((c) => MediaRecorder.isTypeSupported(c.mimeType));
};

// Browsers default to 2.5 Mbps regardless of resolution, which ruins desktop
// text. Scale with the actual stream instead (~0.12 bit/pixel/frame).
const getVideoBitrate = (track?: MediaStreamTrack) => {
  const { width = 1920, height = 1080, frameRate = 30 } = track?.getSettings() ?? {};
  return Math.min(60_000_000, Math.max(8_000_000, Math.round(width * height * frameRate * 0.12)));
};

export const Recorder = () => {
  const { t } = useTranslation();
  const videoMode = useAtomValue(videoModeAtom);
  const [isRecording, setIsRecording] = useState(false);
  const [elapsedMs, setElapsedMs] = useState(0);
  const mediaRecorderRef = useRef<MediaRecorder>();
  const fileWritableRef = useRef<FileSystemWritableFileStream | null>(null);
  const timerRef = useRef<number | null>(null);
  const startTimeRef = useRef<number>(0);

  const isSupported = 'showSaveFilePicker' in window && RECORDABLE_MODES.includes(videoMode);

  const stopTimer = () => {
    if (timerRef.current !== null) {
      window.clearInterval(timerRef.current);
      timerRef.current = null;
    }
  };

  const startTimer = () => {
    stopTimer();
    startTimeRef.current = Date.now();
    setElapsedMs(0);
    timerRef.current = window.setInterval(() => {
      setElapsedMs(Date.now() - startTimeRef.current);
    }, 1000);
  };

  useEffect(() => {
    return () => {
      stopTimer();

      // Finalize the file if unmounted (e.g. video mode switch) while recording
      const recorder = mediaRecorderRef.current;
      if (recorder && recorder.state !== 'inactive') {
        recorder.stop();
      }
    };
  }, []);

  const formatElapsed = (ms: number) => {
    const totalSeconds = Math.floor(ms / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
  };

  const handleStartRecording = async () => {
    try {
      const videoElement = document.getElementById('screen') as HTMLVideoElement | null;
      const stream: MediaStream | undefined = (videoElement as any)?.captureStream?.();
      if (!stream) {
        return;
      }

      const container = pickMimeType();
      if (!container) {
        return;
      }

      const handle = await (window as any).showSaveFilePicker({
        suggestedName: `recording-${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}${container.ext}`,
        types: [
          {
            description: 'Sipeed NanoKVM Recorder',
            accept: container.accept
          }
        ]
      });

      const writable = await handle.createWritable();
      fileWritableRef.current = writable;

      const recorder = new MediaRecorder(stream, {
        mimeType: container.mimeType,
        videoBitsPerSecond: getVideoBitrate(stream.getVideoTracks()[0])
      });

      recorder.ondataavailable = async (event) => {
        if (event.data && event.data.size > 0) {
          if (fileWritableRef.current) {
            await fileWritableRef.current.write(event.data);
          } else {
            recorder.stop();
          }
        }
      };

      recorder.onstop = async () => {
        if (fileWritableRef.current) {
          await fileWritableRef.current.close();
          fileWritableRef.current = null;
        }
        stopTimer();
        setElapsedMs(0);
        setIsRecording(false);
      };

      recorder.start(1000);
      mediaRecorderRef.current = recorder;
      setIsRecording(true);
      startTimer();
    } catch (err) {
      console.error(err);
    }
  };

  const handleStopRecording = () => {
    const recorder = mediaRecorderRef.current;
    if (recorder && recorder.state !== 'inactive') {
      recorder.stop();
      stopTimer();
      setElapsedMs(0);
      setIsRecording(false);
    }
  };

  return (
    <Tooltip
      title={isRecording ? t('recorder.toggleStop') : t('recorder.toggleStart')}
      placement="bottom"
      mouseEnterDelay={0.6}
    >
      <div
        className={`flex h-[28px] cursor-pointer items-center justify-center rounded p-1 text-white hover:bg-neutral-700/70 ${isSupported ? '' : 'pointer-events-none opacity-40'}`}
        onClick={isRecording ? handleStopRecording : handleStartRecording}
      >
        <Video className={isRecording ? 'animate-pulse text-red-400' : ''} size={18} />
        {isRecording && (
          <span className="p-1 text-xs text-red-300">{formatElapsed(elapsedMs)}</span>
        )}
      </div>
    </Tooltip>
  );
};
