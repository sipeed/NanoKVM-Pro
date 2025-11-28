import { useEffect, useState } from 'react';
import { Switch, Tooltip } from 'antd';
import { CircleAlertIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { getHidMode } from '@/api/hid.ts';
import { getMountedImage } from '@/api/storage.ts';
import * as api from '@/api/virtual-device.ts';

export const VirtualDevices = () => {
  const { t } = useTranslation();
  const [isDiskOn, setIsDiskOn] = useState(false);
  const [isNetworkOn, setIsNetworkOn] = useState(false);
  const [isSdCardExist, setIsSdCardExist] = useState(false);
  const [isHidOnlyMode, setIsHidOnlyMode] = useState(true);
  const [isImageMounted, setIsImageMounted] = useState(true);
  const [loading, setLoading] = useState<'' | 'disk' | 'network'>('');

  useEffect(() => {
    api.getVirtualDevice().then((rsp) => {
      if (rsp.code === 0) {
        setIsDiskOn(rsp.data.disk);
        setIsNetworkOn(rsp.data.network);
        setIsSdCardExist(rsp.data.sdCard);
      }
    });

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

  function update(device: 'network' | 'disk') {
    if (loading) return;
    setLoading(device);

    api
      .updateVirtualDevice(device)
      .then((rsp) => {
        if (rsp.code !== 0) {
          return;
        }

        if (device === 'disk') {
          setIsDiskOn(rsp.data.on);
        }
        if (device === 'network') {
          setIsNetworkOn(rsp.data.on);
        }
      })
      .finally(() => {
        setLoading('');
      });
  }

  return (
    <>
      {/* virtual disk */}
      <div className="flex items-center justify-between">
        <div className="flex flex-col">
          <div className="flex items-center space-x-2">
            <span>{t('settings.device.disk')}</span>
            {!isHidOnlyMode && !isSdCardExist && (
              <Tooltip
                title={t('settings.device.noSdCard')}
                className="cursor-pointer text-neutral-500"
                placement="bottom"
                overlayStyle={{ maxWidth: '300px' }}
              >
                <CircleAlertIcon size={15} />
              </Tooltip>
            )}
          </div>

          <span className="text-xs text-neutral-500">{t('settings.device.diskDesc')}</span>
        </div>

        <Switch
          disabled={isHidOnlyMode || isImageMounted || (!isDiskOn && !isSdCardExist)}
          checked={isDiskOn}
          loading={loading === 'disk'}
          onChange={() => update('disk')}
        />
      </div>

      {/* virtual network */}
      <div className="flex items-center justify-between">
        <div className="flex flex-col">
          <span>{t('settings.device.network')}</span>
          <span className="text-xs text-neutral-500">{t('settings.device.networkDesc')}</span>
        </div>

        <Switch
          disabled={isHidOnlyMode || isImageMounted}
          checked={isNetworkOn}
          loading={loading === 'network'}
          onChange={() => update('network')}
        />
      </div>
    </>
  );
};
