import { useState } from 'react';
import { Modal, Tooltip } from 'antd';
import clsx from 'clsx';
import { useAtom } from 'jotai';
import { MicIcon, MicOffIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/virtual-device.ts';
import { microphoneEnabledAtom } from '@/jotai/audio.ts';

export const Microphone = () => {
  const { t } = useTranslation();

  const [isEnabled, setIsEnabled] = useAtom(microphoneEnabledAtom);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const toggleMicrophone = async () => {
    if (isEnabled) {
      setIsEnabled(false);
      return;
    }

    try {
      const rsp = await api.getVirtualDevice();
      if (rsp.code === 0 && rsp.data) {
        if (rsp.data.isMicEnabled) {
          setIsEnabled(true);
        } else {
          setIsModalOpen(true);
        }
      }
    } catch (err) {
      console.log(err);
    }
  };

  const handleConfirm = async () => {
    if (isLoading) return;
    setIsLoading(true);

    try {
      const rsp = await api.updateVirtualDevice('mic');
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      setIsModalOpen(false);
      setIsEnabled(true);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    setIsModalOpen(false);
  };

  return (
    <>
      <Tooltip title={t('microphone.title')} placement="bottom" mouseEnterDelay={0.6}>
        <div
          className={clsx(
            'flex h-[30px] w-[30px] cursor-pointer items-center justify-center rounded text-neutral-300 hover:bg-neutral-700/80 hover:text-white',
            isLoading && 'pointer-events-none opacity-50'
          )}
          onClick={toggleMicrophone}
        >
          {isEnabled ? <MicIcon size={18} className="text-green-400" /> : <MicOffIcon size={18} />}
        </div>
      </Tooltip>

      <Modal
        title={t('microphone.notMounted')}
        open={isModalOpen}
        onOk={handleConfirm}
        onCancel={handleCancel}
        okText={t('microphone.confirm')}
        cancelText={t('microphone.cancel')}
        confirmLoading={isLoading}
      >
        <div className="flex flex-col space-y-3 py-3">
          <span>{t('microphone.tip1')}</span>
          <span>{t('microphone.tip2')}</span>
        </div>
      </Modal>
    </>
  );
};
