import { ReactElement } from 'react';
import { Popover } from 'antd';
import { useAtom } from 'jotai';
import { CheckIcon, PercentIcon, ScalingIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as storage from '@/lib/localstorage.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

const ScaleList = [
  { label: '200', value: 2 },
  { label: '150', value: 1.5 },
  { label: '100', value: 1 },
  { label: '75', value: 0.75 },
  { label: '50', value: 0.5 }
];

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
    <>
      {ScaleList.map((scale) => (
        <div
          key={scale.value}
          className="flex h-[30px] cursor-pointer select-none items-center rounded pl-1 pr-5 hover:bg-neutral-700/70"
          onClick={() => update(scale.value)}
        >
          <div className="flex h-[14px] w-[20px] items-end text-blue-500">
            {scale.value === videoParameters.scale && <CheckIcon size={14} />}
          </div>
          <div className="flex items-center space-x-0.5">
            <span>{scale.label}</span>
            <PercentIcon size={12} />
          </div>
        </div>
      ))}
    </>
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
