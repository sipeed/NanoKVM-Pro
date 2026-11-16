import { useState } from 'react';
import { getDefaultStore, useAtomValue } from 'jotai';
import { useTranslation } from 'react-i18next';

import { isMacroPlayingAtom } from '@/jotai/mouse.ts';

import { MACRO_ALREADY_PLAYING, playMacro } from './player.ts';
import { countScriptCommands } from './script.ts';
import type { Macro } from './types.ts';

type MacroItemProps = {
  macro: Macro;
  macros: Macro[];
};

export const MacroItem = ({ macro, macros }: MacroItemProps) => {
  const { t } = useTranslation();
  const isPlaying = useAtomValue(isMacroPlayingAtom);
  const [isLoading, setIsLoading] = useState(false);

  async function handleClick() {
    if (isLoading || isPlaying || !macro.script.trim()) return;
    if (getDefaultStore().get(isMacroPlayingAtom)) return;

    setIsLoading(true);

    try {
      const error = await playMacro(macro, macros);
      if (error && error !== MACRO_ALREADY_PLAYING) {
        console.log(error);
      }
    } catch (err) {
      console.log(err);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div
      className="flex h-[32px] w-full cursor-pointer items-center justify-between space-x-3 rounded px-3 hover:bg-neutral-700/30"
      onClick={handleClick}
    >
      <span className="truncate">{macro.name}</span>
      <span className="text-xs text-neutral-500">
        {isLoading
          ? t('keyboard.macro.playing')
          : t('keyboard.macro.stepsCount', { count: countScriptCommands(macro.script) })}
      </span>
    </div>
  );
};
