import { useEffect, useRef, useState } from 'react';
import { LockOutlined, WifiOutlined } from '@ant-design/icons';
import { Button, Divider, Input, Modal, Popover } from 'antd';
import {
  LinkIcon,
  LoaderCircleIcon,
  LoaderIcon,
  LockIcon,
  WifiHighIcon,
  WifiIcon,
  WifiLowIcon,
  WifiOffIcon,
  WifiPenIcon
} from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/network.ts';

type WiFi = {
  ssid: string;
  bssid: string;
  signal: number;
  frequency: string;
  security: string;
};

export const Wifi = () => {
  const { t } = useTranslation();

  const [isSupported, setIsSupported] = useState(false);
  const [isAPMode, setIsAPMode] = useState(false);
  const [connectedWifi, setConnectedWifi] = useState<WiFi | null>();

  const [wifiList, setWifiList] = useState<WiFi[]>([]);
  const [isScanning, setIsScanning] = useState(false);

  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedWifi, setSelectedWifi] = useState<WiFi | null>(null);

  const [ssid, setSsid] = useState('');
  const [password, setPassword] = useState('');
  const [isConnecting, setIsConnecting] = useState(false);
  const [connectingBssid, setConnectingBssid] = useState('');
  const [message, setMessage] = useState('');

  const [isDisconnectModalOpen, setIsDisconnectModalOpen] = useState(false);
  const [isDisconnecting, setIsDisconnecting] = useState(false);

  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    getWiFi();
    scanWiFi();
  }, []);

  useEffect(() => {
    if (isPopoverOpen && !intervalRef.current) {
      intervalRef.current = setInterval(scanWiFi, 10 * 1000);
    }

    if (!isPopoverOpen && intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  }, [isPopoverOpen]);

  function getWiFi() {
    api.getWiFi().then((rsp) => {
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      setIsSupported(!!rsp.data?.supported);
      setIsAPMode(!!rsp.data?.apMode);

      if (rsp.data?.connected && rsp.data?.wifi) {
        setConnectedWifi(rsp.data.wifi);
      } else {
        setConnectedWifi(null);
      }
    });
  }

  function scanWiFi() {
    if (isScanning || isConnecting) return;
    setIsScanning(true);

    api
      .scanWiFi()
      .then((rsp) => {
        if (rsp.code !== 0) {
          setIsAPMode(rsp.code === -1);
          return;
        }

        if (rsp.data?.wifiList) {
          const list = sortWifiList(rsp.data.wifiList);
          setWifiList(list);
        }
      })
      .finally(() => {
        setIsScanning(false);
      });
  }

  function connectWifi(name?: string, pwd?: string) {
    setMessage('');

    const wifiSSID = name ? name : ssid;
    const wifiPwd = name ? pwd || '' : password;

    if (isConnecting || !wifiSSID) return;
    setIsConnecting(true);

    api
      .connectWifi(wifiSSID, wifiPwd)
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
          setMessage(t('wifi.failed'));
          return;
        }

        getWiFi();
        scanWiFi();
        setIsModalOpen(false);
      })
      .finally(() => {
        setIsConnecting(false);
        setConnectingBssid('');
      });
  }

  function disconnectWifi() {
    if (isDisconnecting) return;
    setIsDisconnecting(true);

    api
      .disconnectWifi()
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
        }

        getWiFi();
        scanWiFi();
        setIsDisconnectModalOpen(false);
      })
      .finally(() => {
        setIsDisconnecting(false);
      });
  }

  function sortWifiList(list: WiFi[]) {
    const arr = Array.isArray(list) ? [...list] : [];

    arr.sort((a, b) => {
      const sa = (a?.ssid || '').toLowerCase();
      const sb = (b?.ssid || '').toLowerCase();
      if (sa < sb) return -1;
      if (sa > sb) return 1;
      return 0;
    });

    return arr;
  }

  function getWifiSignalIcon(signal: number) {
    if (signal >= -60) {
      return <WifiIcon size={14} />;
    } else if (signal >= -70 && signal < -60) {
      return <WifiHighIcon size={14} />;
    } else {
      return <WifiLowIcon size={14} />;
    }
  }

  function handleOpenChange(open: boolean) {
    if (open) {
      scanWiFi();
      setIsPopoverOpen(true);
      return;
    }

    if (isModalOpen || isDisconnectModalOpen) {
      return;
    }

    setIsPopoverOpen(false);
  }

  function handleSelect(wifi: WiFi | null) {
    setMessage('');

    if (isConnecting) return;
    if (wifi && wifi.bssid === connectedWifi?.bssid) return;

    setSsid(wifi ? wifi.ssid : '');
    setPassword('');
    setConnectingBssid(wifi ? wifi.bssid : '');

    if (wifi && wifi.security.toLowerCase() === 'open') {
      connectWifi(wifi.ssid, '');
      return;
    }

    setSelectedWifi(wifi);
    setIsModalOpen(true);
  }

  // Wi-Fi list
  const content = (
    <div className="flex max-h-[400px] w-[280px] flex-col space-y-1 overflow-y-auto p-3">
      {/* connected Wi-Fi */}
      {connectedWifi && (
        <>
          <div
            className="flex cursor-pointer items-center justify-between rounded-lg px-2 py-1.5 hover:bg-neutral-500/20"
            onClick={() => setIsDisconnectModalOpen(true)}
          >
            <div className="w-2/3 truncate">{connectedWifi.ssid}</div>
            <div className="flex w-1/3 items-center justify-end space-x-2 text-blue-500">
              {connectedWifi.security.toLowerCase() !== 'open' && <LockIcon size={14} />}
              {getWifiSignalIcon(connectedWifi.signal)}
            </div>
          </div>

          <Divider />
        </>
      )}

      {wifiList.map((wifi) => (
        <div
          key={wifi.bssid}
          className="flex cursor-pointer items-center justify-between rounded-lg px-2 py-1.5 hover:bg-neutral-500/20"
          onClick={() => handleSelect(wifi)}
        >
          <div className="w-2/3 truncate">{wifi.ssid}</div>
          <div className="flex w-1/3 items-center justify-end space-x-2">
            {wifi.security.toLowerCase() !== 'open' && <LockIcon size={14} />}
            {isConnecting && connectingBssid === wifi.bssid ? (
              <LoaderCircleIcon className="animate-spin text-green-500" size={14} />
            ) : (
              getWifiSignalIcon(wifi.signal)
            )}
          </div>
        </div>
      ))}

      {isScanning && (
        <div className="flex justify-center py-1.5 text-neutral-500">
          <LoaderIcon className="animate-spin" size={16} />
        </div>
      )}

      <div
        className="flex cursor-pointer items-center rounded-lg px-2 py-1.5 hover:bg-neutral-500/20"
        onClick={() => handleSelect(null)}
      >
        {t('settings.device.wifi.others')}
      </div>
    </div>
  );

  if (!isSupported) {
    return <></>;
  }

  return (
    <>
      {isAPMode ? (
        <div className="flex items-center justify-between">
          <div className="flex flex-col">
            <span>{t('settings.device.wifi.title')}</span>
            <span className="text-xs text-neutral-500">{t('settings.device.wifi.apMode')}</span>
          </div>

          <Button shape="round" onClick={() => handleSelect(null)}>
            <LinkIcon size={14} />
          </Button>
        </div>
      ) : (
        <div className="flex items-center justify-between">
          <div className="flex flex-col space-y-1">
            <span>{t('settings.device.wifi.title')}</span>
            <span className="text-xs text-neutral-500">
              {t('settings.device.wifi.description')}
            </span>
          </div>

          <Popover
            className="p-3"
            styles={{ body: { padding: 0 } }}
            placement="bottomRight"
            trigger="click"
            content={content}
            arrow={false}
            open={isPopoverOpen}
            onOpenChange={handleOpenChange}
          >
            {connectedWifi ? (
              <Button type="primary" shape="round" size="small">
                <div className="flex items-center justify-center px-1.5">
                  <WifiIcon size={16} />
                </div>
              </Button>
            ) : (
              <Button shape="round" size="small">
                <div className="flex items-center justify-center px-1.5">
                  <WifiIcon size={16} />
                </div>
              </Button>
            )}
          </Popover>
        </div>
      )}

      <Modal
        closable={false}
        open={isModalOpen}
        centered={true}
        onOk={() => connectWifi()}
        onCancel={() => {
          if (!isConnecting) {
            setIsModalOpen(false);
          }
        }}
        okText={t('settings.device.wifi.joinBtn')}
        cancelText={t('settings.device.wifi.cancelBtn')}
        confirmLoading={isConnecting}
      >
        {/* title */}
        <div className="flex items-center space-x-5">
          <div className="h-[64px] w-[64px]">
            <WifiPenIcon size={64} className="text-blue-400" />
          </div>

          {!selectedWifi ? (
            <div className="flex flex-col">
              <span className="text-lg font-bold">{t('settings.device.wifi.connect')}</span>
              <span className="text-xs text-neutral-400">
                {t('settings.device.wifi.connectDesc1')}
              </span>
            </div>
          ) : (
            <div className="flex flex-col space-y-1">
              <span className="break-all text-lg font-bold">{selectedWifi.ssid}</span>
              <span className="text-xs text-neutral-400">
                {t('settings.device.wifi.connectDesc2')}
              </span>
            </div>
          )}
        </div>

        {/* form */}
        <div className="flex flex-col items-center space-y-3 py-6">
          {!selectedWifi && (
            <Input
              value={ssid}
              style={{ width: '300px' }}
              prefix={<WifiOutlined />}
              placeholder={t('settings.device.wifi.ssid')}
              onChange={(e) => setSsid(e.target.value)}
            />
          )}
          <Input.Password
            value={password}
            style={{ width: '300px' }}
            prefix={<LockOutlined />}
            placeholder={t('settings.device.wifi.password')}
            onChange={(e) => setPassword(e.target.value)}
          />

          {!!message && <span className="text-sm text-red-500">{message}</span>}
        </div>
      </Modal>

      <Modal
        closable={false}
        open={isDisconnectModalOpen}
        centered={true}
        okText={t('settings.device.wifi.confirmBtn')}
        cancelText={t('settings.device.wifi.cancelBtn')}
        onOk={disconnectWifi}
        onCancel={() => {
          if (!isDisconnecting) {
            setIsDisconnectModalOpen(false);
          }
        }}
        okType="danger"
        confirmLoading={isDisconnecting}
      >
        <div className="flex items-center space-x-5">
          <div className="h-[64px] w-[64px]">
            <WifiOffIcon size={64} className="text-blue-400" />
          </div>

          <div className="flex flex-col space-y-1">
            <span className="break-all text-lg font-bold">{connectedWifi?.ssid}</span>
            <span className="text-xs text-neutral-400">{t('settings.device.wifi.disconnect')}</span>
          </div>
        </div>
      </Modal>
    </>
  );
};
