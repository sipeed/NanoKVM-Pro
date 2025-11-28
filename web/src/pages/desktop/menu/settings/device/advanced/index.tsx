import { Collapse } from 'antd';
import { useTranslation } from 'react-i18next';

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
          children: <StaticIp />
        }
      ]}
    />
  );
};
