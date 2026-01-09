import { Divider } from 'antd';
import { useTranslation } from 'react-i18next';

import { Advanced } from './advanced';
import { Datetime } from './datetime';
import { HdmiCapture } from './hdmi-capture.tsx';
import { HdmiPassthrough } from './hdmi-passthrough.tsx';
import { LedStrip } from './led-strip.tsx';
import { Mdns } from './mdns.tsx';
import { MouseJiggler } from './mouse-jiggler.tsx';
import { Oled } from './oled.tsx';
import { Reboot } from './reboot.tsx';
import { Ssh } from './ssh.tsx';
import { VirtualDevices } from './virtual-devices';
import { Wifi } from './wifi.tsx';

export const Device = () => {
  const { t } = useTranslation();

  return (
    <>
      <div className="text-base font-bold">{t('settings.device.title')}</div>
      <Divider className="opacity-50" />

      <div className="flex flex-col space-y-8">
        <Ssh />
        <Mdns />
        <HdmiCapture />
        <HdmiPassthrough />
        <Divider className="opacity-50" />

        <VirtualDevices />
        <Divider className="opacity-50" />

        <Oled />
        <Wifi />
        <MouseJiggler />
        <LedStrip />
        <Divider className="opacity-50" />

        <Datetime />
        <Divider className="opacity-50" />

        <Advanced />
      </div>

      <Divider className="opacity-50" />

      <Reboot />
    </>
  );
};
