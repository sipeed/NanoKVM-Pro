import { Select } from 'antd';
import { useAtom } from 'jotai';
import { useTranslation } from 'react-i18next';

import type { VideoMode as TVideoMode } from '@/types';
import * as storage from '@/lib/localstorage.ts';
import { videoModeAtom } from '@/jotai/screen.ts';

const videoModes = [
  { value: 'h264-webrtc', label: 'H.264 WebRTC' },
  { value: 'h264-direct', label: 'H.264 Direct ' },
  { value: 'mjpeg', label: 'MJPEG' }
];

export const VideoMode = () => {
  const { t } = useTranslation();

  const [videoMode, setVideoMode] = useAtom(videoModeAtom);

  function update(value: string) {
    if (value === videoMode) {
      return;
    }

    setVideoMode(value as TVideoMode);
    storage.setVideoMode(value);

    setTimeout(() => {
      window.location.reload();
    }, 500);
  }

  return (
    <div className="flex items-center justify-between space-x-5">
      <div className="flex flex-1 flex-col space-y-1">
        <span>{t('settings.screen.videoMode.title')}</span>
        <span className="text-xs text-neutral-500">
          {t('settings.screen.videoMode.description')}
        </span>
      </div>

      <Select value={videoMode} style={{ width: 240 }} options={videoModes} onSelect={update} />
    </div>
  );
};
