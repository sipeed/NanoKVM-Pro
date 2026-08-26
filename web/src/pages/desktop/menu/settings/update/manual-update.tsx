import { ChangeEvent, useRef, useState } from 'react';
import { UploadOutlined } from '@ant-design/icons';
import { Button, Modal, message } from 'antd';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/application.ts';

type ManualUpdateProps = {
  disabled: boolean;
  setIsLocked: (isLocked: boolean) => void;
  onInstallStarted: (jobID: string) => void;
  onInstallFailed: (message: string) => void;
};

function formatBytes(size?: number): string {
  if (size === undefined || size < 0) return '';

  const units = ['B', 'KiB', 'MiB', 'GiB'];
  let value = size;
  let index = 0;
  while (value >= 1024 && index < units.length - 1) {
    value /= 1024;
    index += 1;
  }

  return `${value.toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
}

// Upload is inspected and staged first; installation starts only after the
// confirmation modal has shown the server-detected version and size.
export const ManualUpdate = ({
  disabled,
  setIsLocked,
  onInstallStarted,
  onInstallFailed
}: ManualUpdateProps) => {
  const { t } = useTranslation();
  const inputRef = useRef<HTMLInputElement | null>(null);

  const [isInspecting, setIsInspecting] = useState(false);
  const [isConfirming, setIsConfirming] = useState(false);
  const [inspection, setInspection] = useState<api.ManualUpdateInspection | null>(null);

  function discardInspection() {
    if (!inspection) return;

    const { id } = inspection;
    setInspection(null);
    // Closing a review must not be blocked by cleanup. The server also expires staged jobs.
    api.cancelManualUpdate(id).catch(() => undefined);
  }

  function selectFile() {
    if (disabled || isInspecting || isConfirming) return;

    discardInspection();
    if (inputRef.current) inputRef.current.value = '';
    inputRef.current?.click();
  }

  function inspectFile(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file || disabled || isInspecting || isConfirming) return;

    const formData = new FormData();
    formData.append('file', file, file.name);

    setIsInspecting(true);
    api
      .inspectManualUpdate(formData)
      .then((rsp) => {
        if (rsp.code !== 0 || !rsp.data) {
          message.error(rsp.msg || t('settings.update.manual.inspectFailed'));
          return;
        }

        setInspection(rsp.data as api.ManualUpdateInspection);
      })
      .catch(() => message.error(t('settings.update.manual.inspectFailed')))
      .finally(() => setIsInspecting(false));
  }

  function confirmUpdate() {
    if (!inspection || isConfirming) return;

    setIsConfirming(true);
    api
      .confirmManualUpdate(inspection.id)
      .then((rsp) => {
        if (rsp.code !== 0) {
          const error = rsp.msg || t('settings.update.manual.installFailed');
          discardInspection();
          onInstallFailed(error);
          return;
        }

        setInspection(null);
        setIsLocked(true);
        onInstallStarted(inspection.id);
      })
      .catch(() => {
        discardInspection();
        onInstallFailed(t('settings.update.manual.installFailed'));
      })
      .finally(() => setIsConfirming(false));
  }

  return (
    <>
      <div className="flex items-center justify-between gap-4 py-3">
        <div className="min-w-0 space-y-1">
          <div>{t('settings.update.manual.title')}</div>
          <div className="text-xs text-neutral-500">{t('settings.update.manual.desc')}</div>
        </div>
        <input
          ref={inputRef}
          type="file"
          accept=".tar.gz,.tar.xz,application/gzip,application/x-xz"
          className="hidden"
          onChange={inspectFile}
        />
        <Button
          className="shrink-0"
          icon={<UploadOutlined />}
          loading={isInspecting}
          disabled={disabled || isConfirming}
          onClick={selectFile}
        >
          {t('settings.update.manual.upload')}
        </Button>
      </div>

      <Modal
        open={inspection !== null}
        centered
        closable={!isConfirming}
        maskClosable={!isConfirming}
        title={t('settings.update.manual.confirmTitle')}
        okText={t('settings.update.manual.install')}
        cancelText={t('settings.update.cancel')}
        confirmLoading={isConfirming}
        onOk={confirmUpdate}
        onCancel={() => !isConfirming && discardInspection()}
      >
        {inspection && (
          <div className="space-y-4 py-4">
            <p>{t('settings.update.manual.confirmDesc')}</p>
            <p>{t('settings.update.manual.version', { version: inspection.version })}</p>
            {inspection.size !== undefined && (
              <p>{t('settings.update.manual.size', { size: formatBytes(inspection.size) })}</p>
            )}
          </div>
        )}
      </Modal>
    </>
  );
};
