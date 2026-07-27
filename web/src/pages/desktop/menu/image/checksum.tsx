import { useEffect, useRef, useState } from 'react';
import { Button, Popover, Tooltip } from 'antd';
import { HashIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { checksumImage } from '@/api/storage';
import type { ChecksumAlgorithm } from '@/api/storage';

const algorithms: { label: string; value: ChecksumAlgorithm }[] = [
  { label: 'MD5', value: 'md5' },
  { label: 'SHA-1', value: 'sha1' },
  { label: 'SHA-256', value: 'sha256' }
];

type ChecksumProps = {
  image: string;
  refreshKey: number;
  onError: () => void;
};

export const Checksum = ({ image, refreshKey, onError }: ChecksumProps) => {
  const { t } = useTranslation();
  const [checksums, setChecksums] = useState<Partial<Record<ChecksumAlgorithm, string>>>({});
  const [loading, setLoading] = useState<ChecksumAlgorithm>();
  const requestGeneration = useRef(0);

  useEffect(() => {
    requestGeneration.current += 1;
    setChecksums({});
    setLoading(undefined);

    return () => {
      requestGeneration.current += 1;
    };
  }, [image, refreshKey]);

  function calculate(algorithm: ChecksumAlgorithm) {
    if (loading) return;
    const generation = ++requestGeneration.current;
    setLoading(algorithm);

    checksumImage(image, algorithm)
      .then((rsp) => {
        if (generation !== requestGeneration.current) return;

        const checksum = rsp.data?.checksum;
        if (rsp.code !== 0 || !checksum) {
          onError();
          return;
        }

        setChecksums((current) => ({ ...current, [algorithm]: checksum }));
      })
      .catch(() => {
        if (generation === requestGeneration.current) onError();
      })
      .finally(() => {
        if (generation === requestGeneration.current) setLoading(undefined);
      });
  }

  const content = (
    <div className="flex items-center space-x-1" onClick={(event) => event.stopPropagation()}>
      {algorithms.map(({ label, value }) => (
        <Tooltip key={value} title={checksums[value]}>
          <Button
            size="small"
            type={checksums[value] ? 'primary' : 'default'}
            loading={loading === value}
            disabled={loading !== undefined && loading !== value}
            onClick={() => calculate(value)}
          >
            {label}
          </Button>
        </Tooltip>
      ))}
    </div>
  );

  return (
    <Popover content={content} title={t('image.checksum')} trigger="click" placement="left">
      <Tooltip title={t('image.checksum')}>
        <div
          className="flex h-[24px] w-[24px] items-center justify-center rounded text-neutral-300 hover:bg-neutral-500/50 hover:text-blue-500"
          onClick={(event) => event.stopPropagation()}
        >
          <HashIcon size={16} />
        </div>
      </Tooltip>
    </Popover>
  );
};
