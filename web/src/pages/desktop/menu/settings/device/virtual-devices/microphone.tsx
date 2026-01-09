import { useState } from 'react';
import { Switch, Tooltip } from 'antd';
import { useAtom } from 'jotai';
import { CircleAlertIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/virtual-device.ts';
import { microphoneEnabledAtom } from '@/jotai/audio.ts';

interface NetworkProps {
  isDisabled: boolean;
  isMicEnabled: boolean;
  getVirtualDevices: () => void;
}

export const Microphone = ({ isDisabled, isMicEnabled, getVirtualDevices }: NetworkProps) => {
  const { t } = useTranslation();

  const [microphoneEnabled, setMicrophoneEnabled] = useAtom(microphoneEnabledAtom);

  const [isLoading, setIsLoading] = useState(false);

  async function update() {
    if (isLoading) return;
    setIsLoading(true);

    if (microphoneEnabled) {
      setMicrophoneEnabled(false);
      await new Promise((resolve) => setTimeout(resolve, 1000));
    }

    try {
      const rsp = await api.updateVirtualDevice('mic');
      if (rsp.code !== 0) {
        return;
      }
      getVirtualDevices();
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-col space-y-1">
        <div className="flex items-center space-x-2">
          <span>{t('settings.device.mic')}</span>
          <Tooltip
            title={t('settings.device.micTip')}
            className="cursor-pointer text-neutral-500"
            placement="top"
            styles={{ root: { maxWidth: '450px', whiteSpace: 'pre-wrap' } }}
          >
            <CircleAlertIcon size={15} />
          </Tooltip>
        </div>

        <span className="text-xs text-neutral-500">{t('settings.device.micDesc')}</span>
      </div>

      <Switch disabled={isDisabled} checked={isMicEnabled} loading={isLoading} onChange={update} />
    </div>
  );
};
