import { useState } from 'react';
import { Modal, Segmented, Switch } from 'antd';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/virtual-device.ts';

interface DiskProps {
  isDisabled: boolean;
  isEmmcExist: boolean;
  isSdCardExist: boolean;
  mountedDisk: string;
  getVirtualDevices: () => void;
}

export const Disk = ({
  isDisabled,
  isEmmcExist,
  isSdCardExist,
  mountedDisk,
  getVirtualDevices
}: DiskProps) => {
  const { t } = useTranslation();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  async function update(type: string, confirm?: boolean) {
    if (!type && !mountedDisk) {
      return;
    }

    // mounting eMMC requires a second confirmation
    if (type === 'emmc' && mountedDisk !== 'emmc' && !isEmmcExist && !confirm) {
      setIsModalOpen(true);
      return;
    }

    const typ = type ? type : mountedDisk;
    await updateVirtualDevice(typ);

    await getVirtualDevices();
  }

  async function updateVirtualDevice(type: string) {
    if (isLoading) return;

    setIsLoading(true);
    setIsModalOpen(false);

    try {
      const rsp = await api.updateVirtualDevice('disk', type);
      if (rsp.code !== 0) {
        return;
      }
    } catch (err) {
      console.log(err);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <>
      <div className="flex items-center justify-between">
        <div className="flex flex-col space-y-1">
          <span>{t('settings.device.disk')}</span>
          <span className="text-xs text-neutral-500">{t('settings.device.diskDesc')}</span>
        </div>

        {isSdCardExist ? (
          // sdcard and emmc
          <Segmented
            options={[
              { value: '', label: 'Disable' },
              { value: 'emmc', label: 'eMMC' },
              { value: 'sdcard', label: 'SD Card' }
            ]}
            disabled={isDisabled || isLoading}
            value={mountedDisk}
            onChange={(value) => update(value)}
          />
        ) : (
          // only emmc
          <Switch
            disabled={isDisabled}
            checked={!!mountedDisk}
            loading={isLoading}
            onChange={() => update('emmc')}
          />
        )}
      </div>

      {/* second confirmation */}
      <Modal
        title={t('settings.device.emmc.warning')}
        open={isModalOpen}
        centered={true}
        okText={t('settings.device.okBtn')}
        cancelText={t('settings.device.cancelBtn')}
        onOk={() => update('emmc', true)}
        onCancel={() => setIsModalOpen(false)}
      >
        <p>{t('settings.device.emmc.tip1')}</p>
        <p className="text-red-500">{t('settings.device.emmc.tip2')}</p>
      </Modal>
    </>
  );
};
