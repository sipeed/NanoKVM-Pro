import { useCallback, useEffect, useState } from 'react';
import { InputNumber, Slider, Switch } from 'antd';
import { SquareArrowOutUpRightIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { getLedConfig, setLedConfig } from '@/api/vm';

export const LedStrip = () => {
  const { t } = useTranslation();

  const [enabled, setEnabled] = useState(false);
  const [horizontal, setHorizontal] = useState(1);
  const [vertical, setVertical] = useState(1);
  const [brightness, setBrightness] = useState(0);
  const [needRestart, setNeedRestart] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    getLedConfig().then(({ data, code }) => {
      if (code !== 0) {
        return;
      }
      setEnabled(data.on);
      setBrightness(data.brightness);
      setHorizontal(data.hor);
      setVertical(data.ver);
      setNeedRestart(false);
    });
  }, []);

  function update(on: boolean, hor: number, ver: number, bright: number) {
    if (isLoading) return;
    setIsLoading(true);

    setLedConfig({
      on: on,
      hor: hor,
      ver: ver,
      brightness: bright
    })
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
          return;
        }

        setEnabled(on);
        setHorizontal(hor);
        setVertical(ver);
        setBrightness(bright);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  const horMax = useCallback(() => {
    return 150 - vertical * 2;
  }, [vertical]);

  const verMax = useCallback(() => {
    return Math.floor((150 - horizontal) / 2);
  }, [horizontal]);

  function openDoc() {
    window.open('https://wiki.sipeed.com/hardware/en/kvm/NanoKVM_Pro/ws2812.html', '_blank');
  }

  return (
    <div className="">
      <div className="flex items-center justify-between">
        <div className="flex flex-col">
          <div className="flex items-center space-x-2">
            <span>{t('settings.device.led.title')}</span>
            <div className="cursor-pointer text-neutral-500 hover:text-blue-500" onClick={openDoc}>
              <SquareArrowOutUpRightIcon size={14} />
            </div>
          </div>
          <span className="text-xs text-neutral-500">{t('settings.device.led.description')}</span>
        </div>

        <Switch
          checked={enabled}
          loading={isLoading}
          onChange={(value) => {
            update(value, horizontal, vertical, brightness);
            setNeedRestart(false);
          }}
        />
      </div>

      {enabled && (
        <div className="ml-5 mt-3 space-y-1 rounded bg-neutral-800/20 p-3">
          <div className="flex w-full items-center justify-between gap-2">
            <div className="flex items-center justify-around gap-2">
              <span className="text-neutral-300">{t('settings.device.led.horizontal')}</span>
            </div>
            <InputNumber
              max={horMax()}
              min={1}
              defaultValue={horizontal}
              changeOnWheel
              onChange={(value) => {
                update(enabled, value || 0, vertical, brightness);
                setNeedRestart(true);
              }}
            />
          </div>

          <div className="flex w-full items-center justify-between gap-2">
            <div className="flex items-center justify-around gap-2">
              <span className="text-neutral-300">{t('settings.device.led.vertical')}</span>
            </div>
            <InputNumber
              max={verMax()}
              min={1}
              defaultValue={vertical}
              onChange={(value) => {
                update(enabled, horizontal, value || 0, brightness);
                setNeedRestart(true);
              }}
            />
          </div>

          <div className="flex w-full items-center justify-between gap-2">
            <div className="flex items-center justify-around gap-2">
              <span className="text-neutral-300">{t('settings.device.led.brightness')}</span>
            </div>
            <div className="min-w-[200px]">
              <Slider
                value={brightness}
                range={false}
                onChange={(value) => {
                  update(enabled, horizontal, vertical, value || 0);
                }}
              />
            </div>
          </div>

          {needRestart && (
            <div className="flex justify-center pt-1 text-xs text-red-500">
              {t('settings.device.led.restart')}
            </div>
          )}
        </div>
      )}
    </div>
  );
};
