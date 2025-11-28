import { useAtomValue } from 'jotai';

import { videoModeAtom } from '@/jotai/screen.ts';

import { Volume } from './volume.tsx';

export const Audio = () => {
  const videoMode = useAtomValue(videoModeAtom);

  return <>{videoMode?.endsWith('webrtc') && <Volume />}</>;
};
