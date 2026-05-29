import { useAtomValue } from 'jotai';

import { videoModeAtom } from '@/jotai/screen.ts';

import { H264Direct } from './h264-direct.tsx';
import { H264Webrtc } from './h264-webrtc.tsx';
import { H265Direct } from './h265-direct.tsx';
import { H265Webrtc } from './h265-webrtc.tsx';
import { Mjpeg } from './mjpeg.tsx';

export const Screen = () => {
  const videoMode = useAtomValue(videoModeAtom);

  if (videoMode === 'mjpeg') {
    return <Mjpeg />;
  }

  if (videoMode === 'h264-direct') {
    return <H264Direct />;
  }

  if (videoMode === 'h265-webrtc') {
    return <H265Webrtc />;
  }

  if (videoMode === 'h265-direct') {
    return <H265Direct />;
  }

  return <H264Webrtc />;
};
