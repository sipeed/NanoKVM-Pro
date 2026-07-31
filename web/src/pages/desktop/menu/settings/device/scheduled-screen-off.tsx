import { cloneElement, isValidElement, useEffect, useState, type MouseEvent } from 'react';
import { ConfigProvider, Switch, TimePicker } from 'antd';
import dayjs, { type Dayjs } from 'dayjs';
import { useTranslation } from 'react-i18next';

import { getLCDScreenOff, setLCDScreenOff, type LCDScreenOffConfig } from '@/api/vm.ts';

type ScreenOffSettings = Omit<LCDScreenOffConfig, 'supported'>;
type TimeSelectionStep = 'hour' | 'minute';

const TIME_PICKER_WIDTH = 128;

function minuteToTime(minute: number): Dayjs {
  return dayjs().startOf('day').add(minute, 'minute');
}

function timeToMinute(time: Dayjs): number {
  return time.hour() * 60 + time.minute();
}

interface ScheduledTimePickerProps {
  disabled: boolean;
  minute: number;
  onChange: (minute: number) => void;
}

const ScheduledTimePicker = ({ disabled, minute, onChange }: ScheduledTimePickerProps) => {
  const [open, setOpen] = useState(false);
  const [step, setStep] = useState<TimeSelectionStep>('hour');
  const [draftTime, setDraftTime] = useState(() => minuteToTime(minute));

  useEffect(() => {
    if (!open) {
      setDraftTime(minuteToTime(minute));
    }
  }, [minute, open]);

  function reset() {
    setOpen(false);
    setStep('hour');
    setDraftTime(minuteToTime(minute));
  }

  function selectUnit(unit: number) {
    if (step === 'hour') {
      setDraftTime(draftTime.hour(unit));
      setStep('minute');
      return;
    }

    const nextTime = draftTime.minute(unit);
    setDraftTime(nextTime);
    setStep('hour');
    setOpen(false);
    onChange(timeToMinute(nextTime));
  }

  return (
    <ConfigProvider
      theme={{
        components: {
          DatePicker: {
            timeColumnWidth: TIME_PICKER_WIDTH
          }
        }
      }}
    >
      <TimePicker
        allowClear={false}
        disabled={disabled}
        format="HH:mm"
        inputReadOnly
        minuteStep={1}
        needConfirm={false}
        open={open}
        showHour={step === 'hour'}
        showMinute={step === 'minute'}
        showNow={false}
        showSecond={false}
        style={{ width: TIME_PICKER_WIDTH }}
        styles={{ popup: { root: { width: TIME_PICKER_WIDTH } } }}
        value={open ? draftTime : minuteToTime(minute)}
        cellRender={(current, info) => {
          if (
            info.type !== 'time' ||
            !isValidElement<{ onClick?: (event: MouseEvent<HTMLElement>) => void }>(info.originNode)
          ) {
            return info.originNode;
          }

          return cloneElement(info.originNode, {
            onClick: (event: MouseEvent<HTMLElement>) => {
              event.stopPropagation();
              selectUnit(Number(current));
            }
          });
        }}
        onCalendarChange={(time) => {
          if (time && !Array.isArray(time)) {
            selectUnit(step === 'hour' ? time.hour() : time.minute());
          }
        }}
        onOpenChange={(nextOpen) => {
          if (nextOpen) {
            setDraftTime(minuteToTime(minute));
            setStep('hour');
            setOpen(true);
          } else {
            reset();
          }
        }}
      />
    </ConfigProvider>
  );
};

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

    if (nextSettings.enabled && nextSettings.startMinute === nextSettings.endMinute) {
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

      {settings.enabled && (
        <div className="mx-5 space-y-2 rounded bg-neutral-800/20 p-3">
          <div className="flex items-center justify-between gap-3">
            <span className="text-sm text-neutral-300">
              {t('settings.device.scheduledScreenOff.start')}
            </span>
            <ScheduledTimePicker
              disabled={isBusy}
              minute={settings.startMinute}
              onChange={(startMinute) => update({ ...settings, startMinute })}
            />
          </div>
          <div className="flex items-center justify-between gap-3">
            <span className="text-sm text-neutral-300">
              {t('settings.device.scheduledScreenOff.end')}
            </span>
            <ScheduledTimePicker
              disabled={isBusy}
              minute={settings.endMinute}
              onChange={(endMinute) => update({ ...settings, endMinute })}
            />
          </div>
          {validationError && (
            <div className="text-xs text-red-500">
              {t('settings.device.scheduledScreenOff.invalidRange')}
            </div>
          )}
        </div>
      )}
      {saveError && (
        <div className="mx-5 text-xs text-red-500">
          {t('settings.device.scheduledScreenOff.saveFailed')}
        </div>
      )}
    </div>
  );
};
