import { useEffect, useState } from 'react';
import { LoadingOutlined, RocketOutlined, SmileOutlined } from '@ant-design/icons';
import { Button, Divider, Result, Spin } from 'antd';
import { useTranslation } from 'react-i18next';
import semver from 'semver';

import * as api from '@/api/application.ts';

import { Preview } from './preview.tsx';
import { Updating } from './updating.tsx';

type UpdateProps = {
  setIsLocked: (isClosable: boolean) => void;
};

type Status = '' | 'loading' | 'updating' | 'outdated' | 'latest' | 'failed';

export const Update = ({ setIsLocked }: UpdateProps) => {
  const { t } = useTranslation();

  const [status, setStatus] = useState<Status>('');
  const [currentVersion, setCurrentVersion] = useState('');
  const [latestVersion, setLatestVersion] = useState('');
  const [errMsg, setErrMsg] = useState('');
  const [tipMsg, setTipMsg] = useState('');

  useEffect(() => {
    checkForUpdates();
  }, []);

  function checkForUpdates() {
    if (status === 'loading') return;
    setStatus('loading');

    api
      .getVersion()
      .then((rsp: any) => {
        if (rsp.code !== 0 || !rsp.data) {
          setStatus('failed');
          setErrMsg(t('settings.update.queryFailed'));
          return;
        }

        setCurrentVersion(rsp.data.current);
        setLatestVersion(rsp.data.latest);

        const isLatest = semver.gt(rsp.data.latest, rsp.data.current);
        if (isLatest) {
          setTipMsg(t('settings.update.available'));
        }
        setStatus(!isLatest ? 'latest' : 'outdated');
      })
      .catch(() => {
        setStatus('failed');
        setErrMsg(t('settings.update.queryFailed'));
      });
  }

  function update() {
    if (status !== 'outdated') return;

    setIsLocked(true);
    setStatus('updating');

    api.update().then((rsp: any) => {
      if (rsp.code !== 0) {
        setStatus('failed');
        setErrMsg(t('settings.update.updateFailed'));
      }
    });
  }

  return (
    <>
      <div className="text-base font-bold">{t('settings.update.title')}</div>
      <Divider />

      <Preview checkForUpdates={checkForUpdates} />

      <div className="my-[40px] h-px bg-neutral-500/10" />

      <div className="flex min-h-[400px] flex-col justify-between">
        {status === 'loading' && (
          <div className="flex justify-center pt-24">
            <Spin indicator={<LoadingOutlined spin />} size="large" />
          </div>
        )}

        {status === 'updating' && <Updating />}

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
