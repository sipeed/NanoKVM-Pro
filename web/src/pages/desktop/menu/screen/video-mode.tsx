import { Divider, Popover } from 'antd';
import clsx from 'clsx';
import { useAtomValue } from 'jotai';
import { CheckIcon, TvMinimalPlayIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { setVideoMode as setCookie } from '@/lib/localstorage.ts';
import { getSupportedVideoModes } from '@/lib/video.ts';
import { videoModeAtom } from '@/jotai/screen.ts';

const videoGroups = [
  {
    key: 'h264',
    name: 'H.264',
    modes: [
      { key: 'h264-webrtc', name: 'H.264 WebRTC' },
      { key: 'h264-direct', name: 'H.264 Direct ' }
    ]
  },
  {
    key: 'mjpeg',
    name: 'MJPEG',
    modes: [{ key: 'mjpeg', name: 'MJPEG' }]
  }
];

export const VideoMode = () => {
  const { t } = useTranslation();
  const videoMode = useAtomValue(videoModeAtom);

  const supportedVideoModes = getSupportedVideoModes();

  function update(mode: string) {
    if (mode === videoMode || !supportedVideoModes.includes(mode)) return;

    setCookie(mode);

    // reload after changing video mode
    setTimeout(() => {
      window.location.reload();
    }, 500);
  }

  const content = (
    <>
      {videoGroups.map((group) => (
        <div key={group.key}>
          {group.modes.map((mode) => (
            <div
              key={mode.key}
              className={clsx(
                'flex select-none items-center rounded py-1.5 pl-1 pr-5 hover:bg-neutral-700/70',
                supportedVideoModes.includes(mode.key)
                  ? 'cursor-pointer'
                  : 'cursor-not-allowed text-neutral-500'
              )}
              onClick={() => update(mode.key)}
            >
              <div className="flex h-[14px] w-[20px] items-end text-blue-500">
                {mode.key === videoMode && <CheckIcon size={15} />}
              </div>
              <span>{mode.name}</span>
            </div>
          ))}

          {group.key !== 'mjpeg' && <Divider style={{ margin: '5px 0' }} />}
        </div>
      ))}
    </>
  );

  return (
    <Popover content={content} placement="rightTop" arrow={false} align={{ offset: [14, 0] }}>
      <div className="flex h-[30px] cursor-pointer items-center space-x-2 rounded pl-3 pr-6 text-neutral-300 hover:bg-neutral-700/70">
        <TvMinimalPlayIcon size={18} />
        <span className="select-none text-sm">{t('screen.video')}</span>
      </div>
    </Popover>
  );
};
