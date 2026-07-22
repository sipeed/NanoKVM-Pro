import { useEffect } from 'react';
import { useAtom, useAtomValue } from 'jotai';
import { MonitorIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/stream.ts';
import { VideoParameters } from '@/types';
import * as storage from '@/lib/localstorage.ts';
import { videoModeAtom, videoParametersAtom } from '@/jotai/screen.ts';
import { MenuItem } from '@/components/menu-item.tsx';

import { Advanced } from './advanced.tsx';
import { Bitrate } from './bitrate.tsx';
import { Edid } from './edid.tsx';
import { Quality } from './quality';
import { Rotation } from './rotation.tsx';
import { Scale } from './scale.tsx';
import { VideoMode } from './video-mode.tsx';

export const Screen = () => {
  const { t } = useTranslation();

  const videoMode = useAtomValue(videoModeAtom);
  const [videoParameters, setVideoParameters] = useAtom(videoParametersAtom);

  useEffect(() => {
    if (!videoMode) return;

    api.setMode(videoMode);
    initScreen();
  }, [videoMode]);

  async function initScreen() {
    const parameters = getVideoParameters();
    let { rateControlMode, bitrate, gop, fps, quality } = parameters;

    const scale = getScale(parameters.scale);
    const rotation = getRotation(parameters.rotation);
    setVideoParameters({ ...parameters, scale, rotation });

    try {
      if (videoMode === 'mjpeg') {
        quality = await updateQuality(parameters.quality);
      } else {
        rateControlMode = await updateRateControlMode(parameters.rateControlMode);
        bitrate = await updateBitrate(rateControlMode, parameters.bitrate);
        gop = await updateGop(parameters.gop);
      }

      fps = await updateFps(parameters.fps);

      setVideoParameters((current) => ({
        ...current,
        rateControlMode,
        bitrate,
        gop,
        fps,
        quality
      }));
    } catch (err) {
      console.log(err);
    }
  }

  function getVideoParameters() {
    const parameters = storage.getVideoParameters();
    if (!parameters) {
      return videoParameters;
    }
    return JSON.parse(parameters) as VideoParameters;
  }

  async function updateRateControlMode(mode: string) {
    if (!mode || !['vbr', 'cbr'].includes(mode)) {
      mode = 'vbr';
    }

    const rsp = await api.setRateControlMode(mode);
    if (rsp.code !== 0) {
      return 'vbr';
    }

    return mode;
  }

  async function updateBitrate(mode: string, bitrate: number) {
    if (bitrate < 1000 || bitrate > 20000) {
      bitrate = mode === 'vbr' ? 8000 : 5000;
    }

    const rsp = await api.setQuality(bitrate);
    if (rsp.code !== 0) {
      return 5000;
    }

    return bitrate;
  }

  async function updateFps(fps: number) {
    if (!fps || fps < 0 || fps > 120) {
      fps = 0;
    }

    const rsp = await api.setFps(fps);
    if (rsp.code !== 0) {
      return 0;
    }

    return fps;
  }

  async function updateGop(gop: number) {
    if (gop < 1 || gop > 200) {
      gop = 50;
    }

    const rsp = await api.setGop(gop);
    if (rsp.code !== 0) {
      return 50;
    }

    return gop;
  }

  async function updateQuality(quality: number | undefined) {
    if (!quality || quality < 1 || quality > 100) {
      quality = 80;
    }

    const rsp = await api.setQuality(quality);
    if (rsp.code !== 0) {
      return 80;
    }

    return quality;
  }

  function getScale(scale: number) {
    if (scale < 0.5 || scale > 2) {
      return 1;
    }
    return scale;
  }

  function getRotation(rotation: VideoParameters['rotation']) {
    if (![0, 90, 180, 270].includes(rotation)) {
      return 0;
    }
    return rotation;
  }

  const content = (
    <div className="flex flex-col space-y-1">
      <VideoMode />
      <Edid />
      {videoMode === 'mjpeg' ? <Quality /> : <Bitrate />}
      <Scale />
      <Rotation />
      <Advanced />
    </div>
  );

  return <MenuItem title={t('screen.title')} icon={<MonitorIcon size={18} />} content={content} />;
};
