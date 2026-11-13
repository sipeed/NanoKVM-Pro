export type MouseButtonName = 'left' | 'right' | 'middle';

export type ScriptCommand =
  | { type: 'delay'; ms: number }
  | { type: 'string'; text: string }
  | { type: 'key'; codes: string[] }
  | { type: 'click'; x: number; y: number; button: MouseButtonName; times: number }
  | { type: 'run'; name: string };

export interface Macro {
  id?: string;
  name: string;
  script: string;
}

export type ParseResult =
  | { ok: true; commands: ScriptCommand[] }
  | { ok: false; error: string; line: number };
