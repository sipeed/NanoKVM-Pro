import { useEffect, useState } from 'react';
import { Select } from 'antd';
import { ExternalLinkIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/vm.ts';

const defaultEdidList = [
  { value: 'E18-4K30FPS', label: '3840 x 2160 30Hz' },
  { value: 'E48-4K39FPS', label: '3840 x 2160 39Hz' },
  { value: 'E56-2K60FPS', label: '2560 x 1440 60Hz' },
  { value: 'E54-1080P60FPS', label: '1920 x 1080 60Hz' },
  { value: 'E58-4K16-10', label: '3840 x 2400 30Hz' },
  { value: 'E63-Ultrawide', label: '3440 x 1440 60Hz' }
];

export const Edid = () => {
  const { t } = useTranslation();

  const [currentEdid, setCurrentEdid] = useState('');
  const [edidList, setEdidList] = useState<any[]>(defaultEdidList);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    getCurrentEdid();
    getCustomEdidList();
  }, []);

  function getCurrentEdid() {
    api.getEdid().then((rsp) => {
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      setCurrentEdid(rsp.data.edid);
    });
  }

  function getCustomEdidList() {
    api.getCustomEdidList().then((rsp) => {
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      if (rsp.data?.edidList?.length > 0) {
        setEdidList([
          {
            label: <span>{t('settings.screen.edid.default')}</span>,
            title: 'default',
            options: defaultEdidList
          },
          {
            label: <span>{t('settings.screen.edid.custom')}</span>,
            title: 'custom',
            options: rsp.data.edidList.map((edid: string) => {
              return {
                value: edid,
                label: edid
              };
            })
          }
        ]);
      }
    });
  }

  function update(value: string) {
    if (value === currentEdid || isLoading) return;

    setIsLoading(true);

    api
      .updateEdid(value)
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
          return;
        }
        setCurrentEdid(value);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  return (
    <div className="flex items-center justify-between space-x-5">
      <div className="flex flex-col space-y-1">
        <div className="flex items-center space-x-2">
          <span>EDID</span>

          <a
            className="flex items-center text-neutral-500 hover:text-blue-500"
            href="https://wiki.sipeed.com/hardware/en/kvm/NanoKVM_Pro/extended.html#How-to-Modify-EDID"
            target="_blank"
          >
            <ExternalLinkIcon size={15} />
          </a>
        </div>

        <span className="text-xs text-neutral-500">{t('settings.screen.edid.description')}</span>
      </div>

      <Select
        value={currentEdid}
        style={{ width: 240 }}
        loading={isLoading}
        options={edidList}
        onSelect={update}
      />
    </div>
  );
};
