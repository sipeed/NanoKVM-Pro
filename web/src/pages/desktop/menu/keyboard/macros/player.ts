import { getDefaultStore } from 'jotai';

import { KeyboardReport } from '@/lib/keyboard.ts';
import { MouseReportAbsolute } from '@/lib/mouse.ts';
import { percentToHid } from '@/lib/coords.ts';
import { client, MessageEvent } from '@/lib/websocket.ts';
import { isKeyboardEnableAtom } from '@/jotai/keyboard.ts';
import { isMacroPlayingAtom } from '@/jotai/mouse.ts';

import { charToCodes } from './chars.ts';
import { parseScript } from './script.ts';
import type { Macro, MouseButtonName, ScriptCommand } from './types.ts';

const BUTTON_INDEX: Record<MouseButtonName, number> = {
  left: 0,
  middle: 1,
  right: 2
};

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function sendKeyboard(report: Uint8Array) {
  client.send(new Uint8Array([MessageEvent.Keyboard, ...report]));
}

function sendMouse(report: Uint8Array) {
  client.send(new Uint8Array([MessageEvent.Mouse, ...report]));
}

async function clickAt(
  mouse: MouseReportAbsolute,
  x: number,
  y: number,
  button: MouseButtonName,
  times: number
) {
  const hidX = percentToHid(x);
  const hidY = percentToHid(y);
  const buttonIndex = BUTTON_INDEX[button] ?? 0;

  sendMouse(mouse.buildReport(hidX, hidY, 0));
  await sleep(30);

  for (let i = 0; i < times; i++) {
    mouse.buttonDown(buttonIndex);
    sendMouse(mouse.buildButtonReport(hidX, hidY));
    await sleep(50);

    mouse.buttonUp(buttonIndex);
    sendMouse(mouse.buildButtonReport(hidX, hidY));

    if (i < times - 1) {
      await sleep(80);
    }
  }
}

async function pressCodes(codes: string[]) {
  const keyboard = new KeyboardReport();

  for (const code of codes) {
    sendKeyboard(keyboard.keyDown(code));
    await sleep(20);
  }

  await sleep(50);
  sendKeyboard(keyboard.reset());
}

async function typeText(text: string) {
  for (const char of text) {
    const codes = charToCodes(char);
    if (!codes) {
      continue;
    }
    await pressCodes(codes);
    await sleep(50);
  }
}

export async function playMacro(macro: Macro, macros: Macro[] = []): Promise<string | null> {
  const store = getDefaultStore();
  const previousKeyboard = store.get(isKeyboardEnableAtom);

  store.set(isKeyboardEnableAtom, false);
  store.set(isMacroPlayingAtom, true);

  try {
    return await playCommands(macro, macros, new MouseReportAbsolute(), 0);
  } finally {
    store.set(isMacroPlayingAtom, false);
    store.set(isKeyboardEnableAtom, previousKeyboard);
  }
}

async function playCommands(
  macro: Macro,
  macros: Macro[],
  mouse: MouseReportAbsolute,
  depth: number
): Promise<string | null> {
  const parsed = parseScript(macro.script);
  if (!parsed.ok) {
    return `${macro.name}: line ${parsed.line}: ${parsed.error}`;
  }

  for (const command of parsed.commands) {
    const error = await playCommand(command, macro, macros, mouse, depth);
    if (error) {
      return error;
    }

    if (command.type !== 'delay') {
      await sleep(30);
    }
  }

  return null;
}

async function playCommand(
  command: ScriptCommand,
  macro: Macro,
  macros: Macro[],
  mouse: MouseReportAbsolute,
  depth: number
): Promise<string | null> {
  switch (command.type) {
    case 'click':
      await clickAt(mouse, command.x, command.y, command.button, command.times);
      return null;
    case 'key':
      await pressCodes(command.codes);
      return null;
    case 'string':
      await typeText(command.text);
      return null;
    case 'delay':
      await sleep(command.ms);
      return null;
    case 'run':
      return playNestedMacro(command.name, macro, macros, mouse, depth);
  }
}

async function playNestedMacro(
  name: string,
  parent: Macro,
  macros: Macro[],
  mouse: MouseReportAbsolute,
  depth: number
): Promise<string | null> {
  if (depth >= 1) {
    return `${parent.name}: nested RUN is limited to 1 level`;
  }

  const child = findMacro(macros, name);
  if (!child) {
    return `${parent.name}: macro not found: ${name}`;
  }

  return playCommands(child, macros, mouse, depth + 1);
}

function findMacro(macros: Macro[], name: string): Macro | undefined {
  const needle = name.trim().toLowerCase();
  return macros.find((macro) => macro.name.trim().toLowerCase() === needle);
}
