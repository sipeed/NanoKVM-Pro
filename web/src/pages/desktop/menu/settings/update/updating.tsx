import { useEffect, useState } from 'react';
import { Button, Result, Spin, Steps } from 'antd';
import { useTranslation } from 'react-i18next';
import { IMessageEvent } from 'websocket';

import * as api from '@/api/application.ts';
import { client } from '@/lib/websocket.ts';

type UpdatingProps = {
  jobID?: string;
  onUpdateFailed: (message: string) => void;
  onUpdateSucceeded: () => void;
};

export const Updating = ({ jobID, onUpdateFailed, onUpdateSucceeded }: UpdatingProps) => {
  const { t } = useTranslation();
  const [step, setStep] = useState(0);
  const [progress, setProgress] = useState(0);
  const [status, setStatus] = useState<'process' | 'error' | 'finish'>('process');

  useEffect(() => {
    // Download progress still comes from the HTTP server. Once installation
    // is handed off, the independent worker's persisted job is authoritative.
    if (jobID) return;
    function handleMessage(message: IMessageEvent) {
      try {
        const data = JSON.parse(JSON.parse(message.data as string).data);
        setStep(data.step === 'download' ? 0 : data.step === 'install' ? 1 : 2);
        if (data.progress === -1) setStatus('error');
        else setProgress(data.progress);
      } catch (err) {
        console.log(err);
      }
    }
    client.on('update', handleMessage);
    return () => {
      client.off('update');
    };
  }, [jobID]);

  useEffect(() => {
    if (!jobID) return;
    const id = jobID;
    let active = true;
    let failures = 0;
    let timer: ReturnType<typeof setTimeout>;

    async function checkUpdate() {
      let terminal = false;
      try {
        const rsp = await api.getManualUpdate(id);
        if (!active) return;
        if (rsp.code !== 0 || !rsp.data) throw new Error('invalid update status');
        failures = 0;
        const data = rsp.data as api.ManualUpdateInspection;
        setStep(data.step === 'restart' ? 2 : 1);
        setProgress(data.progress ?? 0);
        if (data.state === 'failed') {
          terminal = true;
          setStatus('error');
          onUpdateFailed(data.error || t('settings.update.manual.installFailed'));
        } else if (data.state === 'succeeded' || data.state === 'reboot_scheduled') {
          terminal = true;
          setStep(2);
          setProgress(100);
          setStatus('finish');
          onUpdateSucceeded();
        }
      } catch {
        if (!active) return;
        // The application is expected to disappear while its DEB restarts it.
        failures += 1;
        if (failures >= 120) {
          terminal = true;
          onUpdateFailed(t('settings.update.updateFailed'));
        }
      }
      if (active && !terminal) timer = setTimeout(checkUpdate, 2000);
    }

    void checkUpdate();
    return () => {
      active = false;
      clearTimeout(timer);
    };
  }, [jobID, onUpdateFailed, onUpdateSucceeded, t]);

  return (
    <div className="flex flex-col items-center justify-center space-y-20 py-20">
      {status === 'finish' ? (
        <Result
          status="success"
          title={t('settings.update.manual.finished')}
          extra={
            <Button onClick={() => window.location.reload()}>
              {t('settings.update.manual.refresh')}
            </Button>
          }
        />
      ) : (
        <>
          <Spin size="large" />
          <Steps
            current={step}
            percent={progress}
            status={status}
            items={[
              { title: t('settings.update.download') },
              { title: t('settings.update.install') },
              { title: t('settings.update.restart') }
            ]}
          />
        </>
      )}
    </div>
  );
};
