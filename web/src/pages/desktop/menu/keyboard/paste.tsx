import { ChangeEvent, useRef, useState } from 'react';
import { Button, Input, Modal, Select, type InputRef } from 'antd';
import clsx from 'clsx';
import { useSetAtom } from 'jotai';
import { ClipboardIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { paste } from '@/api/hid';
import { isKeyboardEnableAtom } from '@/jotai/keyboard.ts';

const { TextArea } = Input;

const languages = [
  { value: 'en', labelKey: 'keyboard.dropdownEnglish' },
  { value: 'ru', labelKey: 'keyboard.dropdownRussian' }
];

// Extended RU → EN translation including punctuation (Russian keyboard layout to US key codes)
function translateRuToEnWithPunctuation(value: string): string {
  const letterMap: Record<string, string> = {
    ё: '`',
    й: 'q', ц: 'w', у: 'e', к: 'r', е: 't', н: 'y', г: 'u', ш: 'i', щ: 'o', з: 'p',
    х: '[', ъ: ']',
    ф: 'a', ы: 's', в: 'd', а: 'f', п: 'g', р: 'h', о: 'j', л: 'k', д: 'l', ж: ';', э: "'",
    я: 'z', ч: 'x', с: 'c', м: 'v', и: 'b', т: 'n', ь: 'm', б: ',', ю: '.'
  };
  const punctuationMap: Record<string, string> = {
    '"': '@', '№': '#', ';': '$', ':': '^', '?': '&', 'Ё': '~', '/': '|', '.': '/', ',': '?'
  };
  return Array.from(value)
    .map((ch) => {
      const lower = ch.toLowerCase();
      if (letterMap[lower]) {
        const translated = letterMap[lower];
        return ch === lower ? translated : translated.toUpperCase();
      }
      if (punctuationMap[ch]) return punctuationMap[ch];
      return ch;
    })
    .join('');
}

function isValidForLanguage(value: string, lang: string): boolean {
  const isRussian = lang === 'ru';
  for (const ch of value) {
    const code = ch.codePointAt(0) ?? 0;
    if (isRussian) {
      if (
        (code >= 0x0410 && code <= 0x042f) || // А-Я
        (code >= 0x0430 && code <= 0x044f) || // а-я
        code === 0x0401 || // Ё
        code === 0x0451 || // ё
        (code >= 0x20 && code <= 0x7e && !(code >= 0x41 && code <= 0x5a) && !(code >= 0x61 && code <= 0x7a))
      ) {
        continue;
      }
      return false;
    }
    if (code <= 0x7f) continue;
    return false;
  }
  return true;
}

export const Paste = () => {
  const { t } = useTranslation();
  const setIsKeyboardEnable = useSetAtom(isKeyboardEnableAtom);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const [language, setLanguage] = useState('en');
  const [status, setStatus] = useState<'' | 'error'>('');
  const [isLoading, setIsLoading] = useState(false);
  const [errMsg, setErrMsg] = useState('');

  const inputRef = useRef<InputRef>(null);

  function onChange(e: ChangeEvent<HTMLTextAreaElement>) {
    const value = e.target.value;
    setStatus(isValidForLanguage(value, language) ? '' : 'error');
    setInputValue(value);
  }

  function onLanguageChange(value: string) {
    setLanguage(value);
    setStatus(inputValue ? (isValidForLanguage(inputValue, value) ? '' : 'error') : '');
  }

  function submit() {
    if (isLoading || !inputValue) return;
    setIsLoading(true);
    const textToSend = language === 'ru' ? translateRuToEnWithPunctuation(inputValue) : inputValue;

    paste(textToSend)
      .then((rsp) => {
        if (rsp.code !== 0) {
          setErrMsg(rsp.msg);
          return;
        }
        setInputValue('');
        setStatus('');
        setErrMsg('');
        setIsModalOpen(false);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  function afterOpenChange(open: boolean) {
    if (open) {
      inputRef.current?.focus();
    }
    setIsKeyboardEnable(!open);
  }

  return (
    <>
      <div
        className={clsx(
          'flex cursor-pointer select-none items-center space-x-2 rounded py-1 pl-2 pr-5 hover:bg-neutral-700/70'
        )}
        onClick={() => setIsModalOpen(true)}
      >
        <ClipboardIcon size={18} />
        <span>{t('keyboard.paste')}</span>
      </div>

      <Modal
        open={isModalOpen}
        centered={false}
        title={t('keyboard.paste')}
        footer={null}
        onCancel={() => setIsModalOpen(false)}
        afterOpenChange={afterOpenChange}
      >
        <div className="flex items-center gap-2 pb-2">
          <span className="text-sm text-neutral-600">{t('keyboard.virtual')}:</span>
          <Select
            size="small"
            style={{ minWidth: 120 }}
            value={language}
            options={languages.map((l) => ({ value: l.value, label: t(l.labelKey) }))}
            onChange={onLanguageChange}
          />
        </div>
        <div className="pb-3 text-xs text-neutral-500">{t('keyboard.tips')}</div>

        <TextArea
          ref={inputRef}
          value={inputValue}
          status={status}
          showCount
          maxLength={1024}
          autoSize={{ minRows: 5, maxRows: 12 }}
          placeholder={t('keyboard.placeholder')}
          onChange={onChange}
        />

        {errMsg && <div className="pt-1 text-sm text-red-500">{errMsg}</div>}

        <div className="flex justify-center py-3">
          <Button type="primary" loading={isLoading} htmlType="submit" onClick={submit}>
            {t('keyboard.submit')}
          </Button>
        </div>
      </Modal>
    </>
  );
};
