import { useEffect, useState } from 'react';
import { Divider, Popover } from 'antd';
import { PlayIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/hid.ts';
import { ScrollArea } from '@/components/ui/scroll-area';

import { Editor } from './editor.tsx';
import { MacroItem } from './item.tsx';
import type { Macro } from './types.ts';

export const Macros = () => {
  const { t } = useTranslation();

  const [isOpen, setIsOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [isPicking, setIsPicking] = useState(false);
  const [macros, setMacros] = useState<Macro[]>([]);

  useEffect(() => {
    if (isOpen) {
      getMacros();
    }
  }, [isOpen]);

  async function getMacros() {
    try {
      const rsp = await api.getMacros();
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      setMacros(rsp.data?.macros ?? []);
    } catch (err) {
      console.log(err);
    }
  }

  async function saveMacro(macro: Macro) {
    const rsp = macro.id
      ? await api.updateMacro(macro.id, macro.name, macro.script)
      : await api.addMacro(macro.name, macro.script);

    if (rsp.code !== 0) {
      console.log(rsp.msg);
      return;
    }

    await getMacros();
  }

  async function delMacro(macro: Macro) {
    if (!macro.id) return;

    try {
      const rsp = await api.deleteMacro(macro.id);
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      await getMacros();
    } catch (err) {
      console.log(err);
    }
  }

  function handleOpenChange(open: boolean) {
    if (open) {
      setIsOpen(true);
      return;
    }
    if (isEditing || isPicking) {
      return;
    }
    setIsOpen(false);
  }

  const content = (
    <ScrollArea className="max-w-[400px] [&>[data-radix-scroll-area-viewport]]:max-h-[350px]">
      {macros.length > 0 ? (
        <>
          {macros.map((macro) => (
            <MacroItem key={macro.id} macro={macro} macros={macros} />
          ))}
          <Divider style={{ margin: '5px 0 5px 0' }} />
        </>
      ) : (
        <div className="px-3 py-2 text-sm text-neutral-500">{t('keyboard.macro.empty')}</div>
      )}

      <Editor
        macros={macros}
        saveMacro={saveMacro}
        delMacro={delMacro}
        setIsEditing={setIsEditing}
        setIsPicking={setIsPicking}
      />
    </ScrollArea>
  );

  return (
    <Popover
      content={content}
      trigger="hover"
      placement="rightTop"
      align={{ offset: [14, 0] }}
      open={isOpen}
      onOpenChange={handleOpenChange}
      destroyOnHidden={false}
      arrow={false}
    >
      <div className="flex cursor-pointer select-none items-center space-x-2 rounded py-1 pl-2 pr-5 hover:bg-neutral-700/70">
        <PlayIcon size={18} />
        <span>{t('keyboard.macro.title')}</span>
      </div>
    </Popover>
  );
};
