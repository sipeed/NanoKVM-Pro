import { Popover, Slider } from 'antd';
import { useAtom } from 'jotai';
import { CircleHelpIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/stream.ts';
import * as storage from '@/lib/localstorage.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

export const Bitrate = () => {
  const { t } = useTranslation();

  const [videoParameters, setVideoParameters] = useAtom(videoParametersAtom);

  async function update(value: number) {
    if (value === videoParameters.bitrate) {
      return;
    }

    try {
      const rsp = await api.setQuality(value);
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      const parameters = { ...videoParameters, bitrate: value };
      setVideoParameters(parameters);
      storage.setVideoParameters(JSON.stringify(parameters));
    } catch (err) {
      console.log(err);
    }
  }

  const help = (
    <ul className="max-w-[550px] list-outside">
      <li>{t('settings.screen.bitrate.lower')}</li>
      <li>{t('settings.screen.bitrate.higher')}</li>
    </ul>
  );

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-1 flex-col">
        <div className="flex items-center space-x-2">
          <span>{t('settings.screen.bitrate.title')}</span>

          <Popover className="cursor-pointer text-neutral-500" content={help}>
            <CircleHelpIcon size={15} />
          </Popover>
        </div>

        <span className="text-xs text-neutral-500">{t('settings.screen.bitrate.description')}</span>
      </div>

      <div className="w-[240px]">
        <Slider
          value={videoParameters.bitrate}
          min={1000}
          max={20000}
          marks={{ 20000: '20000' }}
          onChange={update}
        />
      </div>
    </div>
  );
};
