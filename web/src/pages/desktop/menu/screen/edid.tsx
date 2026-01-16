import { ChangeEvent, useEffect, useRef, useState } from 'react';
import { Divider, Popover } from 'antd';
import {
  ArrowUpFromLineIcon,
  BookIcon,
  CheckIcon,
  LoaderCircleIcon,
  RatioIcon,
  Trash2Icon
} from 'lucide-react';
import { useTranslation } from 'react-i18next';

import * as api from '@/api/vm.ts';

export const Edid = () => {
  const { t } = useTranslation();

  const [currentEdid, setCurrentEdid] = useState('');
  const [customEdidList, setCustomEdidList] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isUploading, setIsUploading] = useState(false);

  const inputRef = useRef<any>(null);

  const edidList = [
    { key: 'E18-4K30FPS', width: 3840, height: 2160, rate: '30Hz' },
    { key: 'E48-4K39FPS', width: 3840, height: 2160, rate: '39Hz' },
    { key: 'E56-2K60FPS', width: 2560, height: 1440, rate: '60Hz' },
    { key: 'E54-1080P60FPS', width: 1920, height: 1080, rate: '60Hz' },
    { key: 'E58-4K16-10', width: 3840, height: 2400, rate: '30Hz' },
    { key: 'E63-Ultrawide', width: 3440, height: 1440, rate: '60Hz' }
  ];

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
        setCustomEdidList(rsp.data.edidList);
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

  function selectFile() {
    if (inputRef.current) {
      inputRef.current.value = null;
    }

    inputRef.current?.click();
  }

  function uploadFile(e: ChangeEvent<HTMLInputElement>) {
    if (!e.target?.files?.length) return;
    const file = e.target.files[0];

    if (isUploading) return;
    setIsUploading(true);

    const formData = new FormData();
    formData.append('file', file);

    api
      .uploadEdid(formData)
      .then((rsp) => {
        if (rsp.code !== 0) {
          console.log(rsp.msg);
          return;
        }

        if (!customEdidList.includes(rsp.data.file)) {
          setCustomEdidList([...customEdidList, rsp.data.file]);
        }
      })
      .finally(() => {
        setIsUploading(false);
      });
  }

  function deleteEdid(e: any, edid: string) {
    e.stopPropagation();

    api.deleteEdid(edid).then((rsp) => {
      if (rsp.code !== 0) {
        console.log(rsp.msg);
        return;
      }

      setCustomEdidList(customEdidList.filter((customEdid) => customEdid !== edid));
    });
  }

  const content = (
    <>
      {/* default EDID list */}
      {edidList.map((edid) => (
        <div
          key={edid.key}
          className="flex cursor-pointer select-none items-center rounded py-1.5 pl-1 pr-6 hover:bg-neutral-700/70"
          onClick={() => update(edid.key)}
        >
          <div className="flex h-[14px] w-[20px] items-end text-blue-500">
            {edid.key === currentEdid && <CheckIcon size={14} />}
          </div>

          <div className="flex w-[80px] items-center justify-between">
            <span>{edid.width}</span>
            <span className="text-xs">x</span>
            <span>{edid.height}</span>
          </div>
          <div className="pl-3">{edid.rate}</div>
        </div>
      ))}

      <Divider style={{ margin: '5px 0' }} />

      {customEdidList.length > 0 ? (
        // custom EDID list
        <>
          {customEdidList.map((edid) => (
            <div
              key={edid}
              className="group flex cursor-pointer select-none items-center justify-between rounded py-1.5 pl-1 pr-6 hover:bg-neutral-700/70"
              onClick={() => update(edid)}
            >
              <div className="flex items-center">
                <div className="flex h-[14px] w-[20px] items-end text-blue-500">
                  {edid === currentEdid && <CheckIcon size={14} />}
                </div>
                <div className="max-w-[100px] select-none truncate">{edid.replace('.bin', '')}</div>
              </div>

              <div
                className="hidden h-[16px] w-[16px] items-center text-red-500 hover:text-red-500/60 group-hover:flex"
                onClick={(e) => deleteEdid(e, edid)}
              >
                <Trash2Icon size={14} />
              </div>
            </div>
          ))}
        </>
      ) : (
        <a
          className="flex cursor-pointer select-none items-center rounded px-1 py-1.5 text-neutral-300 hover:bg-neutral-700/70 hover:text-neutral-300"
          href="https://wiki.sipeed.com/hardware/en/kvm/NanoKVM_Pro/extended.html#How-to-Modify-EDID"
          target="_blank"
        >
          <div className="h-[14px] w-[20px]" />
          <div className="flex items-center space-x-1">
            <BookIcon size={14} />
            <span>{t('settings.screen.edid.document')}</span>
          </div>
        </a>
      )}

      <Divider style={{ margin: '5px 0' }} />

      <input ref={inputRef} type="file" accept=".bin" className="hidden" onChange={uploadFile} />

      <div
        className="flex cursor-pointer select-none items-center rounded px-1 py-1.5 hover:bg-neutral-700/70"
        onClick={selectFile}
      >
        <div className="h-[14px] w-[20px]" />
        <div className="flex items-center space-x-1">
          {isUploading ? (
            <LoaderCircleIcon className="animate-spin" size={14} />
          ) : (
            <ArrowUpFromLineIcon size={14} />
          )}
          <span>{t('screen.upload')}</span>
        </div>
      </div>
    </>
  );

  return (
    <Popover content={content} placement="rightTop" arrow={false} align={{ offset: [14, 0] }}>
      <div className="flex h-[30px] cursor-pointer items-center space-x-2 rounded pl-3 pr-6 text-neutral-300 hover:bg-neutral-700/70">
        <RatioIcon size={18} />
        <span className="select-none text-sm">EDID</span>
      </div>
    </Popover>
  );
};
