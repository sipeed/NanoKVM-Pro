import { useEffect, useRef, useState } from 'react';
import { Button, Divider, message, Modal, Switch } from 'antd';
import { CheckIcon, CopyIcon, EyeIcon, RefreshCcwIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/ai.ts';
import type { AIConfig } from '@/api/ai.ts';

const MCPKeyModalContent = ({ value }: { value: string }) => {
  const { t } = useTranslation();
  const [isCopied, setIsCopied] = useState(false);
  const [copyFailed, setCopyFailed] = useState(false);
  const copyTimer = useRef<ReturnType<typeof setTimeout>>();

  useEffect(() => () => clearTimeout(copyTimer.current), []);

  async function copyKey() {
    setCopyFailed(false);
    try {
      await navigator.clipboard.writeText(value);
      setIsCopied(true);
      clearTimeout(copyTimer.current);
      copyTimer.current = setTimeout(() => setIsCopied(false), 2000);
    } catch {
      setCopyFailed(true);
    }
  }

  return (
    <div className="flex flex-col space-y-3 pt-1">
      <div className="relative rounded-lg border border-neutral-700/60 bg-neutral-900/70 px-3 py-2.5">
        <span className="select-all break-all font-mono text-sm leading-relaxed text-neutral-200">
          {value}
        </span>
        <span aria-hidden="true" className="inline-block h-5 min-h-[26px] w-10 min-w-[28px]" />
        <Button
          type="text"
          size="small"
          className="absolute bottom-2 right-2 text-neutral-400 hover:text-white"
          aria-label={t('settings.ai.copyApiKey')}
          title={t('settings.ai.copyApiKey')}
          icon={
            isCopied ? <CheckIcon size={15} className="!text-green-500" /> : <CopyIcon size={15} />
          }
          onClick={copyKey}
        />
      </div>
      {copyFailed && <span className="text-xs text-red-400">{t('settings.ai.copyManually')}</span>}
    </div>
  );
};

export const AI = () => {
  const { t } = useTranslation();
  const [modal, modalContext] = Modal.useModal();
  const [feedback, feedbackContext] = message.useMessage();
  const [config, setConfig] = useState<AIConfig>({ mcpEnabled: false });
  const [isLoading, setIsLoading] = useState(true);
  const [loadFailed, setLoadFailed] = useState(false);
  const [copied, setCopied] = useState<'endpoint' | 'key' | null>(null);
  const copyTimer = useRef<ReturnType<typeof setTimeout>>();
  const busy = useRef(true);
  const keyModal = useRef<{ destroy: () => void } | null>(null);
  const endpoint = `${window.location.origin}/api/mcp`;

  useEffect(() => {
    let active = true;
    api
      .getAIConfig()
      .then((rsp) => {
        if (!active) return;
        if (rsp.code !== 0) throw new Error(rsp.msg);
        setConfig(rsp.data);
      })
      .catch(() => {
        if (active) {
          setLoadFailed(true);
          void feedback.error(t('settings.ai.loadFailed'));
        }
      })
      .finally(() => {
        if (active) {
          busy.current = false;
          setIsLoading(false);
        }
      });
    return () => {
      active = false;
      clearTimeout(copyTimer.current);
      keyModal.current?.destroy();
    };
  }, [feedback, t]);

  async function run(action: () => Promise<void>) {
    if (busy.current) return;
    busy.current = true;
    setIsLoading(true);
    try {
      await action();
    } catch {
      void feedback.error(t('settings.ai.operationFailed'));
    } finally {
      busy.current = false;
      setIsLoading(false);
    }
  }

  function closeKey() {
    keyModal.current?.destroy();
    keyModal.current = null;
  }

  function showKey(value: string) {
    closeKey();
    keyModal.current = modal.info({
      title: t('settings.ai.apiKey'),
      icon: null,
      centered: true,
      closable: true,
      footer: null,
      width: 420,
      content: <MCPKeyModalContent value={value} />
    });
  }

  async function copy(value: string, target: 'endpoint' | 'key') {
    await navigator.clipboard.writeText(value);
    setCopied(target);
    clearTimeout(copyTimer.current);
    copyTimer.current = setTimeout(() => setCopied(null), 2000);
  }

  function toggle(enabled: boolean) {
    void run(async () => {
      const rsp = await api.setAIConfig(enabled);
      if (rsp.code !== 0) throw new Error(rsp.msg);
      closeKey();
      setConfig(rsp.data);
    });
  }

  function revealKey(shouldCopy: boolean) {
    void run(async () => {
      const rsp = await (shouldCopy ? api.copyMCPAPIKey() : api.getMCPAPIKey());
      if (rsp.code !== 0 || !rsp.data?.mcpApiKey) throw new Error(rsp.msg);
      if (shouldCopy) {
        try {
          await copy(rsp.data.mcpApiKey, 'key');
        } catch {
          showKey(rsp.data.mcpApiKey);
          void feedback.info(t('settings.ai.copyManually'));
        }
      } else {
        showKey(rsp.data.mcpApiKey);
      }
    });
  }

  function regenerateKey() {
    modal.confirm({
      title: t('settings.ai.regenerate'),
      content: (
        <span className="text-xs text-neutral-500">{t('settings.ai.regenerateConfirm')}</span>
      ),
      okText: t('settings.ai.regenerate'),
      cancelText: t('settings.ai.cancel'),
      onOk: () =>
        run(async () => {
          const rsp = await api.regenerateMCPAPIKey();
          if (rsp.code !== 0) throw new Error(rsp.msg);
          closeKey();
          setConfig(rsp.data);
          void feedback.success(t('settings.ai.regenerated'));
        })
    });
  }

  return (
    <>
      <div className="text-base font-bold">{t('settings.ai.title')}</div>
      <Divider className="opacity-50" />
      <div className="flex min-w-0 flex-col space-y-6">
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-col space-y-1">
            <span className="text-sm font-medium">{t('settings.ai.mcpService')}</span>
            <span className="text-xs text-neutral-500">{t('settings.ai.description')}</span>
          </div>
          <Switch
            checked={config.mcpEnabled}
            loading={isLoading}
            disabled={loadFailed}
            onChange={toggle}
          />
        </div>

        {config.mcpEnabled && (
          <div className="flex min-w-0 flex-col overflow-hidden rounded-xl border border-neutral-700/50 bg-neutral-800/40 shadow-sm">
            <div className="group flex min-w-0 flex-col space-y-2 border-b border-neutral-800 px-4 py-3.5 transition-colors hover:bg-neutral-800/40 sm:flex-row sm:items-center sm:justify-between">
              <span className="w-24 shrink-0 text-sm font-medium text-neutral-400">
                {t('settings.ai.endpoint')}
              </span>
              <div className="flex min-w-0 items-center justify-between gap-2">
                <span
                  className="min-w-0 flex-1 select-all truncate font-mono text-sm text-neutral-300"
                  title={endpoint}
                >
                  {endpoint}
                </span>
                <div className="flex shrink-0 items-center opacity-40 transition-opacity group-focus-within:opacity-100 group-hover:opacity-100 sm:ml-4">
                  <Button
                    type="text"
                    size="small"
                    className="text-neutral-400 hover:text-white"
                    aria-label={t('settings.ai.copyEndpoint')}
                    title={t('settings.ai.copyEndpoint')}
                    icon={
                      copied === 'endpoint' ? (
                        <CheckIcon size={15} className="!text-green-500" />
                      ) : (
                        <CopyIcon size={15} />
                      )
                    }
                    onClick={() => {
                      void copy(endpoint, 'endpoint').catch(() => {
                        void feedback.error(t('settings.ai.copyManually'));
                      });
                    }}
                  />
                </div>
              </div>
            </div>
            <div className="group flex min-w-0 flex-col space-y-2 px-4 py-3.5 transition-colors hover:bg-neutral-800/40 sm:flex-row sm:items-center sm:justify-between">
              <span className="w-24 shrink-0 text-sm font-medium text-neutral-400">
                {t('settings.ai.apiKey')}
              </span>
              <div className="flex min-w-0 items-center justify-between gap-2">
                <span className="min-w-0 flex-1 select-all truncate font-mono text-sm text-neutral-300">
                  {config.mcpApiKeyMasked || '-'}
                </span>
                <div className="flex shrink-0 items-center space-x-1 opacity-40 transition-opacity group-focus-within:opacity-100 group-hover:opacity-100 sm:ml-4">
                  <Button
                    type="text"
                    size="small"
                    className="text-neutral-400 hover:text-white"
                    aria-label={t('settings.ai.showApiKey')}
                    title={t('settings.ai.showApiKey')}
                    disabled={isLoading}
                    icon={<EyeIcon size={15} />}
                    onClick={() => revealKey(false)}
                  />
                  <Button
                    type="text"
                    size="small"
                    className="text-neutral-400 hover:text-white"
                    aria-label={t('settings.ai.copyApiKey')}
                    title={t('settings.ai.copyApiKey')}
                    disabled={isLoading}
                    icon={
                      copied === 'key' ? (
                        <CheckIcon size={15} className="!text-green-500" />
                      ) : (
                        <CopyIcon size={15} />
                      )
                    }
                    onClick={() => revealKey(true)}
                  />
                  <Button
                    type="text"
                    size="small"
                    className="text-neutral-400 hover:text-white"
                    aria-label={t('settings.ai.regenerate')}
                    title={t('settings.ai.regenerate')}
                    disabled={isLoading}
                    icon={<RefreshCcwIcon size={15} />}
                    onClick={regenerateKey}
                  />
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
      {modalContext}
      {feedbackContext}
    </>
  );
};
