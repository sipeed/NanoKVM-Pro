import React, { useEffect, useRef, useState } from 'react';
import { useAtom } from 'jotai';

import * as api from '@/api/hid.ts';
import { isModifier } from '@/lib/keymap';
import { leaderKeyAtom } from '@/jotai/keyboard.ts';

type SendKeyEvent = (type: 'keydown' | 'keyup', code: string) => void;

interface LeaderKeyState {
  code: string;
  isPressed: boolean;
  hasOtherKeys: boolean;
  recordMode: boolean;
  recordModifiers: string[];
}

interface LeaderKeyHandlers {
  handleKeyDown: (code: string, sendKeyEvent: SendKeyEvent) => boolean;
  handleKeyUp: (code: string, sendKeyEvent: SendKeyEvent) => boolean;
  reset: (sendKeyEvent?: SendKeyEvent) => void;
  recordMode: boolean;
  recordedKeys: string[];
}

export function useLeaderKey(pressedKeys: React.MutableRefObject<Set<string>>): LeaderKeyHandlers {
  const [leaderKeyCode, setLeaderKeyCode] = useAtom(leaderKeyAtom);
  const [recordMode, setRecordMode] = useState(false);
  const [recordedKeys, setRecordedKeys] = useState<string[]>([]);

  const leaderState = useRef<LeaderKeyState>({
    code: '',
    isPressed: false,
    hasOtherKeys: false,
    recordMode: false,
    recordModifiers: []
  });

  // Get leader key
  useEffect(() => {
    api.getLeaderKey().then((rsp: any) => {
      if (rsp.code === 0 && rsp.data?.key !== undefined) {
        setLeaderKeyCode(rsp.data.key);
      }
    });
  }, [setLeaderKeyCode]);

  // Update leader key
  useEffect(() => {
    leaderState.current.code = leaderKeyCode;
  }, [leaderKeyCode]);

  const updateRecordedKeysUI = () => {
    setRecordedKeys([...leaderState.current.recordModifiers]);
  };

  const handleKeyDown = (code: string, sendKeyEvent: SendKeyEvent): boolean => {
    const leader = leaderState.current;

    // Leader key pressed
    if (leader.code && code === leader.code) {
      leader.isPressed = true;
      leader.hasOtherKeys = false;
      pressedKeys.current.add(code);
      return true;
    }

    // Record mode
    if (leader.recordMode) {
      if (isModifier(code)) {
        // Modifier: send keydown immediately and track it
        if (!leader.recordModifiers.includes(code)) {
          leader.recordModifiers.push(code);
          updateRecordedKeysUI();
          sendKeyEvent('keydown', code);
        }
      } else {
        // Non-modifier: send keydown + keyup immediately
        sendKeyEvent('keydown', code);
        sendKeyEvent('keyup', code);
      }
      pressedKeys.current.add(code);
      return true;
    }

    // Hold mode
    if (leader.isPressed) {
      leader.hasOtherKeys = true;
      return false;
    }

    return false;
  };

  const handleKeyUp = (code: string, sendKeyEvent: SendKeyEvent): boolean => {
    const leader = leaderState.current;

    // Leader key released
    if (leader.code && code === leader.code) {
      if (!leader.hasOtherKeys) {
        if (leader.recordMode) {
          // Exit record mode: release all used modifiers
          for (const key of leader.recordModifiers) {
            sendKeyEvent('keyup', key);
          }
          leader.recordModifiers = [];
          setRecordedKeys([]);
          leader.recordMode = false;
          setRecordMode(false);
        } else {
          // Enter record mode
          leader.recordMode = true;
          setRecordMode(true);
        }
      }

      leader.isPressed = false;
      leader.hasOtherKeys = false;
      pressedKeys.current.delete(code);
      return true;
    }

    // In record mode, handle keyup
    if (leader.recordMode && code !== leader.code) {
      pressedKeys.current.delete(code);
      return true;
    }

    return false;
  };

  const reset = (sendKeyEvent?: SendKeyEvent) => {
    const leader = leaderState.current;

    // Release all used modifiers if sendKeyEvent is provided
    if (sendKeyEvent) {
      for (const key of leader.recordModifiers) {
        sendKeyEvent('keyup', key);
      }
    }

    leader.isPressed = false;
    leader.hasOtherKeys = false;
    leader.recordMode = false;
    leader.recordModifiers = [];
    setRecordMode(false);
    setRecordedKeys([]);
  };

  return { handleKeyDown, handleKeyUp, reset, recordMode, recordedKeys };
}
