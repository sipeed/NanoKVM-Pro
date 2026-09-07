import { useEffect, useState } from 'react';
import { Button, Result, Spin, Steps } from 'antd';
import { useTranslation } from 'react-i18next';
import { IMessageEvent } from 'websocket';

import { client } from '@/lib/websocket.ts';
import * as api from '@/api/application.ts';

type UpdatingProps = {
  manualJobID?: string;
  onManualUpdateFailed?: (message: string) => void;
  onManualUpdateSucceeded?: () => void;
};

export const Updating = ({
  manualJobID,
  onManualUpdateFailed,
  onManualUpdateSucceeded
}: UpdatingProps) => {
  const { t } = useTranslation();

  const [step, setStep] = useState(0);
  const [progress, setProgress] = useState(0);
  const [status, setStatus] = useState<'process' | 'error' | 'finish'>('process');

  useEffect(() => {
    client.on('update', handleMessage);

    return () => {
      client.off('update');
    };
  }, []);

  useEffect(() => {
    // Manual jobs survive a service restart, so polling treats transient API
    // errors as retryable and only terminal states end the update screen.
    if (!manualJobID || !onManualUpdateFailed || !onManualUpdateSucceeded) return;

    const jobID = manualJobID;
    const failManualUpdate = onManualUpdateFailed;
    const finishManualUpdate = onManualUpdateSucceeded;
    let isActive = true;
    let terminal = false;

    function checkManualUpdate() {
      api
        .getManualUpdate(jobID)
        .then((rsp) => {
          if (!isActive || terminal || rsp.code !== 0 || !rsp.data) return;

          const data = rsp.data as {
            state?: api.ManualUpdateState;
            error?: string;
          };
          if (data.state === 'failed') {
            terminal = true;
            setStatus('error');
            failManualUpdate(data.error || t('settings.update.manual.installFailed'));
            return;
          }

          if (data.state === 'succeeded' || data.state === 'reboot_scheduled') {
            terminal = true;
            setStep(2);
            setProgress(100);
            setStatus('finish');
            finishManualUpdate();
          }
        })
        // A restart can temporarily make the API unavailable; wait for the next poll instead.
        .catch(() => undefined);
    }

    checkManualUpdate();
    const intervalID = window.setInterval(checkManualUpdate, 2000);

    return () => {
      isActive = false;
      window.clearInterval(intervalID);
    };
  }, [manualJobID, onManualUpdateFailed, onManualUpdateSucceeded, t]);

  function handleMessage(message: IMessageEvent) {
    try {
      const msg = JSON.parse(message.data as string);
      const data = JSON.parse(msg.data);

      handleStep(data.step);
      handleProgress(data.progress);
    } catch (err) {
      console.log(err);
    }
  }

  function handleStep(value: string) {
    switch (value) {
      case 'download':
        setStep(0);
        break;
      case 'install':
        setStep(1);
        break;
      case 'restart':
        restart();
        break;
      default:
        break;
    }
  }

  function handleProgress(value: number) {
    switch (value) {
      case -1:
        setStatus('error');
        break;
      default:
        setProgress(value);
        break;
    }
  }

  function restart() {
    setStep(2);

    setTimeout(() => {
      window.location.reload();
    }, 20 * 1000);

    let i = 0;

    const intervalId = setInterval(() => {
      setProgress(i * 5);

      if (i === 20) {
        clearInterval(intervalId);
        window.location.reload();
      }
      i++;
    }, 1000);
  }

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
              {
                title: t('settings.update.download')
              },
              {
                title: t('settings.update.install')
              },
              {
                title: t('settings.update.restart')
              }
            ]}
          />
        </>
      )}
    </div>
  );
};
