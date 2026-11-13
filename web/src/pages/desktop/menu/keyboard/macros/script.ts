import { getKeycode, isModifier } from '@/lib/keymap.ts';

import type { MouseButtonName, ParseResult, ScriptCommand } from './types.ts';

const MODIFIERS: Record<string, string> = {
  CTRL: 'ControlLeft',
  CONTROL: 'ControlLeft',
  SHIFT: 'ShiftLeft',
  ALT: 'AltLeft',
  GUI: 'MetaLeft',
  WIN: 'MetaLeft',
  WINDOWS: 'MetaLeft',
  META: 'MetaLeft',
  CMD: 'MetaLeft',
  COMMAND: 'MetaLeft'
};

const SPECIAL_KEYS: Record<string, string> = {
  ENTER: 'Enter',
  RETURN: 'Enter',
  ESC: 'Escape',
  ESCAPE: 'Escape',
  BACKSPACE: 'Backspace',
  BKSP: 'Backspace',
  TAB: 'Tab',
  SPACE: 'Space',
  DELETE: 'Delete',
  DEL: 'Delete',
  INSERT: 'Insert',
  INS: 'Insert',
  HOME: 'Home',
  END: 'End',
  PAGEUP: 'PageUp',
  PGUP: 'PageUp',
  PGDN: 'PageDown',
  PAGEDOWN: 'PageDown',
  UP: 'ArrowUp',
  DOWN: 'ArrowDown',
  LEFT: 'ArrowLeft',
  RIGHT: 'ArrowRight',
  CAPSLOCK: 'CapsLock',
  CAPS: 'CapsLock',
  PRINTSCREEN: 'PrintScreen',
  PRINT: 'PrintScreen',
  SCROLLLOCK: 'ScrollLock',
  PAUSE: 'Pause',
  BREAK: 'Pause',
  MENU: 'ContextMenu',
  APP: 'ContextMenu',
  NUMLOCK: 'NumLock'
};

const BUTTONS: Record<string, MouseButtonName> = {
  LEFT: 'left',
  RIGHT: 'right',
  MIDDLE: 'middle',
  MID: 'middle'
};

export function parseScript(script: string): ParseResult {
  const commands: ScriptCommand[] = [];
  const lines = script.replace(/\r\n/g, '\n').split('\n');

  for (let i = 0; i < lines.length; i++) {
    const lineNo = i + 1;
    const raw = lines[i].trim();
    if (!raw || raw.startsWith('REM') || raw.startsWith('#')) {
      continue;
    }

    const parsed = parseLine(raw);
    if (typeof parsed === 'string') {
      return { ok: false, error: parsed, line: lineNo };
    }
    commands.push(parsed);
  }

  if (!commands.length) {
    return { ok: false, error: 'empty script', line: 1 };
  }

  return { ok: true, commands };
}

function parseLine(raw: string): ScriptCommand | string {
  const [command, rest] = splitCommand(raw);
  const name = command.toUpperCase();

  if (name === 'DELAY') {
    if (!rest) {
      return 'DELAY needs 0..60000 ms';
    }
    const ms = Number(rest.split(/\s+/)[0]);
    if (!Number.isInteger(ms) || ms < 0 || ms > 60000) {
      return 'DELAY needs 0..60000 ms';
    }
    return { type: 'delay', ms };
  }

  if (name === 'STRING' || name === 'STRINGLN') {
    const text = unquote(rest);
    if (!text) {
      return `${name} needs text`;
    }
    return { type: 'string', text: name === 'STRINGLN' ? `${text}\n` : text };
  }

  if (name === 'CLICK' || name === 'DBLCLICK' || name === 'RIGHTCLICK' || name === 'MIDDLECLICK') {
    return parseClick(name, rest);
  }

  if (name === 'RUN') {
    const target = unquote(rest).trim();
    if (!target) {
      return 'RUN needs a macro name';
    }
    if (target.length > 64) {
      return 'RUN name is too long';
    }
    return { type: 'run', name: target };
  }

  return parseKeyLine(raw);
}

function parseClick(name: string, rest: string): ScriptCommand | string {
  const parts = rest.split(/\s+/).filter(Boolean);
  if (parts.length < 2 || parts.length > 3) {
    return `${name} needs X Y [LEFT|RIGHT|MIDDLE]`;
  }

  const x = Number(parts[0]);
  const y = Number(parts[1]);
  if (!Number.isFinite(x) || !Number.isFinite(y) || x < 0 || x > 100 || y < 0 || y > 100) {
    return 'coordinates must be 0..100 percent';
  }

  let button: MouseButtonName = 'left';
  if (name === 'RIGHTCLICK') {
    button = 'right';
  } else if (name === 'MIDDLECLICK') {
    button = 'middle';
  }

  if (parts[2]) {
    const mapped = BUTTONS[parts[2].toUpperCase()];
    if (!mapped) {
      return 'unknown mouse button';
    }
    button = mapped;
  }

  return {
    type: 'click',
    x,
    y,
    button,
    times: name === 'DBLCLICK' ? 2 : 1
  };
}

function parseKeyLine(raw: string): ScriptCommand | string {
  const tokens = raw
    .split(/[\s-]+/)
    .map((token) => token.trim())
    .filter(Boolean);

  if (!tokens.length || tokens.length > 6) {
    return 'invalid key command';
  }

  const codes: string[] = [];
  for (let i = 0; i < tokens.length; i++) {
    const token = tokens[i];
    const isLast = i === tokens.length - 1;
    const code = isLast || !MODIFIERS[token.toUpperCase()] ? tokenToCode(token) : MODIFIERS[token.toUpperCase()];
    if (!code) {
      return `unknown key: ${token}`;
    }
    if (codes.includes(code)) {
      continue;
    }
    codes.push(code);
  }

  const mainKey = codes[codes.length - 1];
  if (isModifier(mainKey) && codes.length === 1) {
    return { type: 'key', codes };
  }

  for (const code of codes) {
    if (getKeycode(code) === undefined && !isModifier(code)) {
      return `unknown key: ${code}`;
    }
  }

  return { type: 'key', codes };
}

function tokenToCode(token: string): string | null {
  const upper = token.toUpperCase();
  if (MODIFIERS[upper]) {
    return MODIFIERS[upper];
  }
  if (SPECIAL_KEYS[upper]) {
    return SPECIAL_KEYS[upper];
  }
  if (/^F([1-9]|1[0-9]|2[0-4])$/.test(upper)) {
    return upper;
  }
  if (/^[A-Z]$/.test(upper)) {
    return `Key${upper}`;
  }
  if (/^[0-9]$/.test(token)) {
    return `Digit${token}`;
  }
  return null;
}

function splitCommand(raw: string): [string, string] {
  const index = raw.search(/\s/);
  if (index < 0) {
    return [raw, ''];
  }
  return [raw.slice(0, index), raw.slice(index + 1).trim()];
}

function unquote(value: string): string {
  if (
    (value.startsWith('"') && value.endsWith('"') && value.length >= 2) ||
    (value.startsWith("'") && value.endsWith("'") && value.length >= 2)
  ) {
    return value.slice(1, -1);
  }
  return value;
}

export function countScriptCommands(script: string): number {
  const parsed = parseScript(script);
  return parsed.ok ? parsed.commands.length : 0;
}
