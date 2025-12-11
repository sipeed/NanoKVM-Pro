import { useEffect } from 'react';
import { useAtom } from 'jotai';
import { useTranslation } from 'react-i18next';
import { useMediaQuery } from 'react-responsive';

import { VideoMode } from '@/types';
import * as storage from '@/lib/localstorage.ts';
import { getSupportedVideoModes } from '@/lib/video.ts';
import { client } from '@/lib/websocket.ts';
import { videoModeAtom } from '@/jotai/screen.ts';
import { Head } from '@/components/head.tsx';

import { Keyboard } from './keyboard';
import { Menu } from './menu';
import { Message } from './message.tsx';
import { Mouse } from './mouse';
import { Notification } from './notification.tsx';
import { Screen } from './screen';
import { VirtualKeyboard } from './virtual-keyboard';

export const Desktop = () => {
  const { t } = useTranslation();
  const isBigScreen = useMediaQuery({ minWidth: 850 });

  const [videoMode, setVideoMode] = useAtom(videoModeAtom);

  useEffect(() => {
    client.connect();

    const mode = getVideoMode() as VideoMode;
    setVideoMode(mode);

    return () => {
      client.close();
    };
  }, []);

  function getVideoMode() {
    const supportedModes = getSupportedVideoModes();

    const cookieVideoMode = storage.getVideoMode();
    if (cookieVideoMode && supportedModes.includes(cookieVideoMode)) {
      return cookieVideoMode;
    }

    return supportedModes.includes('h264-webrtc') ? 'h264-webrtc' : 'mjpeg';
  }

  return (
    <>
      <Head title={t('head.desktop')} />

      {isBigScreen && (
        <>
          <Message />
          <Notification />
        </>
      )}

      {videoMode && (
        <>
          <Menu />
          <Screen />
          <Mouse />
          <Keyboard />
        </>
      )}

      <VirtualKeyboard />
    </>
  );
};
