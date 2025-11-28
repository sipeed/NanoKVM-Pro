import { useEffect, useRef, useState } from 'react';
import { LoaderCircleIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { getBaseUrl } from '@/lib/service.ts';

type TimeProps = {
  timezone: string;
  timeFormat: string;
};

export const Time = ({ timezone, timeFormat }: TimeProps) => {
  const { t, i18n } = useTranslation();

  const [time, setTime] = useState('');

  const formaterRef = useRef<Intl.DateTimeFormat | null>(null);

  useEffect(() => {
    if (!timezone) return;

    formaterRef.current = i18n.language.startsWith('zh')
      ? new Intl.DateTimeFormat('zh-CN', {
          timeZone: timezone,
          year: 'numeric',
          month: 'long',
          day: 'numeric',
          hour: 'numeric',
          minute: 'numeric',
          second: 'numeric',
          hour12: timeFormat === '12H'
        })
      : new Intl.DateTimeFormat('en-US', {
          timeZone: timezone,
          year: 'numeric',
          month: 'long',
          day: 'numeric',
          hour: 'numeric',
          minute: 'numeric',
          second: 'numeric',
          hour12: timeFormat === '12H'
        });
  }, [timezone, timeFormat, i18n.language]);

  useEffect(() => {
    const baseURL = getBaseUrl('http');
    const url = `${baseURL}/api/vm/time`;
    const eventSource = new EventSource(url);

    eventSource.onmessage = (event) => {
      if (!formaterRef.current) return;

      const date = new Date(Number(event.data));
      if (!date) {
        setTime('-');
        return;
      }

      setTime(formaterRef.current.format(date));
    };

    eventSource.onerror = (error) => {
      console.log(error);
      setTime('-');
      eventSource.close();
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return (
    <div className="flex items-center justify-between">
      <span>{t('settings.device.datetime.datetime')}</span>

      {time === '' ? (
        <LoaderCircleIcon className="animate-spin text-neutral-500" size={18} />
      ) : (
        <div className="font-mono">{time}</div>
      )}
    </div>
  );
};
