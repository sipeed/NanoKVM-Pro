import { useEffect, useState } from 'react';
import { Select } from 'antd';
import { useTranslation } from 'react-i18next';

import { getLCDDisplayPolicy, setLCDDisplayPolicy, type LCDDisplayMode, type LCDDisplayPolicy } from '@/api/vm.ts';

export const DisplayPolicy = () => {
  const { t } = useTranslation();
  const [policy, setPolicy] = useState<LCDDisplayPolicy>();
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState(false);

  useEffect(() => {
    getLCDDisplayPolicy().then((rsp) => {
      if (rsp.code === 0 && rsp.data?.supportedModes?.length) setPolicy(rsp.data);
    }).catch(console.error);
  }, []);

  if (!policy) return null;
  const mode = policy.mode ?? 'alwaysOn';
  const save = (update: Parameters<typeof setLCDDisplayPolicy>[0]) => {
    if (saving) return;
    setSaveError(false);
    setSaving(true);
    setLCDDisplayPolicy(update).then((rsp) => {
      if (rsp.code === 0) setPolicy((current) => current && {
        ...current,
        ...(update.mode ? { mode: update.mode } : {}),
        ...(update.schedule ? { schedule: { ...current.schedule, ...update.schedule } } : {})
      });
      else setSaveError(true);
    }).catch((error) => {
      console.error(error);
      setSaveError(true);
    }).finally(() => setSaving(false));
  };

  return <div className="space-y-3 pt-3">
    <div className="flex items-center justify-between gap-3">
      <div className="flex flex-col space-y-1">
        <span>{t('settings.device.displayPolicy.title')}</span>
        <span className="text-xs text-neutral-500">{t('settings.device.displayPolicy.description')}</span>
      </div>
      <Select disabled={saving} value={mode} style={{ width: 150 }} onChange={(value: LCDDisplayMode) => save({ mode: value })}
        options={policy.supportedModes.map((value) => ({ value, label: t(`settings.device.displayPolicy.modes.${value}`) }))} />
    </div>
    {saveError && <div className="mx-5 text-xs text-red-500">{t('settings.device.displayPolicy.saveFailed')}</div>}
  </div>;
};
