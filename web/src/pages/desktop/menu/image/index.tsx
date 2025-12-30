import { DiscIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { MenuItem } from '@/components/menu-item';

import { DownloadImage } from './downloadImage';
import { MounteImage } from './mounteImage';
import { UploaderImages } from './uploader-images';

export const Image = () => {
  const { t } = useTranslation();

  const content = (
    <div className="flex flex-col space-y-1">
      <MounteImage />
      <UploaderImages />
      <DownloadImage />
    </div>
  );

  return (
    <MenuItem
      title={t('image.title')}
      icon={<DiscIcon size={18} />}
      content={content}
      fresh={true}
    />
  );
};
