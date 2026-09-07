import { useEffect, useState } from 'react';
import { Button, Input, Modal, message } from 'antd';
import { TriangleAlertIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/application.ts';

const OFFICIAL_UPDATE_SOURCE = 'https://cdn.sipeed.com/nanokvm';

type UpdateSourceProps = {
  onSourceChanged: () => void;
  onCustomSourceChange: (isCustom: boolean) => void;
};

// The source editor changes one shared root for both application and firmware
// updates; the server remains authoritative for normalization and persistence.
export const UpdateSource = ({ onSourceChanged, onCustomSourceChange }: UpdateSourceProps) => {
  const { t } = useTranslation();

  const [savedSource, setSavedSource] = useState<api.UpdateSource | null>(null);
  const [url, setUrl] = useState(OFFICIAL_UPDATE_SOURCE);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [isResetConfirmOpen, setIsResetConfirmOpen] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    api
      .getUpdateSource()
      .then((rsp) => {
        if (rsp.code !== 0 || !rsp.data) {
          message.error(t('settings.update.source.loadFailed'));
          return;
        }

        const source = rsp.data as api.UpdateSource;
        setSavedSource(source);
        setUrl(source.url || OFFICIAL_UPDATE_SOURCE);
        onCustomSourceChange(!source.isOfficial && source.enabled);
      })
      .catch(() => message.error(t('settings.update.source.loadFailed')))
      .finally(() => setIsLoading(false));
  }, [onCustomSourceChange, t]);

  function validateUrl(value: string): string {
    const trimmed = value.trim();
    if (!trimmed) return t('settings.update.source.invalidUrl');

    try {
      const parsed = new URL(trimmed);
      if (
        (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') ||
        !parsed.hostname ||
        parsed.username ||
        parsed.password ||
        parsed.search ||
        parsed.hash ||
        parsed.pathname.replace(/\/+$/, '').endsWith('.json')
      ) {
        return t('settings.update.source.invalidUrl');
      }
    } catch {
      return t('settings.update.source.invalidUrl');
    }

    return '';
  }

  function requestSave() {
    if (isLoading || isSaving) return;

    const validationError = validateUrl(url);
    setError(validationError);
    if (validationError) return;

    setIsConfirmOpen(true);
  }

  function save() {
    const sourceUrl = url.trim();
    setIsSaving(true);

    api
      .setUpdateSource(sourceUrl)
      .then((rsp) => {
        if (rsp.code !== 0 || !rsp.data) {
          message.error(rsp.msg || t('settings.update.source.saveFailed'));
          return;
        }

        const source = rsp.data as api.UpdateSource;
        setSavedSource(source);
        setUrl(source.url);
        setError('');
        onCustomSourceChange(!source.isOfficial && source.enabled);
        message.success(t('settings.update.source.saved'));
        onSourceChanged();
      })
      .catch(() => message.error(t('settings.update.source.saveFailed')))
      .finally(() => {
        setIsSaving(false);
        setIsConfirmOpen(false);
      });
  }

  function requestReset() {
    if (isLoading || isSaving || savedSource?.isOfficial) return;
    setIsResetConfirmOpen(true);
  }

  function reset() {
    if (isLoading || isSaving || savedSource?.isOfficial) return;

    setIsSaving(true);
    api
      .resetUpdateSource()
      .then((rsp) => {
        if (rsp.code !== 0 || !rsp.data) {
          message.error(rsp.msg || t('settings.update.source.resetFailed'));
          return;
        }

        const source = rsp.data as api.UpdateSource;
        setSavedSource(source);
        setUrl(source.url || OFFICIAL_UPDATE_SOURCE);
        setError('');
        onCustomSourceChange(false);
        message.success(t('settings.update.source.resetDone'));
        onSourceChanged();
      })
      .catch(() => message.error(t('settings.update.source.resetFailed')))
      .finally(() => {
        setIsSaving(false);
        setIsResetConfirmOpen(false);
      });
  }

  const isChanged = Boolean(savedSource && url.trim() !== savedSource.url);
  const canReset = Boolean(savedSource && !savedSource.isOfficial);
  const isHttpSource = url.trim().toLowerCase().startsWith('http://');

  return (
    <>
      <div className="flex flex-col gap-3 py-3">
        <div className="space-y-1">
          <div>{t('settings.update.source.title')}</div>
          <div className="text-xs text-neutral-500">{t('settings.update.source.desc')}</div>
        </div>

        <div className="flex flex-col gap-2 sm:flex-row">
          <div className="min-w-0 flex-1">
            <Input
              value={url}
              status={error ? 'error' : undefined}
              disabled={isLoading || isSaving}
              placeholder={OFFICIAL_UPDATE_SOURCE}
              onChange={(event) => {
                setUrl(event.target.value);
                setError('');
              }}
              onPressEnter={requestSave}
            />
            {error && <div className="mt-1 text-xs text-red-500">{error}</div>}
          </div>

          <Button
            className="shrink-0"
            type="primary"
            loading={isSaving}
            disabled={!isChanged || isLoading}
            onClick={requestSave}
          >
            {t('settings.update.source.save')}
          </Button>
          <Button
            className="shrink-0"
            loading={isSaving}
            disabled={!canReset || isLoading}
            onClick={requestReset}
          >
            {t('settings.update.source.reset')}
          </Button>
        </div>
      </div>

      <Modal
        open={isConfirmOpen}
        centered
        title={
          <div className="flex items-center gap-2 text-amber-500">
            <TriangleAlertIcon size={18} />
            <span>{t('settings.update.source.confirmTitle')}</span>
          </div>
        }
        okText={t('settings.update.source.confirm')}
        cancelText={t('settings.update.cancel')}
        okButtonProps={{ danger: true }}
        confirmLoading={isSaving}
        onOk={save}
        onCancel={() => !isSaving && setIsConfirmOpen(false)}
      >
        <div className="space-y-4 py-4">
          <p>{t('settings.update.source.confirmDesc')}</p>
          <div className="break-all rounded bg-neutral-800 p-3 font-mono text-sm">
            {url.trim()}
          </div>
          {isHttpSource && (
            <div className="rounded bg-red-500/10 p-3 text-sm text-red-400">
              {t('settings.update.source.httpWarning')}
            </div>
          )}
        </div>
      </Modal>

      <Modal
        open={isResetConfirmOpen}
        centered
        title={t('settings.update.source.resetConfirmTitle')}
        okText={t('settings.update.source.resetConfirm')}
        cancelText={t('settings.update.cancel')}
        confirmLoading={isSaving}
        onOk={reset}
        onCancel={() => !isSaving && setIsResetConfirmOpen(false)}
      >
        <p className="py-4">{t('settings.update.source.resetConfirmDesc')}</p>
      </Modal>
    </>
  );
};
