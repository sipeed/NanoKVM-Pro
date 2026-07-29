import { useEffect, useState } from 'react';
import { Switch, TimePicker } from 'antd';
import dayjs, { type Dayjs } from 'dayjs';
import { useTranslation } from 'react-i18next';

import { getLCDScreenOff, setLCDScreenOff, type LCDScreenOffConfig } from '@/api/vm.ts';

type ScreenOffSettings = Omit<LCDScreenOffConfig, 'supported'>;

function minuteToTime(minute: number): Dayjs {
  return dayjs().startOf('day').add(minute, 'minute');
}

function timeToMinute(time: Dayjs): number {
  return time.hour() * 60 + time.minute();
}

export const ScheduledScreenOff = () => {
  const { t } = useTranslation();
  const [savedSettings, setSavedSettings] = useState<ScreenOffSettings>();
  const [settings, setSettings] = useState<ScreenOffSettings>();
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [validationError, setValidationError] = useState(false);
  const [saveError, setSaveError] = useState(false);

  useEffect(() => {
    let cancelled = false;

    getLCDScreenOff()
      .then((rsp) => {
        if (cancelled || rsp.code !== 0 || !rsp.data.supported) {
          return;
        }

        const nextSettings: ScreenOffSettings = {
          enabled: rsp.data.enabled,
          startMinute: rsp.data.startMinute,
          endMinute: rsp.data.endMinute
        };
        setSavedSettings(nextSettings);
        setSettings(nextSettings);
      })
      .catch((error) => {
        if (!cancelled) {
          console.error(error);
        }
      })
      .finally(() => {
        if (!cancelled) {
          setIsLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  function update(nextSettings: ScreenOffSettings) {
    if (isSaving) return;

    setValidationError(false);
    setSaveError(false);
    setSettings(nextSettings);

    if (nextSettings.startMinute === nextSettings.endMinute) {
      setValidationError(true);
      return;
    }

    setIsSaving(true);
    setLCDScreenOff(nextSettings)
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.error(rsp.msg);
          setSaveError(true);
          if (savedSettings) {
            setSettings(savedSettings);
          }
          return;
        }

        setSavedSettings(nextSettings);
      })
      .catch((error) => {
        console.error(error);
        setSaveError(true);
        if (savedSettings) {
          setSettings(savedSettings);
        }
      })
      .finally(() => {
        setIsSaving(false);
      });
  }

  if (isLoading || !settings) {
    return <></>;
  }

  const isBusy = isLoading || isSaving;

  return (
    <div className="space-y-3 pt-3">
      <div className="flex items-center justify-between">
        <div className="flex flex-col space-y-1">
          <span>{t('settings.device.scheduledScreenOff.title')}</span>
          <span className="text-xs text-neutral-500">
            {t('settings.device.scheduledScreenOff.description')}
          </span>
        </div>
        <Switch
          checked={settings.enabled}
          disabled={isBusy}
          loading={isSaving}
          onChange={(enabled) => update({ ...settings, enabled })}
        />
      </div>

      <div className="mx-5 space-y-2 rounded bg-neutral-800/20 p-3">
        <div className="flex items-center justify-between gap-3">
          <span className="text-sm text-neutral-300">
            {t('settings.device.scheduledScreenOff.start')}
          </span>
          <TimePicker
            allowClear={false}
            disabled={!settings.enabled || isBusy}
            format="HH:mm"
            minuteStep={1}
            value={minuteToTime(settings.startMinute)}
            onChange={(time) => {
              if (time) {
                update({ ...settings, startMinute: timeToMinute(time) });
              }
            }}
          />
        </div>
        <div className="flex items-center justify-between gap-3">
          <span className="text-sm text-neutral-300">
            {t('settings.device.scheduledScreenOff.end')}
          </span>
          <TimePicker
            allowClear={false}
            disabled={!settings.enabled || isBusy}
            format="HH:mm"
            minuteStep={1}
            value={minuteToTime(settings.endMinute)}
            onChange={(time) => {
              if (time) {
                update({ ...settings, endMinute: timeToMinute(time) });
              }
            }}
          />
        </div>
        {validationError && (
          <div className="text-xs text-red-500">
            {t('settings.device.scheduledScreenOff.invalidRange')}
          </div>
        )}
        {saveError && (
          <div className="text-xs text-red-500">
            {t('settings.device.scheduledScreenOff.saveFailed')}
          </div>
        )}
      </div>
    </div>
  );
};
