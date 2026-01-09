import { useState } from 'react';
import { Switch } from 'antd';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/virtual-device.ts';

interface NetworkProps {
  isDisabled: boolean;
  isNetworkEnabled: boolean;
  getVirtualDevices: () => void;
}

export const Network = ({ isDisabled, isNetworkEnabled, getVirtualDevices }: NetworkProps) => {
  const { t } = useTranslation();

  const [isLoading, setIsLoading] = useState(false);

  function update() {
    if (isLoading) return;
    setIsLoading(true);

    api
      .updateVirtualDevice('network')
      .then((rsp) => {
        if (rsp.code !== 0) {
          return;
        }

        getVirtualDevices();
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-col space-y-1">
        <span>{t('settings.device.network')}</span>
        <span className="text-xs text-neutral-500">{t('settings.device.networkDesc')}</span>
      </div>

      <Switch
        disabled={isDisabled}
        checked={isNetworkEnabled}
        loading={isLoading}
        onChange={update}
      />
    </div>
  );
};
