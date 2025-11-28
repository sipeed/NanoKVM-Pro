import { useEffect, useRef } from 'react';
import { useAtomValue } from 'jotai';
import { useMediaQuery } from 'react-responsive';

import { client } from '@/lib/websocket.ts';
import { scrollIntervalAtom } from '@/jotai/mouse.ts';

import { MouseButton, MouseEvent } from './constants';

export const Absolute = () => {
  const isBigScreen = useMediaQuery({ minWidth: 850 });

  const scrollInterval = useAtomValue(scrollIntervalAtom);

  const lastScrollTimeRef = useRef(0);
  const touchStartTimeRef = useRef(0);
  const lastTouchYRef = useRef(0);
  const longPressTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isLongPressRef = useRef(false);
  const hasMoveRef = useRef(false);
  const isDraggingRef = useRef(false);
  const hasButtonPressedRef = useRef(false);
  const touchStartPosRef = useRef({ x: 0, y: 0 });

  const TAP_THRESHOLD = 8;
  const DRAG_THRESHOLD = 10;
  const VELOCITY_THRESHOLD = 0.3;

  const mouseButtonMapping = (button: number) => {
    const mappings = [MouseButton.Left, MouseButton.Wheel, MouseButton.Right];
    return mappings[button] || MouseButton.None;
  };

  // listen mouse events
  useEffect(() => {
    const canvas = document.getElementById('screen') as HTMLVideoElement;
    if (!canvas) return;

    canvas.addEventListener('mousedown', handleMouseDown);
    canvas.addEventListener('mouseup', handleMouseUp);
    canvas.addEventListener('mousemove', handleMouseMove);
    canvas.addEventListener('wheel', handleWheel);
    canvas.addEventListener('click', disableEvent);
    canvas.addEventListener('contextmenu', disableEvent);

    if (isBigScreen) {
      canvas.addEventListener('touchstart', handleTouchStart);
      canvas.addEventListener('touchmove', handleTouchMove);
      canvas.addEventListener('touchend', handleTouchEnd);
      canvas.addEventListener('touchcancel', handleTouchCancel);
    }

    // press button
    function handleMouseDown(event: any) {
      disableEvent(event);

      const button: MouseButton = mouseButtonMapping(event.button);
      if (button === MouseButton.None) return;

      const data = [2, MouseEvent.Down, button, 0, 0];
      client.send(data);
    }

    // release button
    function handleMouseUp(event: any) {
      disableEvent(event);

      const data = [2, MouseEvent.Up, MouseButton.None, 0, 0];
      client.send(data);
    }

    // mouse move
    function handleMouseMove(event: any) {
      disableEvent(event);

      const { x, y } = getCoordinate(event);
      const data = [2, MouseEvent.MoveAbsolute, MouseButton.None, x, y];
      client.send(data);
    }

    // mouse scroll
    function handleWheel(event: any) {
      disableEvent(event);

      const delta = Math.floor(event.deltaY);
      if (delta === 0) return;

      const currentTime = Date.now();
      if (currentTime - lastScrollTimeRef.current < scrollInterval) {
        return;
      }
      lastScrollTimeRef.current = currentTime;

      const data = [2, MouseEvent.Scroll, 0, 0, delta];
      client.send(data);
    }

    // touch start
    function handleTouchStart(event: any) {
      disableEvent(event);

      const touch = event.touches[0];
      if (!touch) return;

      // Reset states
      touchStartTimeRef.current = Date.now();
      lastTouchYRef.current = touch.clientY;
      isLongPressRef.current = false;
      hasMoveRef.current = false;
      isDraggingRef.current = false;
      hasButtonPressedRef.current = false;
      touchStartPosRef.current = { x: touch.clientX, y: touch.clientY };

      if (longPressTimerRef.current) {
        clearTimeout(longPressTimerRef.current);
      }

      const { x, y } = getCoordinate(touch);
      const moveData = [2, MouseEvent.MoveAbsolute, MouseButton.None, x, y];
      client.send(moveData);

      if (event.touches.length > 1) {
        return;
      }

      longPressTimerRef.current = setTimeout(() => {
        isLongPressRef.current = true;
        hasButtonPressedRef.current = true;

        if (navigator.vibrate) {
          navigator.vibrate(50);
        }

        const pressData = [2, MouseEvent.Down, MouseButton.Right, 0, 0];
        client.send(pressData);
      }, 500);
    }

    // touch move
    function handleTouchMove(event: any) {
      disableEvent(event);

      const touch = event.touches[0];
      if (!touch) return;

      const { x, y } = getCoordinate(touch);

      // Handle two-finger scroll first
      if (event.touches.length > 1) {
        const deltaY = touch.clientY - lastTouchYRef.current;
        lastTouchYRef.current = touch.clientY;

        if (Math.abs(deltaY) > 2) {
          const delta = Math.floor(deltaY);
          const scrollData = [2, MouseEvent.Scroll, 0, 0, delta];
          client.send(scrollData);
        }
        return;
      }

      const deltaX = Math.abs(touch.clientX - touchStartPosRef.current.x);
      const deltaY = Math.abs(touch.clientY - touchStartPosRef.current.y);
      const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);

      const timeDelta = Date.now() - touchStartTimeRef.current;
      const velocity = timeDelta > 0 ? distance / timeDelta : 0;

      const shouldStartDrag =
        distance > DRAG_THRESHOLD || (distance > TAP_THRESHOLD && velocity > VELOCITY_THRESHOLD);

      if (shouldStartDrag && !isDraggingRef.current && !isLongPressRef.current) {
        if (!hasMoveRef.current) {
          hasMoveRef.current = true;
        }

        if (longPressTimerRef.current) {
          clearTimeout(longPressTimerRef.current);
          longPressTimerRef.current = null;
        }

        if (!hasButtonPressedRef.current) {
          isDraggingRef.current = true;
          hasButtonPressedRef.current = true;
          const pressData = [2, MouseEvent.Down, MouseButton.Left, 0, 0];
          client.send(pressData);
        }
      }

      if (distance > TAP_THRESHOLD && !hasMoveRef.current) {
        hasMoveRef.current = true;
      }

      if (isDraggingRef.current || isLongPressRef.current) {
        const button = isLongPressRef.current ? MouseButton.Right : MouseButton.Left;
        const data = [2, MouseEvent.MoveAbsolute, button, x, y];
        client.send(data);
      } else if (hasMoveRef.current) {
        const data = [2, MouseEvent.MoveAbsolute, MouseButton.None, x, y];
        client.send(data);
      }
    }

    // touch end
    function handleTouchEnd(event: any) {
      disableEvent(event);

      if (longPressTimerRef.current) {
        clearTimeout(longPressTimerRef.current);
        longPressTimerRef.current = null;
      }

      if (!hasMoveRef.current && !isLongPressRef.current) {
        const pressData = [2, MouseEvent.Down, MouseButton.Left, 0, 0];
        client.send(pressData);

        setTimeout(() => {
          const releaseData = [2, MouseEvent.Up, MouseButton.None, 0, 0];
          client.send(releaseData);
        }, 50);
      } else if (hasButtonPressedRef.current) {
        const data = [2, MouseEvent.Up, MouseButton.None, 0, 0];
        client.send(data);
      }

      isLongPressRef.current = false;
      hasMoveRef.current = false;
      isDraggingRef.current = false;
      hasButtonPressedRef.current = false;
    }

    // touch cancel
    function handleTouchCancel(event: any) {
      disableEvent(event);

      if (longPressTimerRef.current) {
        clearTimeout(longPressTimerRef.current);
        longPressTimerRef.current = null;
      }

      if (hasButtonPressedRef.current) {
        const data = [2, MouseEvent.Up, MouseButton.None, 0, 0];
        client.send(data);
      }

      isLongPressRef.current = false;
      hasMoveRef.current = false;
      isDraggingRef.current = false;
      hasButtonPressedRef.current = false;
    }

    function getCorrectedCoords(clientX: number, clientY: number) {
      const rect = canvas.getBoundingClientRect();

      if (!canvas.videoWidth || !canvas.videoHeight) {
        const x = (clientX - rect.left) / rect.width;
        const y = (clientY - rect.top) / rect.height;
        return { x, y };
      }

      const videoRatio = canvas.videoWidth / canvas.videoHeight;
      const elementRatio = rect.width / rect.height;

      let renderedWidth = rect.width;
      let renderedHeight = rect.height;
      let offsetX = 0;
      let offsetY = 0;

      if (videoRatio > elementRatio) {
        renderedHeight = rect.width / videoRatio;
        offsetY = (rect.height - renderedHeight) / 2;
      } else {
        renderedWidth = rect.height * videoRatio;
        offsetX = (rect.width - renderedWidth) / 2;
      }

      const x = (clientX - rect.left - offsetX) / renderedWidth;
      const y = (clientY - rect.top - offsetY) / renderedHeight;

      return { x, y };
    }

    function getCoordinate(event: any): { x: number; y: number } {
      const { x, y } = getCorrectedCoords(event.clientX, event.clientY);

      const finalX = Math.max(0, Math.min(1, x));
      const finalY = Math.max(0, Math.min(1, y));

      const hexX = Math.floor(0x7fff * finalX) + 0x0001;
      const hexY = Math.floor(0x7fff * finalY) + 0x0001;

      return { x: hexX, y: hexY };
    }

    return () => {
      canvas.removeEventListener('mousemove', handleMouseMove);
      canvas.removeEventListener('mousedown', handleMouseDown);
      canvas.removeEventListener('mouseup', handleMouseUp);
      canvas.removeEventListener('wheel', handleWheel);
      canvas.removeEventListener('click', disableEvent);
      canvas.removeEventListener('contextmenu', disableEvent);
      canvas.removeEventListener('touchstart', handleTouchStart);
      canvas.removeEventListener('touchmove', handleTouchMove);
      canvas.removeEventListener('touchend', handleTouchEnd);
      canvas.removeEventListener('touchcancel', handleTouchCancel);

      if (longPressTimerRef.current) {
        clearTimeout(longPressTimerRef.current);
      }
    };
  }, [isBigScreen, scrollInterval]);

  // disable default events
  function disableEvent(event: any) {
    event.preventDefault();
    event.stopPropagation();
  }

  return <></>;
};
