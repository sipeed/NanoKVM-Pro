import { atom } from 'jotai';

import { VideoMode, VideoParameters, VideoStatus } from '@/types';

export const videoModeAtom = atom<VideoMode | null>(null);

export const videoParametersAtom = atom<VideoParameters>({
  rateControlMode: 'vbr', // cbr | vbr
  bitrate: 8000, // 1000 - 20000
  gop: 50, // 1 - 200
  fps: 0, // 0 - 120
  scale: 1,
  quality: 80 // 1-100 (only for mjpeg)
});

export const videoStatusAtom = atom<VideoStatus>(VideoStatus.Normal);

export const videoVolumeAtom = atom(0);
