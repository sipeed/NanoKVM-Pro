import { Image } from 'antd';
import clsx from 'clsx';
import { useAtomValue } from 'jotai';

import MonitorXIcon from '@/assets/images/monitor-x.svg';
import { getBaseUrl } from '@/lib/service.ts';
import { getVideoTransformStyle, isQuarterTurn } from '@/lib/video-transform.ts';
import { mouseStyleAtom } from '@/jotai/mouse.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

export const Mjpeg = () => {
  const videoParameters = useAtomValue(videoParametersAtom);
  const mouseStyle = useAtomValue(mouseStyleAtom);

  return (
    <div
      className={clsx(
        'flex h-screen w-screen justify-center',
        isQuarterTurn(videoParameters.rotation) ? 'items-center' : 'items-start xl:items-center'
      )}
    >
      <Image
        id="screen"
        className={clsx(
          'block max-h-screen min-h-[50vh] min-w-[50vw] select-none object-contain',
          mouseStyle
        )}
        style={getVideoTransformStyle(videoParameters.scale, videoParameters.rotation)}
        src={`${getBaseUrl('http')}/api/stream/mjpeg`}
        fallback={MonitorXIcon}
        preview={false}
      />
    </div>
  );
};
