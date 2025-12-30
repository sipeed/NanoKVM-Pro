import { useEffect, useState } from 'react';
import { Button, Modal, Switch } from 'antd';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/vm.ts';

export const LowPower = () => {
  const { t } = useTranslation();

  const [isEnabled, setIsEnabled] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    getLowPowerMode();
  }, []);

  async function getLowPowerMode() {
    try {
      const rsp = await api.getLowPowerMode();
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }
      setIsEnabled(rsp.data.enabled);
    } catch (err) {
      console.log(err);
    } finally {
      setIsLoading(false);
    }
  }

  async function setLowPowerMode(enable: boolean) {
    setIsModalOpen(false);

    if (isLoading) return;
    setIsLoading(true);

    try {
      const rsp = await api.setLowPowerMode(enable);
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }
      setIsEnabled(enable);
    } catch (err) {
      console.log(err);
    } finally {
      setIsLoading(false);
    }
  }

  function update(enable: boolean) {
    if (isLoading) return;

    if (enable) {
      setIsModalOpen(true);
      return;
    }

    setLowPowerMode(enable);
  }

  const footer = [
    <Button key="cancel" onClick={() => setIsModalOpen(false)}>
      {t('settings.device.lowPower.cancelBtn')}
    </Button>,
    <Button key="ok" type="primary" loading={isLoading} onClick={() => setLowPowerMode(true)}>
      {t('settings.device.lowPower.okBtn')}
    </Button>
  ];

  return (
    <>
      <div className="flex items-center justify-between">
        <span>{t('settings.device.lowPower.title')}</span>

        <Switch value={isEnabled} loading={isLoading} onChange={update} />
      </div>

      <Modal
        title={t('settings.device.lowPower.tips')}
        open={isModalOpen}
        closable={false}
        centered={true}
        footer={footer}
        onCancel={() => setIsModalOpen(false)}
      >
        <div className="flex flex-col space-y-3 py-3">
          <span>{t('settings.device.lowPower.tip1')}</span>
          <span>{t('settings.device.lowPower.tip2')}</span>
        </div>
      </Modal>
    </>
  );
};
