import { useEffect, useState } from 'react';
import { Switch } from 'antd';
import { useAtom } from 'jotai';
import {
  BotMessageSquare,
  DiscIcon,
  FileJsonIcon,
  MaximizeIcon,
  MicIcon,
  NetworkIcon,
  PowerIcon,
  TerminalSquareIcon,
  Volume2Icon,
  XIcon
} from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/vm.ts';
import { menuDisabledItemsAtom } from '@/jotai/settings.ts';

export const MenuBar = () => {
  const { t } = useTranslation();

  const [menuDisabledItems, setMenuDisabledItems] = useAtom(menuDisabledItemsAtom);

  const [loadingItem, setLoadingItem] = useState('');

  const items = [
    { key: 'volume', icon: <Volume2Icon size={16} /> },
    { key: 'microphone', icon: <MicIcon size={16} /> },
    { key: 'image', icon: <DiscIcon size={16} /> },
    { key: 'script', icon: <FileJsonIcon size={16} /> },
    { key: 'assistant', icon: <BotMessageSquare size={16} /> },
    { key: 'terminal', icon: <TerminalSquareIcon size={16} /> },
    { key: 'wol', icon: <NetworkIcon size={16} /> },
    { key: 'power', icon: <PowerIcon size={16} /> },
    { key: 'fullscreen', icon: <MaximizeIcon size={16} />, label: 'fullscreen.toggle' },
    { key: 'collapse', icon: <XIcon size={16} />, label: 'menu.collapse' }
  ];

  useEffect(() => {
    getMenuBarConfig();
  }, []);

  async function getMenuBarConfig() {
    try {
      const rsp = await api.getMenuBarConfig();
      if (rsp.code === 0 && rsp.data?.disabledItems?.length) {
        setMenuDisabledItems(rsp.data.disabledItems);
      }
    } catch (err) {
      console.log(err);
    }
  }

  async function updateItems(item: string) {
    if (loadingItem) return;
    setLoadingItem(item);

    const exist = menuDisabledItems.includes(item);
    const newItems = exist
      ? menuDisabledItems.filter((disabledItem) => disabledItem !== item)
      : [...menuDisabledItems, item];

    try {
      const rsp = await api.setMenuBarConfig(newItems);
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      setMenuDisabledItems(newItems);
    } catch (err) {
      console.log(err);
    } finally {
      setLoadingItem('');
    }
  }

  return (
    <>
      <div className="flex flex-col">
        <span className="text-neutral-400">{t('settings.appearance.menuBar')}</span>
        <span className="text-xs text-neutral-500">{t('settings.appearance.menuBarDesc')}</span>
      </div>

      <div className="mt-6 flex flex-col space-y-5">
        {items.map((item) => (
          <div key={item.key} className="flex items-center justify-between">
            <div className="flex items-center space-x-2 text-neutral-400">
              {item.icon}
              <span className="text-neutral-300">
                {item.label ? t(item.label) : t(`${item.key}.title`)}
              </span>
            </div>

            <Switch
              value={!menuDisabledItems.includes(item.key)}
              loading={item.key === loadingItem}
              onChange={() => updateItems(item.key)}
            />
          </div>
        ))}
      </div>
    </>
  );
};
