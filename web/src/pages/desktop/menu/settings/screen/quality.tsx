import { Slider } from 'antd';
import { useAtom } from 'jotai';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/stream.ts';
import * as storage from '@/lib/localstorage.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

export const Quality = () => {
  const { t } = useTranslation();
  const [videoParameters, setVideoParameters] = useAtom(videoParametersAtom);

  async function update(value: number) {
    if (value === videoParameters.quality) {
      return;
    }

    try {
      const rsp = await api.setQuality(value);
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      const parameters = { ...videoParameters, quality: value };
      setVideoParameters(parameters);
      storage.setVideoParameters(JSON.stringify(parameters));
    } catch (err) {
      console.log(err);
    }
  }

  return (
    <div className="flex items-center justify-between space-x-5">
      <div className="flex flex-1 flex-col">
        <span>{t('settings.screen.quality.title')}</span>
        <span className="text-xs text-neutral-500">{t('settings.screen.quality.description')}</span>
      </div>

      <div className="w-[240px]">
        <Slider
          value={videoParameters.quality}
          min={1}
          max={100}
          marks={{ 100: '100' }}
          onChange={update}
        />
      </div>
    </div>
  );
};
