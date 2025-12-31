import { Divider } from 'antd';
import { useAtomValue } from 'jotai';
import { useTranslation } from 'react-i18next';

import { videoModeAtom } from '@/jotai/screen.ts';

import { Bitrate } from './bitrate.tsx';
import { Edid } from './edid.tsx';
import { Fps } from './fps.tsx';
import { Gop } from './gop.tsx';
import { Quality } from './quality.tsx';
import { RateControl } from './rate-control.tsx';
import { Scale } from './scale.tsx';
import { VideoMode } from './video-mode.tsx';

export const Screen = () => {
  const { t } = useTranslation();

  const videoMode = useAtomValue(videoModeAtom);

  return (
    <>
      <div className="text-base font-bold">{t('settings.screen.title')}</div>
      <Divider className="opacity-50" />

      <div className="flex flex-col space-y-8">
        <VideoMode />
        <Edid />

        {videoMode !== 'mjpeg' && (
          <>
            <RateControl />
            <Bitrate />
            <Gop />
          </>
        )}

        <Fps />

        {videoMode === 'mjpeg' && <Quality />}

        <Scale />
      </div>
    </>
  );
};
