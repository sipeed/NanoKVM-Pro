import { useEffect } from 'react';
import { Button, notification } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

import { isPasswordUpdated } from '@/api/auth.ts';
import { getBrowserInfo } from '@/lib/browser.ts';
import * as storage from '@/lib/localstorage.ts';

export const Notification = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [api, contextHolder] = notification.useNotification();

  useEffect(() => {
    showPasswordNotification();
    showBrowserNotification();
  }, []);

  async function showPasswordNotification() {
    const skip = storage.getSkipModifyPassword();
    if (skip) return;

    try {
      const rsp = await isPasswordUpdated();
      if (rsp.code === 0 && !rsp.data.isUpdated) {
        api.warning({
          key: 'no_change_password',
          message: t('auth.changePassword'),
          description: t('auth.changePasswordDesc'),
          placement: 'topRight',
          btn: (
            <Button type="primary" onClick={changePassword}>
              {t('auth.ok')}
            </Button>
          ),
          duration: null,
          onClose: () => storage.setSkipModifyPassword(true)
        });
      }
    } catch (error) {
      console.error(error);
    }
  }

  function showBrowserNotification() {
    const skip = storage.getSkipBrowserNotification();
    if (skip) return;

    const info = getBrowserInfo();
    if (!info.isChrome) {
      api.info({
        key: 'recommend_browser.',
        message: t('notification.browser.title'),
        description: t('notification.browser.description'),
        placement: 'topRight',
        duration: null,
        onClose: () => storage.setSkipBrowserNotification(true)
      });
    }
  }

  function changePassword() {
    api.destroy();
    navigate('/auth/password');
  }

  return <>{contextHolder}</>;
};
