import { ReactNode, useEffect } from 'react';
import { message } from 'antd';
import { useAtomValue } from 'jotai';
import { MonitorPauseIcon, MonitorXIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { VideoStatus } from '@/types';
import { videoStatusAtom } from '@/jotai/screen.ts';

export const Message = () => {
  const { t } = useTranslation();
  const videoStatus = useAtomValue(videoStatusAtom);

  const [messageApi, contextHolder] = message.useMessage();

  const videoMessageKey = 'video';

  useEffect(() => {
    handleVideoStatus(videoStatus);
  }, [videoStatus]);

  function handleVideoStatus(status: VideoStatus) {
    if (status === VideoStatus.Normal) {
      messageApi.destroy(videoMessageKey);
      return;
    }

    let icon: ReactNode;
    let content: ReactNode;

    if (status === VideoStatus.NoImage) {
      icon = <MonitorXIcon size={16} className="text-red-500/80" />;
      content = <div className="pl-2 text-sm">{t('screen.noSignal')}</div>;
    } else if (status === VideoStatus.InconsistentVideoMode) {
      icon = <MonitorPauseIcon size={16} className="text-yellow-500/80" />;
      content = <div className="pl-2 text-sm">{t('screen.inconsistentVideoMode')}</div>;
    } else {
      return;
    }

    showMessage(videoMessageKey, icon, content);
  }

  function showMessage(key: string, icon: ReactNode, content: ReactNode) {
    messageApi.open({
      key: key,
      type: 'warning',
      icon: icon,
      content: content,
      duration: 0,
      className: 'mt-16'
    });
  }

  return <>{contextHolder}</>;
};
