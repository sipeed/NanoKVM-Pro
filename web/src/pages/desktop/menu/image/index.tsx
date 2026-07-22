import { useEffect, useState } from 'react';
import clsx from 'clsx';
import { DiscIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/storage.ts';
import { MenuItem } from '@/components/menu-item';

import { DownloadImage } from './downloadImage';
import { MounteImage } from './mounteImage';
import { UploaderImages } from './uploader-images';

export const Image = () => {
  const { t } = useTranslation();
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    api.getMountedImage().then((rsp) => {
      if (rsp.code === 0) {
        setIsMounted(!!rsp.data?.file);
      }
    });
  }, []);

  const icon = (
    <div
      className={clsx(
        'h-[18px] w-[18px]',
        isMounted ? 'text-blue-500' : 'text-neutral-300 hover:text-white'
      )}
    >
      <DiscIcon size={18} />
    </div>
  );

  const content = (
    <div className="flex flex-col space-y-1">
      <MounteImage isMounted={isMounted} setIsMounted={setIsMounted} />
      <UploaderImages />
      <DownloadImage />
    </div>
  );

  return (
    <MenuItem
      title={t('image.title')}
      icon={icon}
      content={content}
      fresh={true}
    />
  );
};
