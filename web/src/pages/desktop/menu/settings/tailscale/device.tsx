import { useEffect, useState } from 'react';
import { LogoutOutlined, LockOutlined, UnlockOutlined } from '@ant-design/icons';
import { Button, Divider, Popconfirm, Switch, Tooltip } from 'antd';
import { CircleAlertIcon, EthernetPortIcon, WifiIcon } from 'lucide-react';
import { Trans, useTranslation } from 'react-i18next';

import * as api from '@/api/extensions/tailscale.ts';
import * as vmApi from '@/api/vm.ts';
import { Tailscale } from '@/components/icons/tailscale.tsx';

import { Status } from './types.ts';

type DeviceProps = {
  status: Status;
  onLogout: () => void;
};

export const Device = ({ status, onLogout }: DeviceProps) => {
  const { t } = useTranslation();

  const [isRunning, setIsRunning] = useState(false);
  const [isServing, setIsServing] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [isServingUpdating, setIsServingUpdating] = useState(false);
  const [isLogging, setIsLogging] = useState(false);
  const [errMsg, setErrMsg] = useState('');
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [authConsentUrl, setAuthConsentUrl] = useState('');
  const [serveHttpsUrl, setServeHttpsUrl] = useState('');
  const [ips, setIps] = useState<{ addr: string; type: string }[]>([]);

  const isAwaitingAuth = authConsentUrl !== '';

  useEffect(() => {
    vmApi.getInfo().then((rsp: any) => {
      if (rsp.code === 0 && rsp.data?.ips) {
        setIps(rsp.data.ips);
      }
    });
  }, []);

  useEffect(() => {
    setIsRunning(status.state === 'running');
    setIsServing(status.serve);
    setServeHttpsUrl(status.serveUrl || '');
  }, [status]);

  function handleServeEnabled(newUrl: string) {
    setIsServing(true);
    if (newUrl) {
      setServeHttpsUrl(newUrl);
      if (window.location.hostname === status.ip) {
        window.location.href = newUrl;
      }
    }
  }

  function redirectToFallback() {
    if (window.location.hostname.endsWith('.ts.net')) {
      window.location.href = `http://${status.ip}`;
    }
  }

  async function update() {
    if (isUpdating) return;
    setIsUpdating(true);

    try {
      const rsp = isRunning ? await api.down() : await api.up();
      if (rsp.code !== 0) {
        setErrMsg(rsp.msg);
        return;
      }

      setIsRunning(!isRunning);
    } finally {
      setIsUpdating(false);
    }
  }

  async function updateServe() {
    if (isServingUpdating) return;
    setIsServingUpdating(true);
    setErrMsg('');

    try {
      // if currently serving, disable it
      if (isServing) {
        const rsp = await api.serve(false);
        if (rsp.code !== 0) {
          setErrMsg(rsp.msg);
          return;
        }
        setIsServing(false);

        // Auto-redirect to the Tailscale IP if we were connected via the Tailscale HTTPS URL
        redirectToFallback();

        return;
      }

      // if not serving, enable it
      const rsp = await api.serve(true);
      if (rsp.code !== 0) {
        setErrMsg(rsp.msg);
        return;
      }

      // check if auth url returned
      if (rsp.data && rsp.data.authUrl) {
        setAuthConsentUrl(rsp.data.authUrl);
        window.open(rsp.data.authUrl, '_blank');
        return;
      }

      handleServeEnabled(rsp.data?.serveUrl || '');
    } catch (err) {
      setErrMsg(String(err));
    } finally {
      setIsServingUpdating(false);
    }
  }

  async function onConfirmServe() {
    setIsServingUpdating(true);
    setErrMsg('');
    try {
      // Re-run the serve command now that it's authorized
      const rsp = await api.serve(true);
      if (rsp.code !== 0) {
        setErrMsg(rsp.msg);
        return;
      }

      // If it still wants authorization, the user probably didn't authorize it properly
      if (rsp.data && rsp.data.authUrl) {
        setAuthConsentUrl('');
        setErrMsg(t('settings.tailscale.enableHttpsCertsPlain'));
        return;
      }

      setAuthConsentUrl(''); // Success!
      handleServeEnabled(rsp.data?.serveUrl || '');
    } catch (err) {
      setErrMsg(String(err));
    } finally {
      setIsServingUpdating(false);
    }
  }

  async function logout() {
    if (isLogging) return;
    setIsLogging(true);

    api
      .logout()
      .then((rsp) => {
        if (rsp.code !== 0) {
          setErrMsg(rsp.msg);
          return;
        }

        onLogout();
      })
      .finally(() => {
        setIsLogging(false);
      });
  }

  const isHttpsToggleChecked = !isConfirmOpen && isServing;
  const isHttpsToggleLoading = isConfirmOpen || isServingUpdating || isAwaitingAuth;

  return (
    <div className="flex flex-col space-y-7">
      <div className="flex justify-between">
        <span>{t('settings.tailscale.enable')}</span>
        <Switch checked={isRunning} loading={isUpdating} onClick={update} />
      </div>

      <div className="flex flex-col space-y-2">
        <div className="flex justify-between">
          <div className="flex items-center space-x-2">
            <span>{t('settings.tailscale.enableHttps')}</span>
            <Tooltip
              placement="right"
              arrow={false}
              title={
                <div className="whitespace-pre-line">
                  {t('settings.tailscale.enableHttpsTooltip')}
                </div>
              }
            >
              <CircleAlertIcon size={15} className="cursor-pointer text-neutral-500" />
            </Tooltip>
          </div>
          {isServing ? (
            <Popconfirm
              placement="bottomRight"
              icon={<UnlockOutlined style={{ color: '#ef4444' }} />}
              title={t('settings.tailscale.disableHttpsConfirm')}
              description={
                <div className="mt-1">
                  {t('settings.tailscale.disableHttpsDesc')}
                  <br />
                  {t('settings.tailscale.disableHttpsFallback')}
                  <div className="mt-2 flex flex-col space-y-1">
                    {ips.map((ip) => (
                      <div key={ip.addr} className="flex items-center space-x-2">
                        <div className="size-[16px] text-neutral-500">
                          {ip.type === 'Wireless' ? <WifiIcon size={16} /> : <EthernetPortIcon size={16} />}
                        </div>
                        <a href={`http://${ip.addr}`} target="_blank" rel="noreferrer" className="text-blue-500 hover:underline">
                          http://{ip.addr}
                        </a>
                      </div>
                    ))}
                    <div className="flex items-center space-x-2">
                      <div className="size-[16px] flex items-center justify-center">
                        <Tailscale />
                      </div>
                      <a href={`http://${status.ip}`} target="_blank" rel="noreferrer" className="text-blue-500 hover:underline">
                        http://{status.ip}
                      </a>
                    </div>
                  </div>
                </div>
              }
              open={isConfirmOpen}
              onOpenChange={setIsConfirmOpen}
              onConfirm={updateServe}
              okText={t('settings.tailscale.okBtn')}
              cancelText={t('settings.tailscale.cancelBtn')}
            >
              <Switch checked={isHttpsToggleChecked} loading={isHttpsToggleLoading} disabled={!isRunning} />
            </Popconfirm>
          ) : (
            <Switch
              checked={isServing || isAwaitingAuth}
              loading={isServingUpdating || isAwaitingAuth}
              disabled={!isRunning}
              onClick={updateServe}
            />
          )}
        </div>
        {isAwaitingAuth && (
          <div className="mt-4 flex w-full flex-col items-center justify-center space-y-4 rounded-lg border border-neutral-500/30 p-4">
            <Button type="link" href={authConsentUrl} target="_blank">
              {authConsentUrl}
            </Button>
            <span className="text-xs text-neutral-600">
              <Trans
                i18nKey="settings.tailscale.enableHttpsCerts"
                components={{ 1: <strong /> }}
              />
            </span>
            <div className="flex justify-center w-full">
              <Button type="primary" shape="round" onClick={onConfirmServe}>
                {t('settings.tailscale.authSuccess')}
              </Button>
            </div>
          </div>
        )}
      </div>

      {isServing && serveHttpsUrl && !isAwaitingAuth && (
        <div className="flex justify-between">
          <span>{t('settings.tailscale.httpsUrl')}</span>
          <div className="flex items-center space-x-2">
            <LockOutlined className="text-green-500" />
            <a href={serveHttpsUrl} target="_blank" rel="noreferrer" className="text-blue-500 hover:underline">
              {serveHttpsUrl}
            </a>
          </div>
        </div>
      )}

      <div className="flex justify-between">
        <span>{t('settings.tailscale.deviceName')}</span>
        <span>{status.name}</span>
      </div>

      <div className="flex justify-between">
        <span>{t('settings.tailscale.deviceIP')}</span>
        <span>{status.ip}</span>
      </div>

      <div className="flex justify-between">
        <span>{t('settings.tailscale.account')}</span>
        <span>{status.account}</span>
      </div>

      <Divider className="opacity-50" style={{ margin: '50px 0 0 0' }} />

      <div className="flex justify-center pt-7">
        <Popconfirm
          placement="bottom"
          title={t('settings.tailscale.logoutDesc')}
          okText={t('settings.tailscale.okBtn')}
          cancelText={t('settings.tailscale.cancelBtn')}
          onConfirm={logout}
        >
          <Button
            danger
            type="primary"
            size="large"
            shape="round"
            icon={<LogoutOutlined />}
            loading={isLogging}
          >
            {t('settings.tailscale.logout')}
          </Button>
        </Popconfirm>
      </div>

      {errMsg && <span className="text-red-500">{errMsg}</span>}
    </div>
  );
};
