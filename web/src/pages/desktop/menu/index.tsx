import { useEffect, useRef, useState } from 'react';
import { Divider } from 'antd';
import clsx from 'clsx';
import { useAtom, useAtomValue } from 'jotai';
import { GripVerticalIcon } from 'lucide-react';
import Draggable, { DraggableData, DraggableEvent } from 'react-draggable';

import * as api from '@/api/vm.ts';
import * as ls from '@/lib/localstorage.ts';
import { menuDisabledItemsAtom, submenuOpenCountAtom } from '@/jotai/settings.ts';

import { AIAssistant } from './assistant/index.tsx';
import { Fullscreen } from './fullscreen';
import { Image } from './image';
import { Keyboard } from './keyboard';
import { Mouse } from './mouse';
import { Collapse, Expand } from './operations';
import { Power } from './power';
import { Screen } from './screen';
import { Script } from './script';
import { Settings } from './settings';
import { Terminal } from './terminal';
import { Volume } from './volume';
import { Wol } from './wol';

export const Menu = () => {
  const [menuDisabledItems, setMenuDisabledItems] = useAtom(menuDisabledItemsAtom);
  const submenuOpenCount = useAtomValue(submenuOpenCountAtom);

  const [menuStatus, setMenuStatus] = useState({
    expanded: true,
    moved: false,
    hovered: false,
    hidden: false
  });
  const [menuBounds, setMenuBounds] = useState({ left: 0, right: 0, top: 0, bottom: 0 });

  const nodeRef = useRef<HTMLDivElement | null>(null);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const HIDE_TIMEOUT = 8000;

  // init menu
  useEffect(() => {
    getMenuBarConfig();
    startCountdown();
    handleResize();

    window.addEventListener('resize', handleResize);

    return () => {
      stopCountdown();
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  // handle menu bounds
  useEffect(() => {
    handleResize();
  }, [menuDisabledItems, menuStatus.expanded]);

  // handle hide timer
  useEffect(() => {
    if (submenuOpenCount === 0 && !menuStatus.hovered) {
      startCountdown();
      return;
    }

    stopCountdown();
  }, [submenuOpenCount, menuStatus.hovered]);

  async function getMenuBarConfig() {
    try {
      const rsp = await api.getMenuBarConfig();
      if (rsp.code === 0 && rsp.data?.disabledItems?.length) {
        setMenuDisabledItems(rsp.data.disabledItems);
        return;
      }

      // compatible with old versions
      const items = ls.getMenuDisabledItems();
      setMenuDisabledItems(items);
    } catch (err) {
      console.log(err);
    }
  }

  function handleResize() {
    if (!nodeRef.current) return;

    const elementRect = nodeRef.current.getBoundingClientRect();
    const width = (window.innerWidth - elementRect.width) / 2;

    setMenuBounds({
      left: -width,
      top: -10,
      right: width,
      bottom: window.innerHeight - elementRect.height - 10
    });
  }

  function startCountdown() {
    if (submenuOpenCount > 0 || !menuStatus.expanded || menuStatus.moved) {
      return;
    }

    stopCountdown();

    timerRef.current = setTimeout(() => {
      setMenuStatus((status) => ({
        ...status,
        hidden: true
      }));
    }, HIDE_TIMEOUT);
  }

  function stopCountdown() {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
  }

  function handleExpanded(expanded: boolean) {
    setMenuStatus((status) => ({
      ...status,
      expanded
    }));
  }

  function handleMoved(_e: DraggableEvent, data: DraggableData) {
    if (data.x !== 0 || data.y !== 0) {
      setMenuStatus((status) => ({
        ...status,
        moved: true
      }));
    }
  }

  function handleHovered(hovered: boolean) {
    setMenuStatus((status) => ({
      ...status,
      hovered,
      hidden: hovered ? false : status.hidden
    }));
  }

  function isEnabled(items: string | string[]) {
    if (typeof items === 'string') {
      return !menuDisabledItems.includes(items);
    }
    return items.some((item) => !menuDisabledItems.includes(item));
  }

  return (
    <Draggable
      nodeRef={nodeRef}
      bounds={menuBounds}
      handle="strong"
      positionOffset={{ x: '-50%', y: '0%' }}
      onStop={handleMoved}
    >
      <div
        ref={nodeRef}
        className="fixed left-1/2 top-[10px] z-[1000] -translate-x-1/2"
        onMouseEnter={() => handleHovered(true)}
        onMouseLeave={() => handleHovered(false)}
        onBlur={() => handleHovered(false)}
      >
        {/* Trigger area for auto-show when hidden */}
        {menuStatus.expanded && (
          <div className="absolute -top-[10px] left-0 right-0 h-[46px] w-full bg-transparent" />
        )}

        {/* Menubar */}
        <div className="sticky top-[10px] flex w-full justify-center">
          <div
            className={clsx(
              'h-[36px] items-center rounded bg-neutral-800/80 pl-1 pr-2 transition-all duration-300',
              menuStatus.expanded ? 'flex' : 'hidden',
              menuStatus.hidden ? '-translate-y-[110%] opacity-80' : 'translate-y-0 opacity-100'
            )}
          >
            <strong>
              <div className="flex h-[30px] cursor-move select-none items-center justify-center pl-1 text-neutral-500">
                <GripVerticalIcon size={18} />
              </div>
            </strong>
            <Divider type="vertical" />

            <Screen />
            {isEnabled('volume') && <Volume />}

            <Keyboard />
            <Mouse />
            <Divider type="vertical" />

            {isEnabled('image') && <Image />}
            {isEnabled('script') && <Script />}
            {isEnabled('assistant') && <AIAssistant />}
            {isEnabled('terminal') && <Terminal />}
            {isEnabled('wol') && <Wol />}

            {isEnabled(['image', 'script', 'assistant', 'terminal', 'wol']) && (
              <Divider type="vertical" />
            )}

            {isEnabled('power') && (
              <>
                <Power />
                <Divider type="vertical" />
              </>
            )}

            <Settings />
            {isEnabled('fullscreen') && <Fullscreen />}
            {isEnabled('collapse') && <Collapse toggleMenu={handleExpanded} />}
          </div>
        </div>

        {/* Menubar expand button */}
        {!menuStatus.expanded && <Expand toggleMenu={handleExpanded} />}
      </div>
    </Draggable>
  );
};
