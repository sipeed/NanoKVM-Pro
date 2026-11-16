import { useEffect, useRef, useState } from 'react';
import { Button, Divider, Input, Modal } from 'antd';
import { useSetAtom } from 'jotai';
import { Trash2Icon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { isKeyboardEnableAtom } from '@/jotai/keyboard.ts';
import { ScrollArea } from '@/components/ui/scroll-area.tsx';

import { CoordinatePicker } from './picker.tsx';
import { parseScript } from './script.ts';
import type { Macro } from './types.ts';

interface EditorProps {
  macros: Macro[];
  saveMacro: (macro: Macro) => Promise<void>;
  delMacro: (macro: Macro) => Promise<void>;
  setIsEditing: (isEditing: boolean) => void;
  setIsPicking: (isPicking: boolean) => void;
}

type InsertKind = 'CLICK' | 'DBLCLICK';

const emptyMacro = (): Macro => ({ name: '', script: '' });

export const Editor = ({ macros, saveMacro, delMacro, setIsEditing, setIsPicking }: EditorProps) => {
  const { t } = useTranslation();
  const setIsKeyboardEnable = useSetAtom(isKeyboardEnableAtom);
  const textareaRef = useRef<HTMLTextAreaElement | null>(null);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [draft, setDraft] = useState<Macro>(emptyMacro);
  const [pickingKind, setPickingKind] = useState<InsertKind | null>(null);
  const [isSaving, setIsSaving] = useState(false);
  const [parseError, setParseError] = useState('');

  const isPicking = pickingKind !== null;
  const parsed = parseScript(draft.script);
  const canSave = !!draft.name.trim() && parsed.ok;

  useEffect(() => {
    setIsKeyboardEnable(!isModalOpen && !isPicking);
    setIsEditing(isModalOpen);
    setIsPicking(isPicking);
  }, [isModalOpen, isPicking, setIsEditing, setIsKeyboardEnable, setIsPicking]);

  function openModal(macro?: Macro) {
    setDraft(macro ? { ...macro } : emptyMacro());
    setPickingKind(null);
    setParseError('');
    setIsModalOpen(true);
  }

  function closeModal() {
    if (isPicking) return;
    setDraft(emptyMacro());
    setParseError('');
    setIsModalOpen(false);
  }

  function insertSnippet(snippet: string, withNewline = true) {
    const el = textareaRef.current;
    const value = draft.script;
    const start = el?.selectionStart ?? value.length;
    const end = el?.selectionEnd ?? start;
    const before = value.slice(0, start);
    const after = value.slice(end);
    const prefix = before && !before.endsWith('\n') ? '\n' : '';
    const suffix = withNewline && after && !after.startsWith('\n') ? '\n' : '';
    const next = `${before}${prefix}${snippet}${suffix}${after}`;
    const cursor = before.length + prefix.length + snippet.length;

    setDraft((prev) => ({ ...prev, script: next }));
    setParseError('');

    requestAnimationFrame(() => {
      el?.focus();
      el?.setSelectionRange(cursor, cursor);
    });
  }

  function startPicking(kind: InsertKind) {
    setPickingKind(kind);
    setIsModalOpen(false);
  }

  function handlePicked(x: number, y: number) {
    if (!pickingKind) return;
    const line = `${pickingKind} ${x} ${y}`;
    setDraft((prev) => {
      const script = prev.script.trimEnd();
      const prefix = script ? `${script}\n` : '';
      return { ...prev, script: `${prefix}${line}\n` };
    });
    setPickingKind(null);
    setIsModalOpen(true);
  }

  function cancelPicking() {
    setPickingKind(null);
    setIsModalOpen(true);
  }

  async function handleSave() {
    const result = parseScript(draft.script);
    if (isSaving || !draft.name.trim() || !result.ok) {
      if (!result.ok) {
        setParseError(t('keyboard.macro.parseError', { line: result.line, error: result.error }));
      }
      return;
    }

    setIsSaving(true);
    try {
      await saveMacro({
        id: draft.id,
        name: draft.name.trim(),
        script: draft.script.replace(/\r\n/g, '\n')
      });
      setDraft(emptyMacro());
      setParseError('');
    } catch (err) {
      console.log(err);
    } finally {
      setIsSaving(false);
    }
  }

  return (
    <>
      <div
        className="flex h-[30px] cursor-pointer items-center space-x-1 rounded px-3 text-neutral-300 hover:bg-neutral-700/60"
        onClick={() => openModal()}
      >
        <span>{t('keyboard.macro.custom')}</span>
      </div>

      <Modal
        width={640}
        title={draft.id ? t('keyboard.macro.edit') : t('keyboard.macro.title')}
        keyboard={false}
        footer={null}
        open={isModalOpen}
        onCancel={closeModal}
      >
        <div className="flex flex-col space-y-3">
          <Input
            maxLength={64}
            placeholder={t('keyboard.macro.namePlaceholder')}
            value={draft.name}
            onChange={(e) => setDraft((prev) => ({ ...prev, name: e.target.value }))}
          />

          <pre className="whitespace-pre-wrap rounded bg-neutral-800/70 px-3 py-2 text-xs leading-5 text-neutral-400">
            {t('keyboard.macro.tips')}
          </pre>

          <div className="flex flex-wrap gap-2">
            <Button size="small" onClick={() => insertSnippet('STRING ', false)}>
              STRING
            </Button>
            <Button size="small" onClick={() => insertSnippet('ENTER')}>
              ENTER
            </Button>
            <Button size="small" onClick={() => insertSnippet('TAB')}>
              TAB
            </Button>
            <Button size="small" onClick={() => insertSnippet('CTRL-ALT DELETE')}>
              CTRL-ALT-DEL
            </Button>
            <Button size="small" onClick={() => insertSnippet('DELAY 200')}>
              DELAY
            </Button>
            <Button size="small" onClick={() => insertSnippet('RUN ', false)}>
              RUN
            </Button>
            <Button size="small" onClick={() => startPicking('CLICK')}>
              CLICK
            </Button>
            <Button size="small" onClick={() => startPicking('DBLCLICK')}>
              DBLCLICK
            </Button>
          </div>

          <textarea
            ref={textareaRef}
            spellCheck={false}
            value={draft.script}
            placeholder={t('keyboard.macro.placeholder')}
            className="h-[280px] w-full resize-y rounded border border-neutral-700 bg-neutral-950 p-3 font-mono text-sm leading-6 text-neutral-200 outline-none placeholder:text-neutral-600 focus:border-neutral-500"
            onChange={(e) => {
              setDraft((prev) => ({ ...prev, script: e.target.value }));
              setParseError('');
            }}
          />

          {(parseError || (!parsed.ok && draft.script.trim())) && (
            <div className="text-sm text-red-400">
              {parseError ||
                t('keyboard.macro.parseError', { line: parsed.ok ? 0 : parsed.line, error: parsed.ok ? '' : parsed.error })}
            </div>
          )}

          <div className="flex justify-center pb-1">
            <Button type="primary" className="min-w-24" loading={isSaving} disabled={!canSave} onClick={handleSave}>
              {t('keyboard.macro.save')}
            </Button>
          </div>
        </div>

        {macros.length > 0 && (
          <>
            <Divider />
            <ScrollArea className="[&>[data-radix-scroll-area-viewport]]:max-h-[220px]">
              {macros.map((macro) => (
                <div
                  key={macro.id}
                  className="flex items-center justify-between rounded p-2 hover:bg-neutral-700/50"
                >
                  <div className="min-w-0">
                    <div className="truncate">{macro.name}</div>
                    <div className="truncate font-mono text-xs text-neutral-500">
                      {macro.script.trim().split('\n')[0]}
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Button size="small" onClick={() => openModal(macro)}>
                      {t('keyboard.macro.edit')}
                    </Button>
                    <div
                      className="flex size-[20px] cursor-pointer items-center justify-center rounded-sm text-neutral-500 hover:text-red-500"
                      onClick={() => macro.id && delMacro(macro)}
                    >
                      <Trash2Icon size={16} />
                    </div>
                  </div>
                </div>
              ))}
            </ScrollArea>
          </>
        )}
      </Modal>

      <CoordinatePicker open={isPicking} onPick={handlePicked} onCancel={cancelPicking} />
    </>
  );
};
