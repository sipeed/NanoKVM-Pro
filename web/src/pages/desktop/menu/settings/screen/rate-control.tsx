import { useState } from 'react';
import { Popover, Segmented } from 'antd';
import { useAtom } from 'jotai';
import { CircleHelpIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/stream.ts';
import * as storage from '@/lib/localstorage.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

const rateControlModes = [
  { label: 'CBR', value: 'cbr' },
  { label: 'VBR', value: 'vbr' }
];

export const RateControl = () => {
  const { t } = useTranslation();

  const [videoParameters, setVideoParameters] = useAtom(videoParametersAtom);
  const [isLoading, setIsLoading] = useState(false);

  async function update(mode: string) {
    if (isLoading || mode === videoParameters.rateControlMode) {
      return;
    }

    try {
      const rsp = await api.setRateControlMode(mode);
      if (rsp.code !== 0) {
        return;
      }

      const bitrate = mode === 'vbr' ? 8000 : 5000;
      const isSuccess = await updateBitrate(bitrate);

      const parameters = {
        ...videoParameters,
        rateControlMode: mode,
        bitrate: isSuccess ? bitrate : videoParameters.bitrate
      };

      setVideoParameters(parameters);
      storage.setVideoParameters(JSON.stringify(parameters));
    } catch (err) {
      console.log(err);
    } finally {
      setIsLoading(false);
    }
  }

  async function updateBitrate(bitrate: number) {
    const rsp = await api.setQuality(bitrate);
    return rsp.code === 0;
  }

  const help = (
    <ul className="max-w-[550px] list-outside">
      <li>{t('settings.screen.rateControlMode.cbr')}</li>
      <li>{t('settings.screen.rateControlMode.vbr')}</li>
    </ul>
  );

  return (
    <div className="flex items-center justify-between space-x-5">
      <div className="flex flex-1 flex-col space-y-1">
        <div className="flex items-center space-x-2">
          <span>{t('settings.screen.rateControlMode.title')}</span>

          <Popover className="cursor-pointer text-neutral-500" content={help}>
            <CircleHelpIcon size={15} />
          </Popover>
        </div>

        <span className="text-xs text-neutral-500">
          {t('settings.screen.rateControlMode.description')}
        </span>
      </div>

      <div className="w-[240px]">
        <Segmented<string>
          block
          options={rateControlModes}
          value={videoParameters.rateControlMode}
          disabled={isLoading}
          onChange={update}
        />
      </div>
    </div>
  );
};
