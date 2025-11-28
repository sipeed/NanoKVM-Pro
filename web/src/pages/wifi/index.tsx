import { useState } from 'react';
import { LockOutlined, WifiOutlined } from '@ant-design/icons';
import { Button, Form, Input, Steps } from 'antd';
import { SettingsIcon, ToggleRightIcon, WifiIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/network.ts';
import { Head } from '@/components/head.tsx';

type State = '' | 'loading' | 'success' | 'failed' | 'denied';

export const Wifi = () => {
  const { t } = useTranslation();

  const [state, setState] = useState<State>('');

  async function connect(values: any) {
    if (state === 'loading' || state === 'success') return;

    const ssid = values.ssid;
    const password = values.password;
    if (!ssid) {
      setState('');
      return;
    }

    setState('loading');

    try {
      const rsp = await api.connectWifiNoAuth(ssid, password);
      if (rsp && rsp.code !== 0) {
        if (rsp.code === -1) {
          setState('denied');
        } else if (rsp.code === -2 || rsp.code === -3) {
          setState('failed');
        }
        return;
      }
    } catch (err) {
      console.log(err);
    }

    setState('success');
  }

  return (
    <>
      <Head title={t('head.wifi')} />

      <div className="flex h-screen w-screen flex-col items-center justify-center">
        <Form
          style={{ minWidth: 300, maxWidth: 500 }}
          initialValues={{ remember: true }}
          onFinish={connect}
        >
          <div className="flex flex-col space-y-1 pb-5">
            <span className="text-center text-2xl font-semibold text-neutral-100">
              {t('wifi.title')}
            </span>
            <span className="text-center text-neutral-400">{t('wifi.description')}</span>
          </div>

          <Form.Item name="ssid" rules={[{ required: true }, { type: 'string', min: 1 }]}>
            <Input prefix={<WifiOutlined />} placeholder="SSID" />
          </Form.Item>

          <Form.Item name="password" rules={[{ required: false }, { type: 'string', max: 128 }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="Password" />
          </Form.Item>

          <Form.Item>
            <Button
              className="w-full"
              htmlType="submit"
              type="primary"
              loading={state === 'loading'}
            >
              {t('wifi.confirmBtn')}
            </Button>
          </Form.Item>
        </Form>

        <div className="flex max-w-[520px] justify-center px-5 pt-3 md:px-10">
          {state === 'success' && (
            <span className="text-sm text-green-500">{t('wifi.success')}</span>
          )}

          {state === 'failed' && <span className="text-sm text-red-500">{t('wifi.failed')} </span>}
          {state === 'denied' && (
            <div className="flex flex-col items-center space-y-5">
              <span className="text-sm text-red-500">{t('wifi.invalidMode')} </span>
              <Steps
                className="max-w-[400px]"
                size="small"
                responsive={false}
                items={[
                  { title: 'Settings', status: 'wait', icon: <SettingsIcon size={24} /> },
                  { title: 'Wi-Fi', status: 'wait', icon: <WifiIcon size={24} /> },
                  { title: 'Enable', status: 'wait', icon: <ToggleRightIcon size={24} /> }
                ]}
              />
            </div>
          )}
        </div>
      </div>
    </>
  );
};
