import { ReactElement } from 'react';
import { Popover } from 'antd';
import { useAtom } from 'jotai';
import { CheckIcon, RotateCwIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { VideoRotation } from '@/types';
import * as storage from '@/lib/localstorage.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

const RotationList: VideoRotation[] = [0, 90, 180, 270];

export const Rotation = (): ReactElement => {
  const { t } = useTranslation();
  const [videoParameters, setVideoParameters] = useAtom(videoParametersAtom);

  function update(rotation: VideoRotation) {
    const parameters = { ...videoParameters, rotation };
    setVideoParameters(parameters);
    storage.setVideoParameters(JSON.stringify(parameters));
  }

  const content = (
    <>
      {RotationList.map((rotation) => (
        <div
          key={rotation}
          className="flex h-[30px] cursor-pointer select-none items-center rounded pl-1 pr-5 hover:bg-neutral-700/70"
          onClick={() => update(rotation)}
        >
          <div className="flex h-[14px] w-[20px] items-end text-blue-500">
            {rotation === videoParameters.rotation && <CheckIcon size={14} />}
          </div>
          <span>{rotation}&deg;</span>
        </div>
      ))}
    </>
  );

  return (
    <Popover content={content} placement="rightTop" arrow={false} align={{ offset: [13, 0] }}>
      <div className="flex h-[30px] cursor-pointer items-center space-x-1 rounded pl-3 pr-6 text-neutral-300 hover:bg-neutral-700/50">
        <div className="flex h-[14px] w-[20px] items-end">
          <RotateCwIcon size={16} />
        </div>
        <span>{t('screen.rotation')}</span>
      </div>
    </Popover>
  );
};
