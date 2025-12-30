import { Collapse } from 'antd';
import { useTranslation } from 'react-i18next';

import { LowPower } from './low-power.tsx';
import { StaticIp } from './static-ip.tsx';

export const Advanced = () => {
  const { t } = useTranslation();

  return (
    <Collapse
      ghost
      expandIconPosition="end"
      items={[
        {
          key: 'advanced',
          label: t('settings.device.advanced'),
          children: (
            <div className="mt-5 flex flex-col space-y-8">
              <StaticIp />
              <LowPower />
            </div>
          )
        }
      ]}
    />
  );
};
