import { Popover, Slider } from 'antd';
import { useAtom } from 'jotai';
import { CircleHelpIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/stream.ts';
import * as storage from '@/lib/localstorage.ts';
import { videoParametersAtom } from '@/jotai/screen.ts';

export const Fps = () => {
  const { t } = useTranslation();
  const [videoParameters, setVideoParameters] = useAtom(videoParametersAtom);

  async function update(value: number) {
    if (value === videoParameters.fps) {
      return;
    }

    try {
      const rsp = await api.setFps(value);
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      const parameters = { ...videoParameters, fps: value };
      setVideoParameters(parameters);
      storage.setVideoParameters(JSON.stringify(parameters));
    } catch (err) {
      console.log(err);
    }
  }

  const help = (
    <ul className="max-w-[550px] list-outside">
      <li>{t('settings.screen.fps.tip1')}</li>
      <li>{t('settings.screen.fps.tip2')}</li>
    </ul>
  );

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-1 flex-col space-y-1">
        <div className="flex items-center space-x-2">
          <span>FPS</span>

          <Popover className="cursor-pointer text-neutral-500" content={help}>
            <CircleHelpIcon size={15} />
          </Popover>
        </div>

        <span className="text-xs text-neutral-500">{t('settings.screen.fps.description')}</span>
      </div>

      <div className="w-[240px]">
        <Slider
          value={videoParameters.fps}
          min={0}
          max={120}
          marks={{
            0: videoParameters.fps === 0 ? t('settings.screen.fps.auto') : '0',
            120: '120'
          }}
          onChange={update}
        />
      </div>
    </div>
  );
};
