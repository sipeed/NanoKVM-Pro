import { useEffect, useState } from 'react';
import { Divider, Popover } from 'antd';
import { useAtom } from 'jotai';
import { CheckIcon, SquareActivityIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/stream.ts';
import * as storage from '@/lib/localstorage.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

export const Bitrate = () => {
  const { t } = useTranslation();
  const [videoParameters, setVideoParameters] = useAtom(videoParametersAtom);

  const [customBitrate, setCustomBitrate] = useState(0);
  const [isLoading, setIsLoading] = useState(false);

  const bitrateList = [
    { value: 0, label: t('screen.qualityAuto') },
    { value: 10000, label: t('screen.qualityLossless') },
    { value: 5000, label: t('screen.qualityHigh') },
    { value: 3000, label: t('screen.qualityMedium') },
    { value: 1000, label: t('screen.qualityLow') }
  ];

  useEffect(() => {
    if (videoParameters.rateControlMode === 'vbr') {
      setCustomBitrate(0);
      return;
    }

    const isExist = bitrateList.some((bitrate) => bitrate.value === videoParameters.bitrate);
    if (isExist) {
      setCustomBitrate(0);
      return;
    }

    setCustomBitrate(Math.floor(videoParameters.bitrate / 100));
  }, [videoParameters]);

  async function update(value: number) {
    if (isLoading) {
      return;
    }
    setIsLoading(true);

    const mode = value === 0 ? 'vbr' : 'cbr';
    const bitrate = value === 0 ? 8000 : value;

    try {
      let rsp = await api.setRateControlMode(mode);
      if (rsp.code !== 0) {
        return;
      }

      rsp = await api.setQuality(bitrate);
      if (rsp.code !== 0) {
        return;
      }

      const parameters = { ...videoParameters, rateControlMode: mode, bitrate: bitrate };
      setVideoParameters(parameters);
      storage.setVideoParameters(JSON.stringify(parameters));
    } catch (err) {
      console.log(err);
    } finally {
      setIsLoading(false);
    }
  }

  function isChecked(value: number) {
    if (videoParameters.rateControlMode === 'vbr') {
      return value === 0;
    }
    return value === videoParameters.bitrate;
  }

  const content = (
    <>
      {bitrateList.map((bitrate) => (
        <div
          key={bitrate.value}
          className="flex h-[30px] cursor-pointer select-none items-center rounded pl-1 pr-5 hover:bg-neutral-700/70"
          onClick={() => update(bitrate.value)}
        >
          <div className="flex h-[14px] w-[20px] items-end text-blue-500">
            {isChecked(bitrate.value) && <CheckIcon size={14} />}
          </div>
          <span>{bitrate.label}</span>
        </div>
      ))}

      {customBitrate > 0 && (
        <>
          <Divider style={{ margin: '5px 0' }} />
          <div className="flex h-[30px] cursor-pointer select-none items-center rounded pl-1 pr-5 hover:bg-neutral-700/70">
            <div className="flex h-[14px] w-[20px] items-end text-blue-500">
              <CheckIcon size={14} />
            </div>
            <span>{customBitrate}%</span>
          </div>
        </>
      )}
    </>
  );

  return (
    <Popover content={content} placement="rightTop" arrow={false} align={{ offset: [14, 0] }}>
      <div className="flex h-[30px] cursor-pointer items-center space-x-2 rounded pl-3 pr-6 text-neutral-300 hover:bg-neutral-700/70">
        <SquareActivityIcon size={18} />
        <span className="select-none text-sm">{t('screen.quality')}</span>
      </div>
    </Popover>
  );
};
