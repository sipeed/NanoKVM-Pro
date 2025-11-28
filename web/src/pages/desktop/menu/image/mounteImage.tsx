import { useEffect, useState } from 'react';
import { Divider, Modal, Segmented, Switch } from 'antd';
import clsx from 'clsx';
import { DiscIcon, HardDriveIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/storage';

import { Images } from './images';

export function MounteImage() {
  const { t } = useTranslation();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isMounted, setIsMounted] = useState(false);
  const [mode, setMode] = useState('mass-storage');
  const [readOnly, setReadOnly] = useState(true);

  const modes = [
    {
      value: 'mass-storage',
      label: (
        <div className="flex items-center space-x-1">
          <HardDriveIcon size={16} />
          <span>Mass Storage</span>
        </div>
      )
    },
    {
      value: 'cd-rom',
      label: (
        <div className="flex items-center space-x-1">
          <DiscIcon size={16} />
          <span>CD ROM</span>
        </div>
      )
    }
  ];

  useEffect(() => {
    api.getMountedImage().then((rsp) => {
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }
      if (!rsp?.data?.file) {
        return;
      }

      setIsMounted(true);
      setMode(rsp.data?.cdrom ? 'cd-rom' : 'mass-storage');
      setReadOnly(!!rsp.data?.readOnly);
    });
  }, []);

  function updateMode(value: string) {
    setMode(value);
    if (value === 'cd-rom' && !readOnly) {
      setReadOnly(true);
    }
  }

  function updateReadOnly(value: boolean) {
    if (mode === 'cd-rom' && !value) {
      return;
    }
    setReadOnly(value);
  }

  return (
    <>
      <div
        onClick={() => setIsModalOpen(true)}
        className="flex h-[30px] cursor-pointer select-none items-center space-x-2 rounded py-0.5 pl-3 pr-6 hover:bg-neutral-700/70"
      >
        <HardDriveIcon size={16} className="text-neutral-400" />
        <span
          className={clsx(
            'text-sm',
            isMounted ? 'text-blue-500' : 'text-neutral-300 hover:text-white'
          )}
        >
          {t('image.mount')}
        </span>
      </div>

      <Modal
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        title={t('image.mount')}
        footer={false}
      >
        <Divider style={{ margin: '24px 0' }} />

        <div className="flex flex-col space-y-6">
          {/* mode */}
          <div className="flex items-center justify-between">
            <span>{t('image.mountMode')}</span>
            <Segmented value={mode} options={modes} disabled={isMounted} onChange={updateMode} />
          </div>

          {/* read only */}
          <div className="flex items-center justify-between">
            <span>{t('image.readOnly')}</span>
            <Switch
              value={readOnly}
              disabled={isMounted || mode === 'cd-rom'}
              onChange={updateReadOnly}
            />
          </div>

          <Divider style={{ margin: '24px 0 0 0' }} />

          {/* image list */}
          <Images
            isOpen={isModalOpen}
            cdrom={mode === 'cd-rom'}
            readOnly={readOnly}
            setIsMounted={setIsMounted}
          />
        </div>
      </Modal>
    </>
  );
}
