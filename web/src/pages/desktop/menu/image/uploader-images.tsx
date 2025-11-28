import React, { useEffect, useRef, useState } from 'react';
import { UploadOutlined } from '@ant-design/icons';
import { Button, Divider, Modal, notification } from 'antd';
import clsx from 'clsx';
import { HardDriveUploadIcon, PackageIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { getImages, uploadImage } from '@/api/storage';

export const UploaderImages = () => {
  const { t } = useTranslation();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [notify] = notification.useNotification();
  const [images, setImages] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const inputRef = useRef<any>(null);
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    fetchImages();
  }, []);

  function fetchImages() {
    getImages().then((response) => {
      if (response.code === 0 && response.data && response.data.files) {
        setImages(response.data.files);
      } else {
        notify.error({
          message: t('image.fetchError'),
          description: response.msg,
          duration: 3
        });
      }
    });
  }

  async function uploadFile(event: React.ChangeEvent<HTMLInputElement>) {
    const files = event.target.files;
    if (!files || files.length === 0) return;

    const chunkSize = 8 * 1024 * 1024;
    setIsLoading(true);
    const file = files[0];
    let fileId = '';
    let chunkIndex = 0;
    console.log(file);
    try {
      const totalChunks = Math.ceil(file.size / chunkSize);

      while (chunkIndex < totalChunks) {
        const start = chunkIndex * chunkSize;
        const end = Math.min(start + chunkSize, file.size);
        const blob = file.slice(start, end, file.type || 'application/octet-stream');
        const formData = new FormData();
        formData.append('file', blob, file.name);
        formData.append('chunkSize', chunkSize.toString());
        formData.append('chunkIndex', chunkIndex.toString());
        formData.append('totalChunks', totalChunks.toString());

        const { code, data, msg } = await uploadImage(formData, fileId);

        if (code !== 0) {
          console.log(msg);
          return;
        }

        fileId = data.id;
        chunkIndex++;
        setProgress((chunkIndex / totalChunks) * 100);
      }

      setImages((prev) => [...prev, file.name]);
    } catch (error: any) {
      notify.error({
        message: t('image.uploadError'),
        description: error.message,
        duration: 3
      });
    }
    setIsLoading(false);
    setProgress(0);
  }

  function selectFile(): void {
    if (inputRef.current) {
      inputRef.current.value = null;
    }
    inputRef.current?.click();
  }
  const menuContent = (
    <>
      <div className="flex flex-col gap-2">
        <div className="mt-6 flex items-center justify-between gap-2">
          <div className="ant-modal-title">{t('image.upload')}</div>
          <div className="flex items-center justify-end">
            <input
              ref={inputRef}
              type="file"
              accept=".img,.iso,.bin"
              className="hidden"
              onChange={uploadFile}
            />
            <Button
              ghost
              type="primary"
              size="small"
              icon={<UploadOutlined />}
              loading={isLoading}
              onClick={selectFile}
            ></Button>
          </div>
        </div>

        {progress !== 0 && (
          <div className="h-3 w-full rounded bg-slate-400/50">
            <div className="h-full rounded bg-blue-500" style={{ width: `${progress}%` }}></div>
          </div>
        )}

        <Divider />

        <div className="flex flex-col">
          <span className="text-sm">{t('image.uploaded')}:</span>
          <ul className="flex list-none flex-col">
            {images.length === 0 ? (
              <li>{t('image.empty')}</li>
            ) : (
              images.map((img) => (
                <li
                  key={img}
                  className={clsx(
                    'group flex max-w-[300px] cursor-pointer select-none items-center space-x-1',
                    'rounded p-2 hover:bg-neutral-700/70'
                  )}
                >
                  <div className="flex items-center justify-center gap-2 text-xs text-white/60">
                    <div className="h-[18px] w-[18px]">
                      <PackageIcon size={18} />
                    </div>
                    <span className="text-neutral-400">{img}</span>
                  </div>
                </li>
              ))
            )}
          </ul>
        </div>
      </div>
    </>
  );

  return (
    <>
      <div
        onClick={() => setIsModalOpen(true)}
        className="flex h-[30px] cursor-pointer select-none items-center space-x-2 rounded py-0.5 pl-3 pr-6 hover:bg-neutral-700/70"
      >
        <HardDriveUploadIcon size={16} className="text-neutral-400" />
        <span className="text-sm text-neutral-300">{t('image.upload')}</span>
      </div>

      <Modal
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        //
        footer={null}
        width={600}
      >
        {menuContent}
      </Modal>
    </>
  );
};
