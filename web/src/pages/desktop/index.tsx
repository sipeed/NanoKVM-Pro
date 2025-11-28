import { useEffect } from 'react';
import { useAtom, useAtomValue } from 'jotai';
import { useTranslation } from 'react-i18next';
import { useMediaQuery } from 'react-responsive';

import { VideoMode } from '@/types';
import * as storage from '@/lib/localstorage.ts';
import { getSupportedVideoModes } from '@/lib/video.ts';
import { client } from '@/lib/websocket.ts';
import { isKeyboardEnableAtom } from '@/jotai/keyboard.ts';
import { videoModeAtom } from '@/jotai/screen.ts';
import { Head } from '@/components/head.tsx';

import { Keyboard } from './keyboard';
import { VirtualKeyboard } from './keyboard/virtual-keyboard';
import { Menu } from './menu';
import { Message } from './message.tsx';
import { Mouse } from './mouse';
import { Notification } from './notification.tsx';
import { Screen } from './screen';

export const Desktop = () => {
  const { t } = useTranslation();
  const isBigScreen = useMediaQuery({ minWidth: 850 });

  const [videoMode, setVideoMode] = useAtom(videoModeAtom);
  const isKeyboardEnable = useAtomValue(isKeyboardEnableAtom);

  useEffect(() => {
    const mode = getVideoMode() as VideoMode;
    setVideoMode(mode);

    const timer = setInterval(() => {
      client.send([0]);
    }, 10 * 1000);

    return () => {
      clearInterval(timer);
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
          {isKeyboardEnable && <Keyboard />}
        </>
      )}

      <VirtualKeyboard />
    </>
  );
};
