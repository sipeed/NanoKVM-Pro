import { useEffect, useState } from 'react';
import { Button, Popconfirm, Tooltip } from 'antd';
import { CircleAlertIcon, LoaderCircleIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/extensions/kvmadmin.ts';

type Status = '' | 'notInstall' | 'notRunning' | 'running';
type Failure = '' | 'install' | 'run';

export const Kvmadmin = () => {
  const { t } = useTranslation();

  const [status, setStatus] = useState<Status>('');
  const [failure, setFailure] = useState<Failure>('');
  const [isLoading, setIsLoading] = useState(false);
  const [isPopOpen, setIsPopOpen] = useState(false);

  useEffect(() => {
    getStatus();
  }, []);

  async function getStatus() {
    try {
      const rsp = await api.getStatus();
      if (rsp.code === 0) {
        setStatus(rsp.data.state);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function install() {
    if (isLoading) return;
    setIsLoading(true);

    try {
      setFailure('');

      const rsp = await api.install();

      if (rsp.code !== 0) {
        console.log(rsp.msg);
        setFailure('install');
      } else {
        await getStatus();
      }
    } catch (err) {
      console.log(err);
    }

    setIsLoading(false);
  }

  async function start() {
    if (isLoading) return;
    setIsLoading(true);

    try {
      setFailure('');

      const rsp = await api.start();

      if (rsp.code !== 0) {
        console.log(rsp.msg);
        setFailure('run');
      } else {
        await getStatus();
      }
    } catch (err) {
      console.log(err);
    }

    setIsLoading(false);
  }

  function visit() {
    const protocol = window.location.protocol;
    const hostname = window.location.hostname;
    const url = `${protocol}//${hostname}:8999`;

    window.open(url, '_blank');
  }

  async function uninstall() {
    if (isLoading) return;
    setIsLoading(true);

    setIsPopOpen(false);

    try {
      const rsp = await api.uninstall();
      if (rsp.code !== 0) {
        console.log(rsp.msg);
      }

      await getStatus();
    } catch (err) {
      console.log(err);
    }

    setIsLoading(false);
  }

  const Uninstall = (
    <Popconfirm
      open={isPopOpen}
      placement="bottom"
      title={t('settings.kvmadmin.attention')}
      description={
        <span>
          {t('settings.kvmadmin.confirmUninstall')}
          <br />
          {t('settings.kvmadmin.clearData')}
        </span>
      }
      okType="danger"
      okText={t('settings.kvmadmin.okBtn')}
      cancelText={t('settings.kvmadmin.cancelBtn')}
      onConfirm={uninstall}
      onCancel={() => setIsPopOpen(false)}
    >
      <Button danger loading={isLoading} onClick={() => setIsPopOpen(true)}>
        {t('settings.kvmadmin.uninstall')}
      </Button>
    </Popconfirm>
  );

  return (
    <div className="flex flex-col">
      <div className="flex items-center justify-between">
        <div className="flex flex-col">
          <div className="flex items-center space-x-2">
            <span>{t('settings.kvmadmin.title')}</span>

            <Tooltip
              title={t('settings.kvmadmin.tip')}
              className="cursor-pointer text-neutral-500"
              placement="bottom"
              overlayStyle={{ maxWidth: '350px' }}
            >
              <CircleAlertIcon size={15} />
            </Tooltip>
          </div>

          <span className="text-xs text-neutral-500">{t('settings.kvmadmin.description')}</span>
        </div>

        {/* loading */}
        {status === '' && <LoaderCircleIcon className="animate-spin text-neutral-500" size={18} />}

        {/* not install */}
        {status === 'notInstall' && (
          <Button type="primary" loading={isLoading} onClick={install}>
            {t('settings.kvmadmin.install')}
          </Button>
        )}

        {/* not running */}
        {status === 'notRunning' && (
          <div className="flex items-center space-x-2">
            <Button type="primary" loading={isLoading} onClick={start}>
              {t('settings.kvmadmin.start')}
            </Button>
            {Uninstall}
          </div>
        )}

        {/* running */}
        {status === 'running' && (
          <div className="flex items-center space-x-2">
            <Button type="primary" onClick={visit}>
              {t('settings.kvmadmin.visit')}
            </Button>
            {Uninstall}
          </div>
        )}
      </div>

      {/* error message */}
      {failure !== '' && (
        <div className="flex justify-end pt-2 text-sm text-red-500">
          {failure === 'install' && <span> {t('settings.kvmadmin.installFailed')}</span>}
          {failure === 'run' && <span> {t('settings.kvmadmin.startFailed')}</span>}
        </div>
      )}
    </div>
  );
};
