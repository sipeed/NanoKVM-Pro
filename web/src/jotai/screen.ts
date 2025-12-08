import { atom } from 'jotai';

import { VideoMode, VideoStatus } from '@/types';

export const videoModeAtom = atom<VideoMode | null>(null);

export const videoScaleAtom = atom<number>(1.0);

export const videoStatusAtom = atom<VideoStatus>(VideoStatus.Normal);

export const videoVolumeAtom = atom(0);

export const fullscreenAtom = atom(false);
