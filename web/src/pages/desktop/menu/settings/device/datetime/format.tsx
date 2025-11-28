import { useEffect, useState } from 'react';
import { Segmented } from 'antd';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/vm.ts';

type FormatProps = {
  format: string;
  setFormat: (format: string) => void;
};

export const Format = ({ format, setFormat }: FormatProps) => {
  const { t } = useTranslation();

  const [isLoading, setIsLoading] = useState<boolean>(false);

  useEffect(() => {
    setIsLoading(true);

    api
      .getLCDTimeFormat()
      .then((rsp) => {
        if (rsp.code === 0) {
          const value = rsp.data.format.toUpperCase();
          setFormat(value);
        }
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  function update(timeFormat: string) {
    if (isLoading) return;
    setIsLoading(true);

    api
      .setLCDTimeFormat(timeFormat)
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.error(rsp.msg);
          return;
        }
        setFormat(timeFormat);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  return (
    <div className="flex items-center justify-between">
      <span>{t('settings.device.datetime.format')}</span>

      <Segmented<string>
        options={['12H', '24H']}
        value={format}
        disabled={isLoading}
        onChange={update}
      />
    </div>
  );
};
