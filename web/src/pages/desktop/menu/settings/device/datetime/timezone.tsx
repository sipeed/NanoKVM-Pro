import { useState } from 'react';
import { Select } from 'antd';
import { useTranslation } from 'react-i18next';
import { allTimezones, useTimezoneSelect } from 'react-timezone-select';

import * as api from '@/api/vm.ts';

type TimezoneProps = {
  timezone: string;
  setTimezone: (timezone: string) => void;
};

const labelStyle = 'original';
const timezones = {
  ...allTimezones
};

export const Timezone = ({ timezone, setTimezone }: TimezoneProps) => {
  const { t } = useTranslation();

  const { options } = useTimezoneSelect({ labelStyle, timezones });
  const [isLoading, setIsLoading] = useState(false);
  const [status, setStatus] = useState<'' | 'error'>('');

  function update(value: string) {
    if (isLoading || timezone === value) {
      return;
    }

    setIsLoading(true);

    api
      .setTimezone(value)
      .then((rsp) => {
        if (rsp.code !== 0) {
          setStatus('error');
          console.log(rsp.msg);
          return;
        }

        setStatus('');
        setTimezone(value);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  return (
    <>
      <div className="flex items-center justify-between">
        <span>{t('settings.device.datetime.timezone')}</span>

        <Select
          value={timezone}
          loading={isLoading}
          status={status}
          style={{ width: 250, maxWidth: '50%' }}
          options={options}
          onSelect={update}
        />
      </div>
    </>
  );
};
