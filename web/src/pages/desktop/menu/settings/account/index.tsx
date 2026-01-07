import { useEffect, useState } from 'react';
import { Button, Divider } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

import * as api from '@/api/auth.ts';

import { Kvmadmin } from './kvmadmin.tsx';
import { Logout } from './logout.tsx';

export const Account = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [username, setUsername] = useState('');

  useEffect(() => {
    api.getAccount().then((rsp) => {
      if (rsp.code === 0) {
        setUsername(rsp.data.username);
      }
    });
  }, []);

  function changePassword() {
    navigate('/auth/password');
  }

  return (
    <>
      <div className="text-base font-bold">{t('settings.account.title')}</div>
      <Divider className="opacity-50" />

      <div className="flex flex-col space-y-5">
        <div className="flex items-center justify-between">
          <span>{t('settings.account.webAccount')}</span>
          <Button type="text">{username ? username : '-'}</Button>
        </div>

        <div className="flex items-center justify-between">
          <span>{t('settings.account.password')}</span>
          <Button type="primary" onClick={changePassword}>
            {t('settings.account.updateBtn')}
          </Button>
        </div>
      </div>

      <Divider className="opacity-50" style={{ margin: '40px 0' }} />

      <Kvmadmin />

      <Divider className="opacity-50" style={{ margin: '40px 0' }} />

      <Logout />
    </>
  );
};
