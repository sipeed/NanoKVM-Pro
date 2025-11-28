import { useEffect, useState } from 'react';
import { Button, Modal, Switch } from 'antd';
import TextArea from 'antd/es/input/TextArea';
import clsx from 'clsx';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/network.ts';

type StaticIpConfig = {
  isEnabled: boolean;
  content: string;
};

type Result = {
  type: 'success' | 'failed';
  message: string;
};

export const StaticIp = () => {
  const { t } = useTranslation();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [config, setConfig] = useState<StaticIpConfig>({ isEnabled: false, content: '' });
  const [isEnabled, setIsEnabled] = useState(false);
  const [content, setContent] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<Result | null>(null);

  useEffect(() => {
    getStaticIp();
  }, []);

  function openModal() {
    setIsEnabled(config.isEnabled);
    setContent(config.content);
    setResult(null);
    setIsModalOpen(true);
  }

  function getStaticIp() {
    api.getStaticIp().then((rsp) => {
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      setConfig({
        isEnabled: rsp.data.enabled,
        content: rsp.data.ip
      });
      setIsEnabled(rsp.data.enabled);
      setContent(rsp.data.ip);
    });
  }

  async function updateStaticIP() {
    if (config.isEnabled === isEnabled && config.content === content) {
      setIsModalOpen(false);
      return;
    }

    if (!validContent()) {
      return;
    }

    if (isLoading) return;
    setIsLoading(true);

    const ip = isEnabled ? content : '';

    try {
      const rsp = await api.setStaticIP(isEnabled, ip);
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        setIsLoading(false);
        setResult({ type: 'failed', message: t('settings.device.staticIp.failed') });
        return;
      }

      setTimeout(() => {
        setConfig({
          isEnabled: isEnabled,
          content: ip
        });
        setContent(ip);
        setIsLoading(false);
        setResult({ type: 'success', message: t('settings.device.staticIp.success') });
      }, 10 * 1000);
    } catch (err) {
      console.log(err);
      setIsEnabled(false);
      return;
    }
  }

  function validContent() {
    if (content.trim() === '') {
      setResult({ type: 'failed', message: t('settings.device.staticIp.notEmpty') });
      return false;
    }

    const lines = content.split(/\r?\n/);

    for (const line of lines) {
      const trimmedLine = line.trim();

      if (trimmedLine === '' || trimmedLine.startsWith('#')) {
        continue;
      }

      const parts = trimmedLine.split(' ');
      const ipConfig = parts[0];

      const ipCidrRegex = /^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d{1,2}$/;
      if (!ipCidrRegex.test(ipConfig)) {
        setResult({ type: 'failed', message: t('settings.device.staticIp.notValid') });
        return false;
      }
    }

    return true;
  }

  const footer = [
    <Button key="cancel" onClick={() => setIsModalOpen(false)}>
      {t('settings.device.staticIp.cancelBtn')}
    </Button>,
    <Button key="ok" type="primary" loading={isLoading} onClick={updateStaticIP}>
      {t('settings.device.staticIp.okBtn')}
    </Button>
  ];

  return (
    <>
      <div className="flex items-center justify-between pt-5">
        <span>{t('settings.device.staticIp.title')}</span>

        <Button type={config.isEnabled ? 'primary' : 'default'} onClick={openModal}>
          {t('settings.device.staticIp.config')}
        </Button>
      </div>

      <Modal
        open={isModalOpen}
        closable={false}
        centered={true}
        footer={footer}
        onCancel={() => setIsModalOpen(false)}
      >
        <div className="flex flex-col space-y-3">
          <div className="flex flex-col rounded-lg bg-neutral-800 p-3">
            <div className="flex items-center">
              <span className="text-lg">{t('settings.device.staticIp.title')}</span>
            </div>
            <span className="text-sm text-neutral-500">{t('settings.device.staticIp.tip')}</span>
            <a
              className="w-fit text-sm"
              href="https://wiki.sipeed.com/hardware/en/kvm/NanoKVM_Pro/extended.html#How-to-Set-a-Static-IP"
              target="_blank"
            >
              {t('settings.device.staticIp.document')}
            </a>
          </div>

          <div className="flex justify-between rounded-lg bg-neutral-800 p-3">
            <span>{t('settings.device.staticIp.enable')}</span>
            <Switch value={isEnabled} onChange={setIsEnabled} />
          </div>

          {isEnabled && (
            <TextArea
              autoFocus
              value={content}
              placeholder={t('settings.device.staticIp.placeholder')}
              autoSize={{ minRows: 5, maxRows: 10 }}
              onFocus={() => setResult(null)}
              onChange={(e) => setContent(e.target.value)}
            />
          )}

          <div className="h-7">
            {result && (
              <span
                className={clsx(
                  'text-sm',
                  result.type === 'success' ? 'text-green-500' : 'text-red-500'
                )}
              >
                {result.message}
              </span>
            )}
          </div>
        </div>
      </Modal>
    </>
  );
};
