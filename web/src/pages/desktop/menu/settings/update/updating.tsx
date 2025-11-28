import { useEffect, useState } from 'react';
import { Spin, Steps } from 'antd';
import { useTranslation } from 'react-i18next';
import { IMessageEvent } from 'websocket';

import { client } from '@/lib/websocket.ts';

export const Updating = () => {
  const { t } = useTranslation();

  const [step, setStep] = useState(0);
  const [progress, setProgress] = useState(0);
  const [status, setStatus] = useState<'process' | 'error' | 'finish'>('process');

  useEffect(() => {
    client.register('update', handleMessage);

    return () => {
      client.unregister('update');
    };
  }, []);

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
    </div>
  );
};
