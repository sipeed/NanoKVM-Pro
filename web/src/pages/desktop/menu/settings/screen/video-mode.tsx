import { useEffect, useState } from 'react';
import { Select } from 'antd';
import { useAtom } from 'jotai';
import Cookies from 'js-cookie';
import { useTranslation } from 'react-i18next';

import type { VideoMode as TVideoMode } from '@/types';
import * as storage from '@/lib/localstorage.ts';
import { getSupportedVideoModes } from '@/lib/video.ts';
import { videoModeAtom } from '@/jotai/screen.ts';

export const VideoMode = () => {
  const { t } = useTranslation();

  const [videoMode, setVideoMode] = useAtom(videoModeAtom);

  const [isSupported, setIsSupported] = useState(false);
  const [isEnabled, setIsEnabled] = useState(false);

  const videoModes = [
    { value: 'h265-webrtc', label: 'H.265 WebRTC', disabled: !isSupported || !isEnabled },
    { value: 'h265-direct', label: 'H.265 Direct', disabled: !isSupported || !isEnabled },
    { value: 'h264-webrtc', label: 'H.264 WebRTC' },
    { value: 'h264-direct', label: 'H.264 Direct ' },
    { value: 'mjpeg', label: 'MJPEG' }
  ];

  useEffect(() => {
    const videoModes = getSupportedVideoModes();
    const supported = videoModes.some((mode) => mode.startsWith('h265'));
    setIsSupported(supported);

    const h265 = Cookies.get('h265');
    if (h265) {
      setIsEnabled(true);
    }
  }, []);

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
      <div className="flex flex-1 flex-col">
        <span>{t('settings.screen.videoMode.title')}</span>
        <span className="text-xs text-neutral-500">
          {t('settings.screen.videoMode.description')}
        </span>
      </div>

      <Select value={videoMode} style={{ width: 240 }} options={videoModes} onSelect={update} />
    </div>
  );
};
