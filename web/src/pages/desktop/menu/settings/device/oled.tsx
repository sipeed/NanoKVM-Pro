import { useEffect, useState } from 'react';
import { Select } from 'antd';
import clsx from 'clsx';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/vm.ts';

export const Oled = ({ disable = false }: { disable?: boolean }) => {
  const { t } = useTranslation();
  const [isLoading, setIsLoading] = useState(false);
  const [screenType, SetScreenType] = useState('');
  const [sleep, setSleep] = useState('');

  useEffect(() => {
    api.getOLED().then((rsp) => {
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      if (!rsp.data.exist) {
        return;
      }

      SetScreenType(rsp.data.type);
      setSleep(rsp.data.sleep.toString());
    });
  }, []);

  const options = [0, 15, 30, 60, 180, 300, 600, 1800, 3600].map((duration) => ({
    value: duration.toString(),
    label: t(`settings.device.oled.${duration}`)
  }));

  function update(value: string) {
    if (isLoading) return;
    setIsLoading(true);

    api
      .setOLED(parseInt(value))
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
          return;
        }

        setSleep(value);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  if (screenType !== 'ATX') {
    return <></>;
  }

  return (
    <div
      className={clsx(
        'flex items-center justify-between',
        disable && 'cursor-not-allowed opacity-50',
        'pt-3',
        isLoading ? 'cursor-pointer' : 'cursor-not-allowed'
      )}
    >
      <div className="flex flex-col space-y-1">
        <span>{t('settings.device.oled.title')}</span>
        <span className="text-xs text-neutral-500">{t('settings.device.oled.description')}</span>
      </div>

      <Select
        style={{ width: 150 }}
        value={sleep}
        options={options}
        loading={isLoading}
        onChange={update}
      />
    </div>
  );
};
