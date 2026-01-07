import { Slider } from 'antd';
import { useAtom } from 'jotai';
import { Volume1Icon, Volume2Icon, VolumeIcon, VolumeXIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { videoVolumeAtom } from '@/jotai/screen.ts';
import { MenuItem } from '@/components/menu-item.tsx';

export const Volume = () => {
  const { t } = useTranslation();

  const [videoVolume, setVideoVolume] = useAtom(videoVolumeAtom);

  const icon = (
    <>
      {videoVolume === 0 ? (
        <VolumeXIcon size={18} />
      ) : videoVolume < 50 ? (
        <Volume1Icon size={18} />
      ) : (
        <Volume2Icon size={18} />
      )}
    </>
  );

  const content = (
    <div className="flex items-center">
      <div
        className="mr-0.5 flex h-[16px] w-[16px] cursor-pointer items-center justify-center"
        onClick={() => setVideoVolume(0)}
      >
        <VolumeIcon size={18} />
      </div>

      <div className="min-w-[230px]">
        <Slider value={videoVolume} range={false} onChange={setVideoVolume} />
      </div>

      <div
        className="ml-2 flex h-[16px] w-[16px] cursor-pointer items-center justify-center"
        onClick={() => setVideoVolume(100)}
      >
        <Volume2Icon size={18} />
      </div>
    </div>
  );

  return <MenuItem title={t('volume.title')} icon={icon} content={content} />;
};
