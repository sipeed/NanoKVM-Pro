import { Image } from 'antd';
import clsx from 'clsx';
import { useAtomValue } from 'jotai';

import MonitorXIcon from '@/assets/images/monitor-x.svg';
import { getBaseUrl } from '@/lib/service.ts';
import { mouseStyleAtom } from '@/jotai/mouse.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

export const Mjpeg = () => {
  const videoParameters = useAtomValue(videoParametersAtom);
  const mouseStyle = useAtomValue(mouseStyleAtom);

  return (
    <div className="flex h-screen w-screen items-start justify-center xl:items-center">
      <Image
        id="screen"
        className={clsx(
          'block max-h-screen min-h-[50vh] min-w-[50vw] select-none object-scale-down',
          mouseStyle
        )}
        style={{ transform: `scale(${videoParameters.scale})` }}
        src={`${getBaseUrl('http')}/api/stream/mjpeg`}
        fallback={MonitorXIcon}
        preview={false}
      />
    </div>
  );
};
