import { useEffect, useState } from 'react';
import { Switch, Tooltip } from 'antd';
import { CircleAlertIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/vm.ts';

export const HdmiPassthrough = () => {
  const { t } = useTranslation();

  const [isEnabled, setIsEnabled] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    getPassthrough();
  }, []);

  function getPassthrough() {
    setIsLoading(true);

    api
      .getHdmiPassthrough()
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
          return;
        }

        setIsEnabled(!!rsp.data?.enabled);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  function update() {
    if (isLoading) return;
    setIsLoading(true);

    api
      .setHdmiPassthrough(!isEnabled)
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
          return;
        }
        setIsEnabled(!isEnabled);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-col">
        <div className="flex items-center space-x-2">
          <span>{t('settings.device.passthrough.title')}</span>
          <Tooltip
            title={t('settings.device.passthrough.tip')}
            className="cursor-pointer text-neutral-500"
            placement="bottom"
            styles={{ root: { maxWidth: '450px' } }}
          >
            <CircleAlertIcon size={15} />
          </Tooltip>
        </div>

        <span className="text-xs text-neutral-500">
          {t('settings.device.passthrough.description')}
        </span>
      </div>

      <Switch checked={isEnabled} loading={isLoading} onChange={update} />
    </div>
  );
};
