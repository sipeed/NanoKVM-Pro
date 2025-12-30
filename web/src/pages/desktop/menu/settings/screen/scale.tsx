import { Slider } from 'antd';
import { useAtom } from 'jotai';
import { useTranslation } from 'react-i18next';

import * as storage from '@/lib/localstorage.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

export const Scale = () => {
  const { t } = useTranslation();
  const [videoParameters, setVideoParameters] = useAtom(videoParametersAtom);

  async function update(scale: number) {
    const parameters = { ...videoParameters, scale };

    setVideoParameters(parameters);
    storage.setVideoParameters(JSON.stringify(parameters));
  }

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-1 flex-col space-y-1">
        <span>{t('settings.screen.scale.title')}</span>
        <span className="text-xs text-neutral-500">{t('settings.screen.scale.description')}</span>
      </div>

      <div className="w-[240px]">
        <Slider
          value={videoParameters.scale}
          min={0.5}
          max={2}
          step={0.1}
          marks={{ 2: '2' }}
          onChange={update}
        />
      </div>
    </div>
  );
};
