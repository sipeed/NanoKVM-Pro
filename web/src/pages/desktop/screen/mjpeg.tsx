import { Image } from 'antd';
import clsx from 'clsx';
import { useAtomValue } from 'jotai';

import MonitorXIcon from '@/assets/images/monitor-x.svg';
import { getBaseUrl } from '@/lib/service.ts';
import { mouseStyleAtom } from '@/jotai/mouse.ts';
import { videoScaleAtom } from '@/jotai/screen.ts';

export const Mjpeg = () => {
  const mouseStyle = useAtomValue(mouseStyleAtom);
  const videoScale = useAtomValue(videoScaleAtom);

  return (
    <div className="flex h-screen w-screen items-start justify-center xl:items-center origin-center"
         style={{
            transform: `scale(${videoScale})`,
          }}>
      <Image
        id="screen"
        className={clsx(
          'block max-h-screen min-h-[480px] min-w-[640px] select-none object-scale-down',
          mouseStyle
        )}
        src={`${getBaseUrl('http')}/api/stream/mjpeg`}
        fallback={MonitorXIcon}
        preview={false}
      />
    </div>
  );
};
