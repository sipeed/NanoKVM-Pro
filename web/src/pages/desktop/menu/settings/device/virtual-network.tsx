import { useEffect, useState } from 'react';
import { Switch } from 'antd';
import { useTranslation } from 'react-i18next';

import { getHidMode } from '@/api/hid.ts';
import { getMountedImage } from '@/api/storage.ts';
import * as api from '@/api/virtual-device.ts';

export const VirtualNetwork = () => {
  const { t } = useTranslation();

  const [isNetworkOn, setIsNetworkOn] = useState(false);
  const [isHidOnlyMode, setIsHidOnlyMode] = useState(false);
  const [isImageMounted, setIsImageMounted] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    getVirtualDevices();

    getHidMode().then((rsp) => {
      if (rsp.code === 0) {
        setIsHidOnlyMode(rsp.data.mode === 'hid-only');
      }
    });

    getMountedImage().then((rsp) => {
      if (rsp.code === 0) {
        setIsImageMounted(!!rsp.data?.file);
      }
    });
  }, []);

  async function getVirtualDevices() {
    const rsp = await api.getVirtualDevice();
    if (rsp.code !== 0) {
      return;
    }
    setIsNetworkOn(rsp.data.isNetworkEnabled);
  }

  function update() {
    if (isLoading) return;
    setIsLoading(true);

    api
      .updateVirtualDevice('network', '')
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
        disabled={isHidOnlyMode || isImageMounted}
        checked={isNetworkOn}
        loading={isLoading}
        onChange={update}
      />
    </div>
  );
};
