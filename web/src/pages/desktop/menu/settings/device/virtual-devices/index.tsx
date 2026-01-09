import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getHidMode } from '@/api/hid.ts';
import { getMountedImage } from '@/api/storage.ts';
import * as api from '@/api/virtual-device.ts';

import { Disk } from './disk.tsx';
import { Microphone } from './microphone.tsx';
import { Network } from './network.tsx';

export const VirtualDevices = () => {
  const { t } = useTranslation();

  const [isHidOnlyMode, setIsHidOnlyMode] = useState(false);
  const [isImageMounted, setIsImageMounted] = useState(false);

  const [isNetworkEnabled, setIsNetworkEnabled] = useState(false);
  const [isMicEnabled, setIsMicEnabled] = useState(false);

  const [mountedDisk, setMountedDisk] = useState('');
  const [isSdCardExist, setIsSdCardExist] = useState(false);
  const [isEmmcExist, setIsEmmcExist] = useState(false);

  useEffect(() => {
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

    getVirtualDevices();
  }, []);

  async function getVirtualDevices() {
    const rsp = await api.getVirtualDevice();
    if (rsp.code !== 0 || !rsp.data) {
      return;
    }

    setIsNetworkEnabled(rsp.data.isNetworkEnabled);
    setIsMicEnabled(rsp.data.isMicEnabled);

    setMountedDisk(rsp.data.mountedDisk);
    setIsEmmcExist(rsp.data.isEmmcExist);
    setIsSdCardExist(rsp.data.isSdCardExist);
  }

  return (
    <>
      {isHidOnlyMode ? (
        <div className="flex flex-col space-y-1">
          <span className="text-neutral-400">{t('settings.device.hidOnly')}</span>
          <span className="text-xs text-neutral-500">{t('settings.device.hidOnlyDesc')}</span>
        </div>
      ) : (
        isImageMounted && (
          <div className="flex flex-col space-y-1">
            <span className="text-neutral-400">{t('settings.device.imageMounted')}</span>
            <span className="text-xs text-neutral-500">
              {t('settings.device.imageMountedDesc')}
            </span>
          </div>
        )
      )}

      <Disk
        isDisabled={isHidOnlyMode || isImageMounted}
        isEmmcExist={isEmmcExist}
        isSdCardExist={isSdCardExist}
        mountedDisk={mountedDisk}
        getVirtualDevices={getVirtualDevices}
      />

      <Network
        isDisabled={isHidOnlyMode || isImageMounted}
        isNetworkEnabled={isNetworkEnabled}
        getVirtualDevices={getVirtualDevices}
      />

      <Microphone
        isDisabled={isHidOnlyMode || isImageMounted}
        isMicEnabled={isMicEnabled}
        getVirtualDevices={getVirtualDevices}
      />
    </>
  );
};
