import { ChangeEvent, useEffect, useRef, useState } from 'react';
import { Button, Divider, Input, InputRef, Modal } from 'antd';
import clsx from 'clsx';
import { useSetAtom } from 'jotai';
import { HardDriveDownloadIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { downloadImage, imageEnabled, statusImage } from '@/api/storage';
import { isKeyboardEnableAtom } from '@/jotai/keyboard';

export const DownloadImage = () => {
  const [modal, setModal] = useState(false);
  const { t } = useTranslation();
  const setIsKeyboardEnable = useSetAtom(isKeyboardEnableAtom);

  const [input, setInput] = useState('');
  const [status, setStatus] = useState('');
  const [log, setLog] = useState('');
  const [diskEnabled, setDiskEnabled] = useState(false);
  const [popoverKey, setPopoverKey] = useState(0);

  const inputRef = useRef<InputRef>(null);

  const intervalId = useRef<ReturnType<typeof setInterval> | undefined>(undefined);

  useEffect(() => {
    checkDiskEnabled();
  }, []);

  function checkDiskEnabled() {
    imageEnabled()
      .then((res) => {
        setDiskEnabled(res.data.enabled);
      })
      .catch(() => {
        setDiskEnabled(false);
      });
  }

  function handleOpenChange(open: boolean) {
    if (open) {
      clearInterval(intervalId.current);
      checkDiskEnabled();
      getDownloadStatus();
      if (!intervalId.current) {
        intervalId.current = setInterval(getDownloadStatus, 2500);
      }
      setIsKeyboardEnable(false);
      setPopoverKey((prevKey) => prevKey + 1); // Force re-render
    } else {
      setInput('');
      setStatus('');
      setLog('');

      setIsKeyboardEnable(true);
      clearInterval(intervalId.current);
      intervalId.current = undefined;
    }
  }

  function handleChange(e: ChangeEvent<HTMLInputElement>) {
    setInput(e.target.value);
  }

  function getDownloadStatus() {
    statusImage().then((rsp) => {
      if (rsp.data.status) {
        setStatus(rsp.data.status);
        if (rsp.data.status === 'in_progress') {
          // Check if rsp has a percentage value
          if (rsp.data.percentage) {
            setLog('Downloading (' + rsp.data.percentage + ')' + ': ' + rsp.data.file);
          } else {
            setLog('Downloading' + ': ' + rsp.data.file);
          }
          setInput(rsp.data.file);
        }
        if (rsp.data.status === 'failed') {
          setLog('Failed');
          clearInterval(intervalId.current);
        }
        if (rsp.data.status === 'idle') {
          setLog(''); // Clear the log
          clearInterval(intervalId.current);
        }
      }
    });
  }

  function download(url?: string) {
    if (!url) return;

    setStatus('in_progress');
    setLog('Downloading: ' + url);
    // start the getDownloadStatus to tick every 5 seconds

    downloadImage(url)
      .then((rsp) => {
        if (rsp.code !== 0) {
          setStatus('failed');
          setLog(rsp.msg || 'Failed');
          return;
        }

        getDownloadStatus();
        // Start the interval to check the download status
        if (!intervalId.current) {
          intervalId.current = setInterval(getDownloadStatus, 2500);
        }
      })
      .catch(() => {
        clearInterval(intervalId.current); // Clear the interval when the download is complete or fails
        setStatus('failed');
        setLog('Failed');
      });
  }
  return (
    <>
      <div
        onClick={() => setModal(true)}
        className="flex h-[30px] cursor-pointer select-none items-center space-x-2 rounded pl-3 pr-6 text-neutral-300 hover:bg-neutral-700/70"
      >
        <HardDriveDownloadIcon size={18} />
        <span className="select-none text-sm">{t('image.download')}</span>
      </div>

      <Modal
        open={modal}
        afterOpenChange={handleOpenChange}
        onCancel={() => setModal(false)}
        title={t('download.title')}
        footer={false}
      >
        <div key={popoverKey} className="min-w-[300px]">
          <Divider style={{ margin: '10px 0 10px 0' }} />

          {!diskEnabled ? (
            <div className="text-red-500">{t('download.disabled')}</div>
          ) : (
            <>
              <div className="pb-1 text-neutral-500">{t('download.input')}</div>
              <div className="flex items-center space-x-1">
                <Input
                  ref={inputRef}
                  value={input}
                  onChange={handleChange}
                  disabled={status === 'in_progress'}
                />
                <Button
                  type="primary"
                  onClick={() => download(input)}
                  disabled={status === 'in_progress'}
                >
                  {t('download.ok')}
                </Button>
              </div>
            </>
          )}
          <div className={clsx('py-2')}>
            {status && (
              <div
                className={clsx(
                  'max-w-[300px] break-words text-sm',
                  status === 'failed' ? 'text-red-500' : 'text-green-500'
                )}
              >
                {log}
              </div>
            )}
          </div>
        </div>
      </Modal>
    </>
  );
};
