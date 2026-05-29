import { useEffect, useRef, useState } from 'react';
import { Button, Modal } from 'antd';
import clsx from 'clsx';
import Cookies from 'js-cookie';
import { CheckIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { VideoMode } from '@/types';
import { getSupportedVideoModes } from '@/lib/video.ts';

type VideoModeH265Props = {
  videoMode: VideoMode | null;
  update: (mode: string) => void;
};

const videoModes = [
  { key: 'h265-webrtc', name: 'H.265 WebRTC' },
  { key: 'h265-direct', name: 'H.265 Direct' }
];

export const VideoModeH265 = ({ videoMode, update }: VideoModeH265Props) => {
  const { t } = useTranslation();

  const [isSupported, setIsSupported] = useState(false);
  const [isEnabled, setIsEnabled] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const divRef = useRef<any>(null);

  useEffect(() => {
    const videoModes = getSupportedVideoModes();
    const supported = videoModes.some((mode) => mode.startsWith('h265'));
    setIsSupported(supported);

    const h265 = Cookies.get('h265');
    if (h265) {
      setIsEnabled(true);
      return;
    }

    const targetNode = divRef.current;
    if (!targetNode) return;

    const observer = new MutationObserver((mutationsList, obs) => {
      for (const mutation of mutationsList) {
        if (mutation.type === 'childList') {
          mutation.removedNodes.forEach((removedNode) => {
            if (removedNode === targetNode) {
              Cookies.set('h265', 'Enable');
              setIsEnabled(true);
              obs.disconnect();
            }
          });
        }
      }
    });

    const config = { childList: true };
    if (targetNode.parentNode) {
      observer.observe(targetNode.parentNode, config);
    }

    return () => {
      observer.disconnect();
    };
  }, []);

  function handleClick(mode: string) {
    if (!isSupported || !isEnabled) {
      setIsModalOpen(true);
      return;
    }

    update(mode);
  }

  return (
    <>
      <div id="disable-H.265" ref={divRef} />
      {videoModes.map((mode) => (
        <div
          key={mode.key}
          className={clsx(
            'flex select-none items-center rounded py-1.5 pl-1 pr-5 hover:bg-neutral-700/70',
            isSupported && isEnabled ? 'cursor-pointer' : 'text-neutral-500 opacity-70'
          )}
          onClick={() => handleClick(mode.key)}
        >
          <div className="flex h-[14px] w-[20px] items-end text-blue-500">
            {mode.key === videoMode && <CheckIcon size={15} />}
          </div>
          <span>{mode.name}</span>
        </div>
      ))}

      <Modal
        centered
        title={t('notification.h265.title')}
        open={isModalOpen}
        footer={null}
        closeIcon={false}
      >
        <div className="whitespace-pre-line py-5">
          {isSupported ? t('notification.h265.description') : t('notification.h265.notSupported')}
        </div>

        <div className="flex justify-center py-3">
          <Button
            type="primary"
            size="large"
            className="w-32"
            onClick={() => setIsModalOpen(false)}
          >
            {t('notification.btnOk')}
          </Button>
        </div>
      </Modal>
    </>
  );
};
