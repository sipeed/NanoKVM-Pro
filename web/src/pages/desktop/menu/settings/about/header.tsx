import { useState } from 'react';
import { Button, Divider, Modal, Typography } from 'antd';
import { useTranslation } from 'react-i18next';

import { switchPiKVM } from '@/api/vm.ts';

const { Paragraph, Text } = Typography;

export const Header = () => {
  const { t } = useTranslation();

  const [isLoading, setIsLoading] = useState(false);
  const [open, setOpen] = useState(false);
  const [errMsg, setErrMsg] = useState('');

  function switchSystem() {
    if (isLoading) return;
    setIsLoading(true);

    const timeoutId = setTimeout(() => {
      window.location.reload();
    }, 60000);

    switchPiKVM()
      .then((rsp) => {
        if (rsp.code !== 0) {
          setErrMsg(rsp.msg);
          setIsLoading(false);
          clearTimeout(timeoutId);
        }
      })
      .catch((err) => {
        console.log(err);
        setIsLoading(false);
        clearTimeout(timeoutId);
      });
  }

  return (
    <>
      <div className="flex items-center justify-between">
        <div className="text-base font-bold">{t('settings.about.title')}</div>

        <Button danger onClick={() => setOpen(true)}>
          {t('settings.about.pikvm.title')}
        </Button>
      </div>

      <Modal
        open={open}
        title={t('settings.about.pikvm.attention')}
        onCancel={() => setOpen(false)}
        footer={null}
      >
        <Divider />

        <Paragraph>{t('settings.about.pikvm.desc1')}</Paragraph>

        <Paragraph>
          <Text strong> {t('settings.about.pikvm.desc2')}</Text>
        </Paragraph>

        {errMsg && <div className="pt-1 text-sm text-red-500">{errMsg}</div>}

        <div className="flex justify-center pt-5">
          <Button danger type="primary" loading={isLoading} onClick={switchSystem}>
            {t('settings.about.pikvm.confirm')}
          </Button>
        </div>
      </Modal>
    </>
  );
};
