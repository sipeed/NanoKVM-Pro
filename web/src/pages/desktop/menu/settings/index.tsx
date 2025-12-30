import { useEffect, useState } from 'react';
import { Badge, Modal, Tooltip } from 'antd';
import clsx from 'clsx';
import { useAtom, useSetAtom } from 'jotai';
import {
  BadgeInfoIcon,
  CircleArrowUpIcon,
  MonitorIcon,
  PaletteIcon,
  SettingsIcon,
  SmartphoneIcon,
  UserRoundIcon
} from 'lucide-react';
import { useTranslation } from 'react-i18next';
import semver from 'semver';

import * as api from '@/api/application.ts';
import * as ls from '@/lib/localstorage.ts';
import { isKeyboardEnableAtom } from '@/jotai/keyboard.ts';
import { isSettingsOpenAtom, settingTabAtom } from '@/jotai/settings.ts';
import { Tailscale as TailscaleIcon } from '@/components/icons/tailscale';

import { About } from './about';
import { Account } from './account';
import { Appearance } from './appearance';
import { Device } from './device';
import { Screen } from './screen';
import { Tailscale } from './tailscale';
import { Update } from './update';

export const Settings = () => {
  const { t } = useTranslation();
  const setIsKeyboardEnable = useSetAtom(isKeyboardEnableAtom);
  const [isSettingsOpen, setIsSettingsOpen] = useAtom(isSettingsOpenAtom);
  const [settingTab, setSettingTab] = useAtom(settingTabAtom);

  const [isLocked, setIsLocked] = useState(false);

  const [isUpdateAvailable, setIsUpdateAvailable] = useState(false);

  const tabs = [
    { id: 'about', icon: <BadgeInfoIcon size={16} />, component: <About /> },
    { id: 'appearance', icon: <PaletteIcon size={16} />, component: <Appearance /> },
    { id: 'screen', icon: <MonitorIcon size={16} />, component: <Screen /> },
    { id: 'device', icon: <SmartphoneIcon size={16} />, component: <Device /> },
    {
      id: 'tailscale',
      icon: <TailscaleIcon />,
      component: <Tailscale setIsLocked={setIsLocked} />
    },
    {
      id: 'update',
      icon: <CircleArrowUpIcon size={16} />,
      component: <Update setIsLocked={setIsLocked} />
    },
    { id: 'account', icon: <UserRoundIcon size={16} />, component: <Account /> }
  ];

  useEffect(() => {
    const skip = ls.getSkipUpdate();
    if (!skip) {
      checkForUpdates();
    }
  }, []);

  function checkForUpdates() {
    api.getVersion().then((rsp: any) => {
      if (rsp.code !== 0) {
        return;
      }

      if (semver.gt(rsp.data.latest, rsp.data.current)) {
        setIsUpdateAvailable(true);
      }
    });
  }

  function changeTab(tab: string) {
    if (isLocked) {
      return;
    }

    setSettingTab(tab);

    if (isUpdateAvailable && tab === 'update') {
      setIsUpdateAvailable(false);
      ls.setSkipUpdate(true);
    }
  }

  function openModal() {
    setIsSettingsOpen(true);
    setIsKeyboardEnable(false);
  }

  function closeModal() {
    if (isLocked) {
      return;
    }

    setIsKeyboardEnable(true);
    setIsSettingsOpen(false);
    setSettingTab('about');
  }

  return (
    <>
      <Tooltip title={t('settings.title')} placement="bottom" mouseEnterDelay={0.6}>
        <div
          className="flex h-[30px] w-[30px] cursor-pointer items-center justify-center rounded text-white hover:bg-neutral-700/80"
          onClick={openModal}
        >
          <Badge dot={isUpdateAvailable} color="blue" offset={[0, 2]}>
            <div className="pt-[3px]">
              <SettingsIcon size={18} />
            </div>
          </Badge>
        </div>
      </Tooltip>

      <Modal
        open={isSettingsOpen}
        width={'80%'}
        centered={true}
        footer={null}
        destroyOnHidden={true}
        onCancel={closeModal}
        style={{ maxWidth: '1080px' }}
        styles={{ content: { padding: 0 } }}
      >
        <div className="flex h-[80vh] max-h-[700px] rounded-lg outline outline-1 outline-neutral-700">
          <div className="flex h-full max-w-[260px] flex-col space-y-0.5 rounded-l-lg bg-neutral-800 px-1 sm:w-1/5 md:w-1/4 md:px-2">
            <div className="hidden px-3 pt-10 text-xl sm:block">{t('settings.title')}</div>
            <div className="h-10 sm:h-5" />
            {tabs.map((tab) => (
              <div
                key={tab.id}
                className={clsx(
                  'flex cursor-pointer select-none items-center space-x-2 rounded-lg p-2 sm:px-3',
                  tab.id === settingTab ? 'bg-neutral-700/70' : 'hover:bg-neutral-700'
                )}
                onClick={() => changeTab(tab.id)}
              >
                <div className="h-[16px] w-[16px]">{tab.icon}</div>

                {isUpdateAvailable && tab.id === 'update' ? (
                  <Badge dot color="blue" offset={[6, 3]}>
                    <span className="hidden truncate text-sm sm:block">
                      {t(`settings.${tab.id}.title`)}
                    </span>
                  </Badge>
                ) : (
                  <span className="hidden truncate text-sm sm:block">
                    {t(`settings.${tab.id}.title`)}
                  </span>
                )}
              </div>
            ))}
          </div>

          <div className="flex h-full w-full flex-col items-center overflow-y-auto rounded-r-lg bg-neutral-900 px-3">
            <div className="w-full max-w-[600px] pb-10 pt-14">
              <>{tabs.find((tab) => tab.id === settingTab)?.component}</>
            </div>
          </div>
        </div>
      </Modal>
    </>
  );
};
