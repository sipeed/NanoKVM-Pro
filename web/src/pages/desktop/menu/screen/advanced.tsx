import { useSetAtom } from 'jotai';
import { SlidersVerticalIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { isKeyboardEnableAtom } from '@/jotai/keyboard.ts';
import { isSettingsOpenAtom, settingTabAtom } from '@/jotai/settings.ts';

export const Advanced = () => {
  const { t } = useTranslation();

  const setIsKeyboardEnable = useSetAtom(isKeyboardEnableAtom);
  const setIsSettingsOpen = useSetAtom(isSettingsOpenAtom);
  const setSettingTab = useSetAtom(settingTabAtom);

  function openSettings() {
    setIsSettingsOpen(true);
    setSettingTab('screen');
    setIsKeyboardEnable(false);
  }

  return (
    <div
      className="flex h-[30px] cursor-pointer items-center space-x-2 rounded pl-3 pr-6 text-neutral-300 hover:bg-neutral-700/70"
      onClick={openSettings}
    >
      <SlidersVerticalIcon size={18} />
      <span className="select-none text-sm">{t('screen.advanced')}</span>
    </div>
  );
};
