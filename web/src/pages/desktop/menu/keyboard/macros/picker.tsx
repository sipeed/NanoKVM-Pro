import { type MouseEvent, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { useTranslation } from 'react-i18next';

import { clientToScreenPercent } from '@/lib/coords.ts';

type PickerProps = {
  open: boolean;
  onPick: (x: number, y: number) => void;
  onCancel: () => void;
};

export const CoordinatePicker = ({ open, onPick, onCancel }: PickerProps) => {
  const { t } = useTranslation();

  useEffect(() => {
    if (!open) return;

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        event.preventDefault();
        event.stopPropagation();
        onCancel();
      }
    };

    window.addEventListener('keydown', handleKeyDown, true);
    return () => window.removeEventListener('keydown', handleKeyDown, true);
  }, [open, onCancel]);

  if (!open) {
    return null;
  }

  function handleClick(event: MouseEvent<HTMLDivElement>) {
    event.preventDefault();
    event.stopPropagation();

    const coords = clientToScreenPercent(event.clientX, event.clientY);
    if (!coords) {
      onCancel();
      return;
    }

    onPick(coords.x, coords.y);
  }

  return createPortal(
    <div
      className="fixed inset-0 z-[3000] cursor-crosshair bg-black/20"
      onClick={handleClick}
      onContextMenu={(event) => {
        event.preventDefault();
        onCancel();
      }}
    >
      <div className="pointer-events-none absolute left-1/2 top-4 -translate-x-1/2 rounded bg-neutral-900/90 px-4 py-2 text-sm text-neutral-100 shadow">
        {t('keyboard.macro.pickHint')}
      </div>
    </div>,
    document.body
  );
};
