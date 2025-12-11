import { ReactElement } from 'react';
import { Popover, Slider } from 'antd';
import { useAtom } from 'jotai';
import { ScalingIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as storage from '@/lib/localstorage.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

export const Scale = (): ReactElement => {
  const { t } = useTranslation();

  const [videoParameters, setVideoParameters] = useAtom(videoParametersAtom);

  async function update(scale: number): Promise<void> {
    if (scale < 0.5 || scale > 2) {
      return;
    }

    const parameters = { ...videoParameters, scale };
    setVideoParameters(parameters);
    storage.setVideoParameters(JSON.stringify(parameters));
  }

  const content = (
    <div className="m-3 h-[180px] w-[70px]">
      <Slider
        vertical
        marks={{
          0.5: <span>x0.5</span>,
          1: <span>x1.0</span>,
          1.5: <span>x1.5</span>,
          2: <span>x2.0</span>
        }}
        range={false}
        included={false}
        min={0.5}
        max={2}
        step={0.1}
        defaultValue={videoParameters.scale}
        onChange={update}
      />
    </div>
  );

  return (
    <Popover content={content} placement="rightTop" arrow={false} align={{ offset: [13, 0] }}>
      <div className="flex h-[30px] cursor-pointer items-center space-x-1 rounded pl-3 pr-6 text-neutral-300 hover:bg-neutral-700/50">
        <div className="flex h-[14px] w-[20px] items-end">
          <ScalingIcon size={16} />
        </div>
        <span>{t('screen.scale')}</span>
      </div>
    </Popover>
  );
};
