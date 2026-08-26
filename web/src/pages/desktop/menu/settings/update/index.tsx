import { useCallback, useEffect, useRef, useState } from 'react';
import { LoadingOutlined, RocketOutlined, SmileOutlined } from '@ant-design/icons';
import { Button, Divider, Result, Spin } from 'antd';
import { useTranslation } from 'react-i18next';
import semver from 'semver';

import * as api from '@/api/application.ts';

import { ManualUpdate } from './manual-update.tsx';
import { Preview } from './preview.tsx';
import { UpdateSource } from './update-source.tsx';
import { Updating } from './updating.tsx';

type UpdateProps = {
  setIsLocked: (isClosable: boolean) => void;
};

type Status = '' | 'loading' | 'updating' | 'outdated' | 'latest' | 'failed';

export const Update = ({ setIsLocked }: UpdateProps) => {
  const { t } = useTranslation();

  const [status, setStatus] = useState<Status>('loading');
  const [currentVersion, setCurrentVersion] = useState('');
  const [latestVersion, setLatestVersion] = useState('');
  const [errMsg, setErrMsg] = useState('');
  const [tipMsg, setTipMsg] = useState('');
  const [isCustomSource, setIsCustomSource] = useState(false);
  const [manualJobID, setManualJobID] = useState('');
  // Ignore late responses from an earlier check after the source changes.
  const versionRequestID = useRef(0);

  const checkForUpdates = useCallback(() => {
    const requestID = ++versionRequestID.current;
    setStatus('loading');

    api
      .getVersion()
      .then((rsp: any) => {
        if (requestID !== versionRequestID.current) return;

        if (rsp.code !== 0 || !rsp.data) {
          setStatus('failed');
          setErrMsg(t('settings.update.queryFailed'));
          return;
        }

        const current = typeof rsp.data.current === 'string' ? semver.valid(rsp.data.current) : null;
        const latest = typeof rsp.data.latest === 'string' ? semver.valid(rsp.data.latest) : null;
        if (!current || !latest) {
          setStatus('failed');
          setErrMsg(t('settings.update.queryFailed'));
          return;
        }

        setCurrentVersion(rsp.data.current);
        setLatestVersion(rsp.data.latest);

        const isLatest = semver.gt(latest, current);
        if (isLatest) {
          setTipMsg(t('settings.update.available'));
        }
        setStatus(!isLatest ? 'latest' : 'outdated');
      })
      .catch(() => {
        if (requestID !== versionRequestID.current) return;

        setStatus('failed');
        setErrMsg(t('settings.update.queryFailed'));
      });
  }, [t]);

  useEffect(() => {
    checkForUpdates();
  }, [checkForUpdates]);

  function update() {
    if (status !== 'outdated') return;

    setManualJobID('');
    setIsLocked(true);
    setStatus('updating');

    api.update().then((rsp: any) => {
      if (rsp.code !== 0) {
        setStatus('failed');
        setErrMsg(t('settings.update.updateFailed'));
      }
    });
  }

  // Manual jobs use the same progress screen but are polled instead of using
  // the online updater's WebSocket progress events.
  const startManualUpdate = useCallback((jobID: string) => {
    setManualJobID(jobID);
    setErrMsg('');
    setStatus('updating');
  }, []);

  const failManualUpdate = useCallback((message: string) => {
    setIsLocked(false);
    setStatus('failed');
    setErrMsg(message);
  }, [setIsLocked]);

  const finishManualUpdate = useCallback(() => {
    setIsLocked(false);
    setManualJobID('');
  }, [setIsLocked]);

  const handleUpdateSourceChanged = useCallback(() => {
    setCurrentVersion('');
    setLatestVersion('');
    setErrMsg('');
    setTipMsg('');
    setStatus('');
  }, []);

  return (
    <>
      <div className="text-base font-bold">{t('settings.update.title')}</div>
      <Divider className="opacity-50" />

      <Preview checkForUpdates={checkForUpdates} disabled={isCustomSource} />
      <Divider className="opacity-50" />

      <UpdateSource
        onSourceChanged={handleUpdateSourceChanged}
        onCustomSourceChange={setIsCustomSource}
      />
      <Divider className="opacity-50" />

      <ManualUpdate
        disabled={status === 'loading' || status === 'updating'}
        setIsLocked={setIsLocked}
        onInstallStarted={startManualUpdate}
        onInstallFailed={failManualUpdate}
      />
      <Divider className="opacity-50" />

      <div className="flex min-h-[400px] flex-col justify-between">
        {status === 'loading' && (
          <div className="flex justify-center pt-24">
            <Spin indicator={<LoadingOutlined spin />} size="large" />
          </div>
        )}

        {status === 'updating' && (
          <Updating
            manualJobID={manualJobID}
            onManualUpdateFailed={failManualUpdate}
            onManualUpdateSucceeded={finishManualUpdate}
          />
        )}

        {status === '' && (
          <Result
            status="info"
            title={t('settings.update.source.ready')}
            extra={
              <Button onClick={checkForUpdates}>{t('settings.update.title')}</Button>
            }
          />
        )}

        {status === 'latest' && (
          <Result
            status="success"
            icon={<SmileOutlined />}
            title={currentVersion}
            subTitle={t('settings.update.isLatest')}
            extra={[
              <Button key="confirm" onClick={checkForUpdates}>
                {t('settings.update.title')}
              </Button>
            ]}
          />
        )}

        {status === 'outdated' && (
          <Result
            status="warning"
            icon={<RocketOutlined />}
            title={`${currentVersion} -> ${latestVersion}`}
            subTitle={tipMsg}
            extra={[
              <Button key="confirm" type="primary" onClick={update}>
                {t('settings.update.confirm')}
              </Button>
            ]}
          />
        )}

        {status === 'failed' && <Result subTitle={errMsg} />}

        <div className="flex justify-center">
          <Button
            type="link"
            size="small"
            href="https://github.com/sipeed/NanoKVM-Pro/blob/main/CHANGELOG.md"
            target="_blank"
          >
            {t('settings.update.changelog')}
          </Button>
        </div>
      </div>
    </>
  );
};
