import { useEffect, useState } from 'react';
import { useAtomValue } from 'jotai';
import { MonitorIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/stream.ts';
import * as storage from '@/lib/localstorage.ts';
import { videoModeAtom } from '@/jotai/screen.ts';
import { MenuItem } from '@/components/menu-item.tsx';

import { BitRateMap, QualityMap } from './constants.ts';
import { Edid } from './edid.tsx';
import { Gop } from './gop.tsx';
import { Quality } from './quality';
import { VideoMode } from './video-mode.tsx';

export const Screen = () => {
  const { t } = useTranslation();

  const videoMode = useAtomValue(videoModeAtom);
  const [quality, setQuality] = useState(2);
  const [gop, setGop] = useState(50);

  useEffect(() => {
    if (!videoMode) return;

    api.setMode(videoMode);
    updateQuality();
    updateGop();
  }, [videoMode]);

  async function updateQuality() {
    const cookieQuality = storage.getQuality();
    if (!cookieQuality) return;

    const key = cookieQuality >= 1 && cookieQuality <= 4 ? cookieQuality : 2;
    const value = videoMode === 'mjpeg' ? QualityMap.get(key) : BitRateMap.get(key);
    if (!value) return;

    const rsp = await api.setQuality(value);
    if (rsp.code === 0) {
      setQuality(key);
    }
  }

  async function updateGop() {
    const cookieGop = storage.getGop();
    if (!cookieGop) return;

    const rsp = await api.setGop(cookieGop);
    if (rsp.code === 0) {
      setGop(cookieGop);
    }
  }

  const content = (
    <div className="flex flex-col space-y-1">
      <VideoMode />
      <Edid />
      <Quality quality={quality} setQuality={setQuality} />
      {videoMode !== 'mjpeg' && <Gop gop={gop} setGop={setGop} />}
    </div>
  );

  return <MenuItem title={t('screen.title')} icon={<MonitorIcon size={18} />} content={content} />;
};
